package main

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
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

	log.Debug("Registering prometheus metrics")

	m := &metrics{
		channelBalanceGauge: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "liquidator.channel_balance",
			Help: "The total number of processed events",
		},
			[]string{"channel_id", "local_node_pubkey", "remote_node_pubkey", "local_node_alias", "remote_node_alias", "active", "initiator"},
		),
	}

	//Golang collector
	reg.MustRegister(collectors.NewGoCollector())
	//Process collector
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	reg.MustRegister(m.channelBalanceGauge)

	prometheusMetrics = m
}

// Entrypoint of liquidator main logic
func startLiquidator() {

	var wg = &sync.WaitGroup{}

	//Create a new prometheus registry
	reg := prometheus.NewRegistry()

	InitMetrics(reg)

	log.Debug("Prometheus metrics registered")

	log.Debug("Starting /metrics endpoint")
	//Start a http server to expose metrics
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg, EnableOpenMetrics: true}))
	go http.ListenAndServe(":9000", nil)

	log.Debug("/metrics endpoint started")

	//For each node in nodesHosts, connect to the node and get the list of channels

	for i, node := range nodesHosts {

		log.Infof("Starting monitoring for node: %v", node)

		//Generate TLS credentials from directory
		tlsCertEncoded := nodesTLSCerts[i]
		tlsCertDecoded, err := base64.StdEncoding.DecodeString(tlsCertEncoded)
		if err != nil {
			log.Fatalf("Failed to decode TLS cert: %v", err)
		}

		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(tlsCertDecoded) {
			log.Fatalf("credentials: failed to append certificates")
		}

		creds := credentials.NewClientTLSFromCert(cp, "")

		if err != nil {
			log.Fatalf("Failed to load credentials: %v", err)
		}

		conn, err := grpc.Dial(node, grpc.WithTransportCredentials(creds))

		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		lightningClient := lnrpc.NewLightningClient(conn)

		ctx := context.Background()

		macaroon := nodesMacaroons[i]

		if macaroon == "" {
			log.Fatalf("No macaroon provided for node %v", node)
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
		
		err := fmt.Errorf("channel capacity is <= 0")
		log.Error(err)
		return -1, err
	}

	localBalance := float64(channel.GetLocalBalance())

	channelBalanceRatioInt := localBalance / capacity

	//Truncate channelbalance to 2 decimal places
	channelBalanceRatio := float64(int(channelBalanceRatioInt*100)) / 100

	//Check that the ratio is between 0 and 1
	if channelBalanceRatio > 1 || channelBalanceRatio < 0 {
		
		err := fmt.Errorf("channel balance ratio is not between 0 and 1")
		log.Error(err)
		return -1, err
	}

	return channelBalanceRatio, nil

}

// Locking fuction to be used in a goroutine to monitor channels
func monitorChannels(nodeHost string, macaroon string, lightningClient lnrpc.LightningClient, ctx context.Context) {

	//Defer a recover function to catch panics
	defer func() {
		if r := recover(); r != nil {
			log.Error("Recovered panic in monitorChannels function, retrying...", r)

			//Sleep for 10 seconds before restarting the monitoring
			time.Sleep(10 * time.Second)

			//Restart monitoring via recursion
			monitorChannels(nodeHost, macaroon, lightningClient, ctx)
		}

	}()

	//Check that nodehost matches host:port string
	if nodeHost == "" {
		error := fmt.Errorf("nodeHost is empty")
		log.Error(error)
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
			log.Errorf("Error listing channels: %v", err)
		}

		if response == nil || len(response.Channels) == 0 {
			log.Errorf("No channels found for node %v", nodeHost)

			time.Sleep(1 * time.Second)

			continue
		}

		//Iterate over response channels
		for _, channel := range response.Channels {
			log.Debugf("Monitoring Channel Id: %v", channel.ChanId)
			//Record the channel balance in a prometheus gauge
			channelBalanceRatio, err := recordChannelBalance(channel)
			//Log channel balance ratio in debug mode
			log.Debugf("Channel balance ratio for node %v channel %v is %v", nodeHost, channel.GetChanId(), channelBalanceRatio)

			if err != nil {
				log.Errorf("Error calculating channel balance: %v", err)
			}

			//ChannelId uint to string
			channelId := fmt.Sprint(channel.GetChanId())

			//Get the node info
			localNodeInfo, err := getInfo(&lightningClient, &context)

			if err != nil {
				log.Errorf("Error getting local node info: %v", err)
			}

			remoteNodeInfo, err := getNodeInfo(channel.RemotePubkey, &lightningClient, &context)

			if err != nil {
				log.Errorf("Error getting remote node info: %v", err)

			}

			//Set the channel balance in the gauge
			localPubKey := localNodeInfo.IdentityPubkey
			remotePubKey := channel.GetRemotePubkey()
			localAlias := localNodeInfo.Alias
			remoteAlias := remoteNodeInfo.GetNode().Alias
			//Channel Active to string
			active := strconv.FormatBool(channel.GetActive())
			initiator := strconv.FormatBool(channel.GetInitiator())
			
			prometheusMetrics.channelBalanceGauge.With(prometheus.Labels{
				"channel_id":         channelId,
				"local_node_pubkey":  localPubKey,
				"remote_node_pubkey": remotePubKey,
				"local_node_alias":   localAlias,
				"remote_node_alias":  remoteAlias,
				"active":	active,
				"initiator":	initiator,
				}).Set(channelBalanceRatio)

			time.Sleep(pollingInterval)
		}

	}

}

func getInfo(lightningClient *lnrpc.LightningClient, context *context.Context) (*lnrpc.GetInfoResponse, error) {

	//Call GetInfo method of lightning client
	response, err := (*lightningClient).GetInfo(*context, &lnrpc.GetInfoRequest{})

	if err != nil {
		log.Errorf("Error getting info: %v", err)
		return nil, err
	}

	return response, nil
}

// Gets the info from the node with the given pubkey
func getNodeInfo(pubkey string, lightningClient *lnrpc.LightningClient, context *context.Context) (*lnrpc.NodeInfo, error) {

	if pubkey == "" {
		error := fmt.Errorf("pubkey is empty")
		log.Error(error, "pubkey is empty")
		return nil, error

	}

	//Call GetNodeInfo method of lightning client with metadata headers
	response, err := (*lightningClient).GetNodeInfo(*context, &lnrpc.NodeInfoRequest{
		PubKey: pubkey,
	})

	if err != nil {
		log.Errorf("Error getting node info: %v", err)
		return nil, err
	}

	return response, nil
}
