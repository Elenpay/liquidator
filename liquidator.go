package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Elenpay/liquidator/cache"
	"github.com/Elenpay/liquidator/helper"
	"github.com/Elenpay/liquidator/nodeguard"
	"github.com/Elenpay/liquidator/provider"
	"github.com/Elenpay/liquidator/rpc"
	"github.com/lightninglabs/loop/looprpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	prometheusMetrics *metrics
	rulesCache        cache.Cache
)

type metrics struct {
	channelBalanceGauge prometheus.GaugeVec
}

// Inits the prometheusMetrics global metric struct
func initMetrics(reg prometheus.Registerer) {

	log.Debug("Registering prometheus metrics")

	m := &metrics{
		channelBalanceGauge: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "liquidator_channel_balance",
			Help: "The total number of processed events",
		},
			[]string{"chan_id", "local_node_pubkey", "remote_node_pubkey", "local_node_alias", "remote_node_alias", "active", "initiator"},
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

	initMetrics(reg)

	log.Debug("Prometheus metrics registered")

	log.Debug("Starting /metrics endpoint")
	//Start a http server to expose metrics
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg, EnableOpenMetrics: true}))
	go http.ListenAndServe(":9000", nil)

	log.Debug("/metrics endpoint started")

	//Start cache to store liquidation rules in case of nodeguard failure
	cache, err := cache.NewCache()
	if err != nil {
		log.Fatalf("failed to create cache: %v", err)
	}
	rulesCache = cache

	//For each node in nodesHosts, connect to the node and get the list of channels

	for i, nodeEndpoint := range nodesHosts {

		log.Infof("starting monitoring for node: %v", nodeEndpoint)

		nodeTLSCertEncoded := nodesTLSCerts[i]

		loopdTLSCertEncoded := loopdTLSCerts[i]

		//Create a lightning client to connect to the node
		lightningClient, conn, err := rpc.CreateLightningClient(nodeEndpoint, nodeTLSCertEncoded)
		if err != nil {
			log.Fatalf("failed to create lightning client: %v", err)
		}
		defer conn.Close()

		//TODO Add support for multiple providers in the future

		//Create SwapClient to communicate with loopd
		loopdHost := loopdHosts[i]

		swapClient, swapConn, err := rpc.CreateSwapClientClient(loopdHost, loopdTLSCertEncoded)
		if err != nil {
			log.Fatalf("failed to create swap client: %v", err)
		}
		defer swapConn.Close()

		//Create NodeGuardClient to communicate with nodeguard
		nodeguardClient, nodeguardConn, err := rpc.CreateNodeGuardClient(nodeguardHost)
		if err != nil {
			log.Fatalf("failed to create nodeguard client: %v", err)
		}
		defer nodeguardConn.Close()

		if err != nil {
			log.Fatalf("failed to create lightning client: %v", err)
		}

		//Macaroon of the lnd node

		nodeMacaroon := nodesMacaroons[i]

		if nodeMacaroon == "" {
			log.Fatalf("no macaroon provided for node %v", nodeEndpoint)
		}

		nodeContext, err := helper.GenerateContextWithMacaroon(nodeMacaroon)
		if err != nil {
			log.Fatal("failed to generate context with macaroon")
		}

		//Macaroon of the loopd of this lnd node
		//TODO Make this optional when loop is not the provider
		loopdMacaroon := loopdMacaroons[i]

		if loopdMacaroon == "" {
			log.Fatalf("no macaroon provided for loopd %v", loopdHost)
		}

		//Get the local node info
		nodeInfo, err := getLocalNodeInfo(&lightningClient, nodeContext)
		if err != nil {
			log.Fatalf("failed to get local node info: %v", err)
		}

		wg.Add(1)

		//Start a goroutine to poll nodeguard for liquidation rules for this node
		go startNodeGuardPolling(*nodeInfo, nodeguardClient, nodeContext)

		//Start a goroutine to monitor the channels of the node
		go monitorChannels(nodeEndpoint, *nodeInfo, nodeMacaroon, loopdMacaroon, lightningClient, nodeguardClient, swapClient, nodeContext)

	}

	wg.Wait()

	//TODO Graceful shutdown
}

// Start a goroutine to poll nodeguard for liquidation rules
func startNodeGuardPolling(nodeInfo lnrpc.GetInfoResponse, nodeguardClient nodeguard.NodeGuardServiceClient, ctx context.Context) {

	pubkey := nodeInfo.GetIdentityPubkey()

	for {

		//Get liquidation rules from nodeguard
		liquidationRules, err := nodeguardClient.GetLiquidityRules(ctx, &nodeguard.GetLiquidityRulesRequest{
			NodePubkey: pubkey,
		})

		if liquidationRules == nil {
			log.Warnf("no liquidation rules found for node %v", pubkey)
			time.Sleep(10 * time.Second)
			continue
		}

		//TODO Maybe there is a better way to do this
		var derefRules []nodeguard.LiquidityRule

		for _, rule := range liquidationRules.LiquidityRules {
			derefRules = append(derefRules, *rule)
		}

		if err != nil {
			log.Errorf("failed to get liquidation rules from nodeguard: %v", err)
		} else {
			//Store liquidation rules in cache
			x := nodeInfo.GetIdentityPubkey()
			rulesCache.SetLiquidityRules(x, derefRules)
		}

		//Sleep for 10 seconds
		time.Sleep(10 * time.Second)

	}

}

// Record the channel balance in a prometheus gauge
func getChannelBalanceRatio(channel *lnrpc.Channel) (float64, error) {

	capacity := float64(channel.GetCapacity())

	if capacity <= 0 {

		err := fmt.Errorf("channel capacity is <= 0")
		log.Error(err)
		return -1, err
	}

	remoteBalance := float64(channel.GetRemoteBalance())

	channelBalanceRatioInt := remoteBalance / capacity

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
func monitorChannels(nodeHost string, nodeInfo lnrpc.GetInfoResponse, nodeMacaroon string, loopdMacaroon string, lightningClient lnrpc.LightningClient, nodeguardClient nodeguard.NodeGuardServiceClient, swapClient looprpc.SwapClientClient, nodeCtx context.Context) {

	//Defer a recover function to catch panics
	defer func() {
		if r := recover(); r != nil {
			log.Error("recovered panic in monitorChannels function, retrying...", r)

			//Sleep for 10 seconds before restarting the monitoring
			time.Sleep(10 * time.Second)

			//Restart monitoring via recursion
			monitorChannels(nodeHost, nodeInfo, nodeMacaroon, loopdMacaroon, lightningClient, nodeguardClient, swapClient, nodeCtx)
		}

	}()

	//Check that node host matches host:port string
	if nodeHost == "" {
		err := fmt.Errorf("nodeHost is empty")
		log.Error(err)
	}

	//Loop provider
	loopProvider := provider.LoopProvider{}
	//Infinite loop to monitor channels
	for {

		//Call ListChannels method of lightning client with metadata headers
		response, err := lightningClient.ListChannels(nodeCtx, &lnrpc.ListChannelsRequest{
			ActiveOnly: false,
		})

		if err != nil {
			log.Errorf("error listing channels: %v", err)
		}

		if response == nil || len(response.Channels) == 0 {
			log.Errorf("no channels found for node %v", nodeHost)

			time.Sleep(1 * time.Second)

			continue
		}

		//TODO Support rules without nodeguard in the future

		//Get liquidation rule from cache
		nodePubKey := nodeInfo.GetIdentityPubkey()
		liquidationRules, err := rulesCache.GetLiquidityRules(nodePubKey)
		if err != nil {
			log.Errorf("failed to get liquidation rules from cache: %v for node %s", err, nodePubKey)
			time.Sleep(1 * time.Second)
			continue
		}

		//Iterate over response channels
		for _, channel := range response.Channels {
			log.Debugf("monitoring Channel Id: %v", channel.ChanId)

			go monitorChannel(channel, nodeHost, lightningClient, nodeCtx, liquidationRules, swapClient, nodeguardClient, loopProvider, loopdMacaroon)

		}
		//Sleep
		time.Sleep(pollingInterval)
	}
}

func monitorChannel(channel *lnrpc.Channel, nodeHost string, lightningClient lnrpc.LightningClient, context context.Context, liquidationRules map[uint64][]nodeguard.LiquidityRule, swapClient looprpc.SwapClientClient, nodeguardClient nodeguard.NodeGuardServiceClient, loopProvider provider.LoopProvider, loopdMacaroon string) {
	//Record the channel balance in a prometheus gauge
	channelBalanceRatio, err := getChannelBalanceRatio(channel)

	if err != nil {
		log.Errorf("error calculating channel balance: %v", err)
	}
	log.Debugf("channel balance ratio for node %v channel %v is %v", nodeHost, channel.GetChanId(), channelBalanceRatio)

	recordChannelBalanceMetric(nodeHost, channel, channelBalanceRatio, lightningClient, context)

	channelRules := liquidationRules[channel.GetChanId()]

	//Manage liquidity
	err = manageChannelLiquidity(channel, channelBalanceRatio, &channelRules, swapClient, nodeguardClient, &loopProvider, loopdMacaroon)

	if err != nil {
		log.Errorf("error managing channel liquidity: %v", err)
	}
}

// Manage channel liquidity and perform swaps if necessary
func manageChannelLiquidity(channel *lnrpc.Channel, channelBalanceRatio float64, channelRules *[]nodeguard.LiquidityRule, swapClientClient looprpc.SwapClientClient, nodeguardClient nodeguard.NodeGuardServiceClient, loopProvider *provider.LoopProvider, loopdMacaroon string) error {

	//Check if channel is active

	//TODO Review if the following checks should return an error or not
	if !channel.GetActive() {
		log.Debugf("channel %v is inactive, skipping", channel.GetChanId())
		return nil
	}

	if len(*channelRules) == 0 {
		log.Debugf("no rules found for channel %v, skipping", channel.GetChanId())
		return nil
	}

	if len(*channelRules) > 1 {
		log.Warnf("multiple rules found for channel %v, only the first rule will be used", channel.GetChanId())
	}

	//TODO Discuss support multiple rules per channel
	rule := (*channelRules)[0]

	//TODO Add some more logic to determine the swap amount properly
	swapTargetRatio := (rule.MinimumRemoteBalance + rule.MinimumLocalBalance) / 2 // Average of the minimum local and remote balance

	var swapTarget int64 = int64(float64(channel.GetCapacity()) * float64(swapTargetRatio))

	loopdCtx, err := helper.GenerateContextWithMacaroon(loopdMacaroon)
	if err != nil {
		return err
	}

	switch {

	case channelBalanceRatio < float64(rule.MinimumLocalBalance):
		{
			//If the balance ratio is below the minimum threshold, perform a reverse swap to increase the local balance

			//Calculate the swap amount
			swapAmount := helper.AbsInt64((channel.LocalBalance - swapTarget))

			//Request nodeguard a new destination address for the reverse swap
			walletRequest := &nodeguard.GetNewWalletAddressRequest{
				WalletId: rule.WalletId,
			}

			addrResponse, err := nodeguardClient.GetNewWalletAddress(context.Background(), walletRequest)
			if err != nil || addrResponse.GetAddress() == "" {
				log.Errorf("error requesting nodeguard a new wallet address: %v", err)
				return err
			}

			//Perform the swap
			swapRequest := provider.ReverseSubmarineSwapRequest{
				SatsAmount:         swapAmount,
				ChannelSet:         []uint64{channel.GetChanId()},
				ReceiverBTCAddress: addrResponse.Address,
			}

			resp, err := loopProvider.RequestReverseSubmarineSwap(loopdCtx, swapRequest, swapClientClient)
			if err != nil {
				log.Errorf("error performing reverse swap: %v", err)
				return err
			}

			log.Infof("reverse swap performed for channel %v, swap id: %v", channel.GetChanId(), resp.SwapId)

			//TODO Monitor the swap status and lock future swaps if the swap is pending until it is completed/failed

		}
	case channelBalanceRatio > float64(rule.MinimumRemoteBalance):
		//The balance ratio is above the maximum threshold, perform a swap to increase the remote balance
		{
			//Calculate the swap amount
			swapAmount := helper.AbsInt64((channel.RemoteBalance - swapTarget))

			//Perform the swap
			swapRequest := provider.SubmarineSwapRequest{
				SatsAmount:    swapAmount,
				LastHopPubkey: channel.RemotePubkey,
			}

			resp, err := loopProvider.RequestSubmarineSwap(loopdCtx, swapRequest, swapClientClient)
			if err != nil {
				log.Errorf("error performing swap: %v", err)
				return err
			}

			if resp.InvoiceBTCAddress == "" {
				err := fmt.Errorf("invoice BTC address is empty for swap id: %v", resp.SwapId)
				log.Errorf("error performing swap: %v", err)
				return err
			}

			//Request nodeguard to send the swap amount to the invoice address

			withdrawalRequest := nodeguard.RequestWithdrawalRequest{
				WalletId:    rule.WalletId,
				Address:     resp.InvoiceBTCAddress,
				Amount:      swapAmount,
				Description: fmt.Sprintf("Swap %v", resp.SwapId),
			}

			withdrawalResponse, err := nodeguardClient.RequestWithdrawal(context.Background(), &withdrawalRequest)
			if err != nil {
				err = fmt.Errorf("error requesting nodeguard to send the swap amount to the invoice address: %v", err)
				log.Errorf("error performing swap: %v", err)

				return err
			}

			if withdrawalResponse.IsHotWallet {
				log.Infof("Swap request sent to nodeguard hot wallet with id: %d for swap id: %v", rule.GetWalletId(), resp.SwapId)
			} else {
				log.Infof("Swap request sent to nodeguard cold wallet with id: %d for swap id: %v", rule.GetWalletId(), resp.SwapId)
			}

			//TODO Monitor the swap status and lock future swaps if the swap is pending until it is completed/failed

		}

	}

	return nil

}

// Record the channel balance in a prometheus gauge
func recordChannelBalanceMetric(nodeHost string, channel *lnrpc.Channel, channelBalanceRatio float64, lightningClient lnrpc.LightningClient, context context.Context) {

	channelId := fmt.Sprint(channel.GetChanId())

	localNodeInfo, err := getInfo(&lightningClient, &context)

	if err != nil {
		log.Errorf("error getting local node info: %v", err)
	}

	remoteNodeInfo, err := getNodeInfo(channel.RemotePubkey, &lightningClient, &context)

	if err != nil {
		log.Errorf("error getting remote node info: %v", err)

	}

	localPubKey := localNodeInfo.GetIdentityPubkey()
	remotePubKey := channel.GetRemotePubkey()
	localAlias := localNodeInfo.GetAlias()
	remoteAlias := remoteNodeInfo.GetNode().GetAlias()

	active := strconv.FormatBool(channel.GetActive())
	initiator := strconv.FormatBool(channel.GetInitiator())

	prometheusMetrics.channelBalanceGauge.With(prometheus.Labels{
		"chan_id":            channelId,
		"local_node_pubkey":  localPubKey,
		"remote_node_pubkey": remotePubKey,
		"local_node_alias":   localAlias,
		"remote_node_alias":  remoteAlias,
		"active":             active,
		"initiator":          initiator,
	}).Set(channelBalanceRatio)
}

// Gets the info from the node which we have the macaroon
func getInfo(lightningClient *lnrpc.LightningClient, context *context.Context) (*lnrpc.GetInfoResponse, error) {

	//Call GetInfo method of lightning client
	response, err := (*lightningClient).GetInfo(*context, &lnrpc.GetInfoRequest{})

	if err != nil {
		log.Errorf("error getting info: %v", err)
		return nil, err
	}

	return response, nil
}

// Gets the info from the node which we have the macaroon
func getLocalNodeInfo(lightningClient *lnrpc.LightningClient, context context.Context) (*lnrpc.GetInfoResponse, error) {

	//Call GetInfo method of lightning client
	response, err := (*lightningClient).GetInfo(context, &lnrpc.GetInfoRequest{})

	if err != nil {
		log.Errorf("error getting info: %v", err)
		return nil, err
	}

	return response, nil
}

// Gets the info from the node with the given pubkey
func getNodeInfo(pubkey string, lightningClient *lnrpc.LightningClient, context *context.Context) (*lnrpc.NodeInfo, error) {

	if pubkey == "" {
		err := fmt.Errorf("pubkey is empty")
		log.Error(err, "pubkey is empty")
		return nil, err

	}

	//Call GetNodeInfo method of lightning client with metadata headers
	response, err := (*lightningClient).GetNodeInfo(*context, &lnrpc.NodeInfoRequest{
		PubKey: pubkey,
	})

	if err != nil {
		log.Errorf("error getting node info: %v", err)
		return nil, err
	}

	return response, nil
}
