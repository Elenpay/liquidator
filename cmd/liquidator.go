package cmd

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type metrics struct {
	channelBalanceGauge prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		channelBalanceGauge: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "liquidator_channel_balance",
			Help: "The total number of processed events",
		},
			[]string{"channel_id", "local_node_pubkey", "remote_node_pubkey", "local_node_alias", "remote_node_alias"},
		),
	}
	reg.MustRegister(m.channelBalanceGauge)
	return m
}

func startLiquidator() {

	//Create a new prometheus registry
	reg := prometheus.NewRegistry()

	metrics := NewMetrics(reg)

	//Start a http server to expose metrics
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	go http.ListenAndServe(":9000", nil)

	//For each node in nodesHosts, connect to the node and get the list of channels
	for _, node := range nodesHosts {
		//Generate TLS credentials from directory
		creds, err := credentials.NewClientTLSFromFile("/Users/joseap/.polar/networks/1/volumes/lnd/alice/tls.cert", "")

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

		//Call connectToServer function with url of the server every minute
		for {
			err = monitorChannels(node, metrics, lightningClient, ctx)
			
			if err != nil {
				handleErr(err,"Error monitoring channels: %v")
			}

			time.Sleep(1 * time.Second)
		}
	}

}

// Record the channel balance in a prometheus gauge
func recordChannelBalance(metrics *metrics, channel *lnrpc.Channel) {

	channelBalanceRatioInt := float64(channel.GetLocalBalance()) / float64(channel.GetCapacity())

	//Truncate channelbalance to 2 decimal places
	channelBalanceRatio := float64(int(channelBalanceRatioInt*100)) / 100

	//Check that the ration is between 0 and 1
	if channelBalanceRatio > 1 || channelBalanceRatio < 0 {
		ErrorLog.Fatalf("Channel balance ratio is not between 0 and 1")
	}

	//ChannelId uint to string
	channelId := fmt.Sprint(channel.GetChanId())

	//Set the channel balance in the gauge
	metrics.channelBalanceGauge.With(prometheus.Labels{
		"channel_id":         channelId,
		"local_node_pubkey":  "123", //TODO
		"remote_node_pubkey": channel.GetRemotePubkey(),
		"local_node_alias":   "alice",
		"remote_node_alias":  "bob"}).Set(channelBalanceRatio)

	time.Sleep(2 * time.Second)

}

// func that receives a url and connects to a grpc server
func monitorChannels(nodeHost string, metrics *metrics, lightningClient lnrpc.LightningClient, ctx context.Context) error {

	//Check that nodehost matches host:port string
	if nodeHost == "" {
		error := fmt.Errorf("nodeHost is empty")
		handleErr(error,"nodeHost is empty")
		return error

	}

	//TODO Fix this
	md := metadata.New(map[string]string{"macaroon": "0201036c6e6402ac01030a108ae5b2928f746a822b04a9b2848eb0321201301a0f0a07616464726573731204726561641a0c0a04696e666f1204726561641a100a08696e766f696365731204726561641a100a086d616361726f6f6e1204726561641a0f0a076d6573736167651204726561641a100a086f6666636861696e1204726561641a0f0a076f6e636861696e1204726561641a0d0a0570656572731204726561641a0e0a067369676e657212047265616400000620dfb922212dc2831973c90712913f0bfea68916a640c4d8475359ea593e6789ea"})

	context := metadata.NewOutgoingContext(ctx, md)

	//Call ListChannels method of lightning client with metadata headers
	response, err := lightningClient.ListChannels(context, &lnrpc.ListChannelsRequest{
		ActiveOnly: false,
	})

	if err != nil {
		ErrorLog.Printf("ListChannels fail", err)
	}

	//Iterate over response channels
	for _, channel := range response.Channels {
		//Record the channel balance in a prometheus gauge
		recordChannelBalance(metrics, channel)

	}

	return nil
}
