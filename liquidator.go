package main

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var (
	prometheusMetrics *metrics
)

type metrics struct {
	channelBalanceGauge prometheus.GaugeVec
}

// Inits the prometheusMetrics global metric struct
func InitMetrics(reg prometheus.Registerer) {
	m := &metrics{
		channelBalanceGauge: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "liquidator_channel_balance",
			Help: "The total number of processed events",
		},
			[]string{"channel_id", "local_node_pubkey", "remote_node_pubkey", "local_node_alias", "remote_node_alias"},
		),
	}

	//Golang collector
	reg.MustRegister(collectors.NewGoCollector())
	//Process collector
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	reg.MustRegister(m.channelBalanceGauge)

	prometheusMetrics = m
}

func startLiquidator() {

	var wg = &sync.WaitGroup{}

	//Create a new prometheus registry
	reg := prometheus.NewRegistry()

	InitMetrics(reg)

	//Start a http server to expose metrics
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg, EnableOpenMetrics: true}))
	go http.ListenAndServe(":9000", nil)

	//For each node in nodesHosts, connect to the node and get the list of channels

	for i, node := range nodesHosts {

		InfoLog.Println("starting monitoring for node: %v", node)

		//Generate TLS credentials from directory
		tlsCertEncoded := nodesTLSCerts[i]
		tlsCertDecoded, err := base64.StdEncoding.DecodeString(tlsCertEncoded)
		if err != nil {
			ErrorLog.Fatalf("Failed to decode TLS cert: %v", err)
		}

		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(tlsCertDecoded) {
			ErrorLog.Fatalf("credentials: failed to append certificates")
		}

		creds := credentials.NewClientTLSFromCert(cp, "")

		if err != nil {
			ErrorLog.Fatalf("Failed to load credentials: %v", err)
		}

		conn, err := grpc.Dial(node, grpc.WithTransportCredentials(creds))

		if err != nil {
			ErrorLog.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		lightningClient := lnrpc.NewLightningClient(conn)

		ctx := context.Background()

		macaroon := nodesMacaroons[i]

		if macaroon == "" {
			ErrorLog.Fatalf("No macaroon provided for node %v", node)
		}

		wg.Add(1)
		go monitorChannels(node, macaroon, lightningClient, ctx)

	}

	wg.Wait()

	//TODO Graceful shutdown
}

// Record the channel balance in a prometheus gauge
func recordChannelBalance(channel *lnrpc.Channel) (float64, error) {

	capacity := float64(channel.GetCapacity())

	if capacity <= 0 {
		ErrorLog.Printf("Channel capacity is <= 0")
		err := fmt.Errorf("channel capacity is <= 0")
		return -1, err
	}

	localBalance := float64(channel.GetLocalBalance())

	channelBalanceRatioInt := localBalance / capacity

	//Truncate channelbalance to 2 decimal places
	channelBalanceRatio := float64(int(channelBalanceRatioInt*100)) / 100

	//Check that the ration is between 0 and 1
	if channelBalanceRatio > 1 || channelBalanceRatio < 0 {
		ErrorLog.Println("Channel balance ratio is not between 0 and 1")
		err := fmt.Errorf("channel balance ratio is not between 0 and 1")
		return -1, err
	}

	return channelBalanceRatio, nil

}

// Locking fuction to be used in a goroutine to monitor channels
func monitorChannels(nodeHost string, macaroon string, lightningClient lnrpc.LightningClient, ctx context.Context) {

	//Check that nodehost matches host:port string
	if nodeHost == "" {
		error := fmt.Errorf("nodeHost is empty")
		ErrorLog.Println(error, "nodeHost is empty")
	}

	md := metadata.New(map[string]string{"macaroon": macaroon})

	context := metadata.NewOutgoingContext(ctx, md)

	//Infinite loop to monitor channels
	for {

		//Call ListChannels method of lightning client with metadata headers
		response, err := lightningClient.ListChannels(context, &lnrpc.ListChannelsRequest{
			ActiveOnly: false,
		})

		if err != nil {
			ErrorLog.Printf("Error listing channels: %v", err)
		}

		//Iterate over response channels
		for _, channel := range response.Channels {
			//Record the channel balance in a prometheus gauge
			channelBalanceRatio, err := recordChannelBalance(channel)

			if err != nil {
				ErrorLog.Printf("Error calculating channel balance: %v", err)
			}

			//ChannelId uint to string
			channelId := fmt.Sprint(channel.GetChanId())

			//Get the node info
			localNodeInfo, err := getInfo(&lightningClient, &context)

			if err != nil {
				ErrorLog.Printf("Error getting local node info: %v", err)
			}

			remoteNodeInfo, err := getNodeInfo(channel.RemotePubkey, &lightningClient, &context)

			if err != nil {
				ErrorLog.Printf("Error getting remote node info: %v", err)

			}

			//Set the channel balance in the gauge
			prometheusMetrics.channelBalanceGauge.With(prometheus.Labels{
				"channel_id":         channelId,
				"local_node_pubkey":  localNodeInfo.IdentityPubkey,
				"remote_node_pubkey": channel.GetRemotePubkey(),
				"local_node_alias":   localNodeInfo.Alias,
				"remote_node_alias":  remoteNodeInfo.GetNode().Alias}).Set(channelBalanceRatio)

			time.Sleep(pollingInterval)
		}
	}

}

func getInfo(lightningClient *lnrpc.LightningClient, context *context.Context) (*lnrpc.GetInfoResponse, error) {

	//Call GetInfo method of lightning client
	response, err := (*lightningClient).GetInfo(*context, &lnrpc.GetInfoRequest{})

	if err != nil {
		ErrorLog.Printf("Error getting info: %v", err)
		return nil, err
	}

	return response, nil
}

// Gets the info from the node with the given pubkey
func getNodeInfo(pubkey string, lightningClient *lnrpc.LightningClient, context *context.Context) (*lnrpc.NodeInfo, error) {

	//TODO
	if pubkey == "" {
		error := fmt.Errorf("pubkey is empty")
		ErrorLog.Println(error, "pubkey is empty")
		return nil, error

	}

	//Call GetNodeInfo method of lightning client with metadata headers
	response, err := (*lightningClient).GetNodeInfo(*context, &lnrpc.NodeInfoRequest{
		PubKey: pubkey,
	})

	if err != nil {
		ErrorLog.Printf("Error getting node info: %v", err)
		return nil, err
	}

	return response, nil
}
