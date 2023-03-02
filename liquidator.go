package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Elenpay/liquidator/cache"
	"github.com/Elenpay/liquidator/errors"
	"github.com/Elenpay/liquidator/helper"
	"github.com/Elenpay/liquidator/nodeguard"
	"github.com/Elenpay/liquidator/provider"
	"github.com/Elenpay/liquidator/rpc"
	"github.com/lightninglabs/loop/looprpc"
	"github.com/lightningnetwork/lnd/lnrpc"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var (
	prometheusMetrics *metrics
	rulesCache        cache.Cache
)

// Entrypoint of liquidator main logic
func startLiquidator() {

	//Init opentelemetry tracer
	initTracer(context.TODO())

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

		nodeContext, err := helper.GenerateContextWithMacaroon(nodeMacaroon, context.Background())
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
		nodeInfo, err := getLocalNodeInfo(lightningClient, nodeContext)
		if err != nil {
			log.Fatalf("failed to get local node info: %v", err)
		}

		wg.Add(1)

		//Start a goroutine to poll nodeguard for liquidation rules for this node
		go startNodeGuardPolling(nodeInfo, nodeguardClient, nodeContext)

		//Start a goroutine to monitor the channels of the node
		go monitorChannels(MonitorChannelsInfo{
			nodeHost:        nodeEndpoint,
			nodeInfo:        nodeInfo,
			nodeMacaroon:    nodeMacaroon,
			loopdMacaroon:   loopdMacaroon,
			lightningClient: lightningClient,
			nodeguardClient: nodeguardClient,
			swapClient:      swapClient,
			nodeCtx:         nodeContext,
		})

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

		if liquidationRules == nil || len(liquidationRules.LiquidityRules) == 0 {
			log.Debugf("no liquidation rules found for node %v", pubkey)
			time.Sleep(pollingInterval)
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
		time.Sleep(pollingInterval)

	}

}

// Calculate the ratio of remote balance to capacity of a channel
func getChannelBalanceRatio(channel *lnrpc.Channel, spanCtx context.Context) (float64, error) {

	//Start span
	spanCtx, span := otel.Tracer("monitorChannel").Start(spanCtx, "getChannelBalanceRatio")
	defer span.End()

	capacity := float64(channel.GetCapacity())

	if capacity <= 0 {

		err := fmt.Errorf("channel capacity is <= 0")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.WithField("span", span).Error(err)
		return -1, err
	}

	remoteBalance := float64(channel.GetRemoteBalance())

	channelBalanceRatioInt := remoteBalance / capacity

	//Truncate channelbalance to 2 decimal places
	channelBalanceRatio := float64(int(channelBalanceRatioInt*100)) / 100

	//Check that the ratio is between 0 and 1
	if channelBalanceRatio > 1 || channelBalanceRatio < 0 {

		err := fmt.Errorf("channel balance ratio is not between 0 and 1")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.WithField("span", span).Error(err)
		return -1, err
	}

	return channelBalanceRatio, nil

}

// Locking fuction to be used in a goroutine to monitor channels
func monitorChannels(info MonitorChannelsInfo) {

	//Defer a recover function to catch panics
	defer func() {
		if r := recover(); r != nil {
			log.Error("recovered panic in monitorChannels function, retrying...", r)

			//Sleep for 10 seconds before restarting the monitoring
			time.Sleep(10 * time.Second)

			//Restart monitoring via recursion
			monitorChannels(info)
		}

	}()

	//Check that node host matches host:port string
	if info.nodeHost == "" {
		err := fmt.Errorf("nodeHost is empty")
		log.Error(err)
	}

	//Loop provider
	loopProvider := provider.LoopProvider{}
	//Infinite loop to monitor channels
	for {

		//Call ListChannels method of lightning client with metadata headers
		response, err := info.lightningClient.ListChannels(info.nodeCtx, &lnrpc.ListChannelsRequest{
			ActiveOnly: false,
		})

		if err != nil {
			log.Errorf("error listing channels: %v", err)
		}

		if response == nil || len(response.Channels) == 0 {
			log.Debugf("no channels found for node %v", info.nodeHost)

			time.Sleep(pollingInterval)

			continue
		}

		//TODO Support rules without nodeguard in the future

		//Get liquidation rule from cache
		nodePubKey := info.nodeInfo.GetIdentityPubkey()
		liquidationRules, err := rulesCache.GetLiquidityRules(nodePubKey)
		if err != nil {
			log.Debugf("failed to get liquidation rules from cache: %v for node %s", err, nodePubKey)
			time.Sleep(pollingInterval)
			continue
		}

		//Iterate over response channels
		for _, channel := range response.Channels {
			log.Debugf("monitoring Channel Id: %v", channel.ChanId)

			go monitorChannel(MonitorChannelInfo{
				channel:          channel,
				nodeHost:         info.nodeHost,
				lightningClient:  info.lightningClient,
				context:          info.nodeCtx,
				liquidationRules: liquidationRules,
				swapClient:       info.swapClient,
				nodeguardClient:  info.nodeguardClient,
				loopProvider:     &loopProvider,
				loopdMacaroon:    info.loopdMacaroon,
				nodeInfo:         info.nodeInfo,
			})

		}
		//Sleep
		time.Sleep(pollingInterval)
	}
}

// Monitor a single channel liquidity and perform actions if needed
func monitorChannel(info MonitorChannelInfo) {

	//Start span
	spanCtx, span := otel.Tracer("monitorChannel").Start(info.context, "monitorChannel")
	defer span.End()

	//Add atrributes to span
	span.SetAttributes(attribute.String("nodeHost", info.nodeHost), attribute.Int64("channelId", int64(info.channel.GetChanId())))

	spanId := span.SpanContext().SpanID().String()
	traceId := span.SpanContext().TraceID().String()

	log.WithField("span", span).Debugf("monitorChannel SpanId: %v TraceId: %v", spanId, traceId)
	//Record the channel balance in a prometheus gauge
	channelBalanceRatio, err := getChannelBalanceRatio(info.channel, spanCtx)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.WithField("span", span).Errorf("error calculating channel balance: %v", err)
	}
	log.Debugf("channel balance ratio for node %v channel %v is %v", info.nodeHost, info.channel.GetChanId(), channelBalanceRatio)

	recordChannelBalanceMetric(info.nodeHost, info.channel, channelBalanceRatio, info.lightningClient, spanCtx)

	channelRules := info.liquidationRules[info.channel.GetChanId()]

	//Manage liquidity
	err = manageChannelLiquidity(ManageChannelLiquidityInfo{
		channel:             info.channel,
		channelBalanceRatio: channelBalanceRatio,
		channelRules:        &channelRules,
		swapClientClient:    info.swapClient,
		nodeguardClient:     info.nodeguardClient,
		loopProvider:        info.loopProvider,
		loopdMacaroon:       info.loopdMacaroon,
		nodeInfo:            info.nodeInfo,
		ctx:                 spanCtx,
	})

	if err != nil {

		//if err is not of type SwapInProgressError, record it

		switch err.(type) {
		case *errors.SwapInProgressError:

		default:
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		log.WithField("span", span).Errorf("error managing channel liquidity: %v", err)
	}
}

// Manage channel liquidity and perform swaps if necessary
func manageChannelLiquidity(info ManageChannelLiquidityInfo) error {

	//Start span
	_, span := otel.Tracer("monitorChannel").Start(info.ctx, "manageChannelLiquidity")
	defer span.End()

	//Check if channel is active
	channel := info.channel
	channelRules := info.channelRules

	//TODO Review if the following checks should return an error or not
	if !channel.GetActive() {
		log.WithField("span", span).Debugf("channel %v is inactive, skipping", channel.GetChanId())
		return nil
	}

	if len(*channelRules) == 0 {
		log.WithField("span", span).Debugf("no rules found for channel %v, skipping", channel.GetChanId())
		return nil
	}

	if len(*channelRules) > 1 {
		log.WithField("span", span).Warnf("multiple rules found for channel %v, only the first rule will be used", channel.GetChanId())
	}

	//TODO Discuss support multiple rules per channel
	rule := (*channelRules)[0]

	//Divide all rule percents by 100 to get ratios
	rule.MinimumLocalBalance = rule.MinimumLocalBalance / 100
	rule.MinimumRemoteBalance = rule.MinimumRemoteBalance / 100
	rule.RebalanceTarget = rule.RebalanceTarget / 100

	// If rebalance target is 0, the swap target balance is the average of the minimum local and remote balance

	swapTargetRatio := (rule.MinimumRemoteBalance + rule.MinimumLocalBalance) / 2 // Average of the minimum local and remote balance

	// If rebalance target is > 0, the swap target balance is the rebalance target
	if rule.RebalanceTarget > 0 {
		swapTargetRatio = rule.RebalanceTarget
	}

	var swapAmountTarget int64 = int64(float64(channel.GetCapacity()) * float64(swapTargetRatio))

	loopdCtx, err := helper.GenerateContextWithMacaroon(info.loopdMacaroon, info.ctx)
	if err != nil {
		return err
	}

	switch {

	case info.channelBalanceRatio < float64(rule.MinimumLocalBalance):
		{
			//If the balance ratio is below the minimum threshold, perform a reverse swap to increase the local balance

			//Add attribute to the span of swap requested
			span.SetAttributes(attribute.String("swapRequestedType", "reverse"))
			span.SetAttributes(attribute.String("chanId", fmt.Sprintf("%v", channel.GetChanId())))
			span.SetAttributes(attribute.String("nodePubkey", info.nodeInfo.GetIdentityPubkey()))
			span.SetAttributes(attribute.String("nodeAlias", info.nodeInfo.GetAlias()))

			//Calculate the swap amount
			swapAmount := helper.AbsInt64((channel.RemoteBalance - swapAmountTarget))

			//Request nodeguard a new destination address for the reverse swap
			walletRequest := &nodeguard.GetNewWalletAddressRequest{
				WalletId: rule.WalletId,
			}

			addrResponse, err := info.nodeguardClient.GetNewWalletAddress(info.ctx, walletRequest)
			if err != nil || addrResponse.GetAddress() == "" {
				log.WithField("span", span).Errorf("error requesting nodeguard a new wallet address: %v on node: %v", err, info.nodeInfo.GetIdentityPubkey())
				return err
			}

			//Perform the swap
			swapRequest := provider.ReverseSubmarineSwapRequest{
				SatsAmount:         swapAmount,
				ChannelSet:         []uint64{channel.GetChanId()},
				ReceiverBTCAddress: addrResponse.Address,
			}

			resp, err := info.loopProvider.RequestReverseSubmarineSwap(loopdCtx, swapRequest, info.swapClientClient)
			if err != nil {
				log.WithField("span", span).Errorf("error performing reverse swap: %v on node: %v", err, info.nodeInfo.GetIdentityPubkey())
				return err
			}

			//Monitor the swap
			swapStatus, err := info.loopProvider.MonitorSwap(loopdCtx, resp.SwapId, info.swapClientClient)
			if err != nil {
				return err
			}

			if swapStatus.State == looprpc.SwapState_FAILED {
				//Error log: The swap was failed
				err := fmt.Errorf("failed reverse swap, failure reason: %v channel: %v on node: %v", swapStatus.GetFailureReason(), channel.GetChanId(), info.nodeInfo.GetIdentityPubkey())
				log.WithField("span", span).Error(err)
				return err
			} else if swapStatus.State == looprpc.SwapState_SUCCESS {
				//Success log: The swap was successful
				log.WithField("span", span).Infof("reverse swap performed for channel %v, swap id: %v on node %v", channel.GetChanId(), resp.SwapId, info.nodeInfo.GetIdentityPubkey())

				swapType := "rswap"

				//Add fees to counter
				recordSwapFees(swapStatus, info, swapType, channel)

			}

		}
	case info.channelBalanceRatio > float64(rule.MinimumRemoteBalance):
		{
			//The balance ratio is above the maximum threshold, perform a swap to increase the remote balance

			//Add attribute to the span of swap requested
			span.SetAttributes(attribute.String("swapRequestedType", "swap"))
			span.SetAttributes(attribute.String("chanId", fmt.Sprintf("%v", channel.GetChanId())))
			span.SetAttributes(attribute.String("nodePubkey", info.nodeInfo.GetIdentityPubkey()))
			span.SetAttributes(attribute.String("nodeAlias", info.nodeInfo.GetAlias()))
			//Calculate the swap amount
			swapAmount := helper.AbsInt64((channel.RemoteBalance - swapAmountTarget))

			//Perform the swap
			swapRequest := provider.SubmarineSwapRequest{
				SatsAmount:    swapAmount,
				LastHopPubkey: channel.RemotePubkey,
			}

			resp, err := info.loopProvider.RequestSubmarineSwap(loopdCtx, swapRequest, info.swapClientClient)
			if err != nil {
				log.WithField("span", span).Errorf("error performing swap: %v on node: %v", err, info.nodeInfo.GetIdentityPubkey())
				return err
			}

			if resp.InvoiceBTCAddress == "" {
				err := fmt.Errorf("invoice BTC address is empty for swap id: %v on node: %v", resp.SwapId, info.nodeInfo.GetIdentityPubkey())
				log.WithField("span", span).Errorf("error performing swap: %v", err)
				return err
			}

			//Request nodeguard to send the swap amount to the invoice address

			withdrawalRequest := nodeguard.RequestWithdrawalRequest{
				WalletId:    rule.WalletId,
				Address:     resp.InvoiceBTCAddress,
				Amount:      swapAmount,
				Description: fmt.Sprintf("Swap %v", resp.SwapId),
			}

			withdrawalResponse, err := info.nodeguardClient.RequestWithdrawal(info.ctx, &withdrawalRequest)
			if err != nil {
				err = fmt.Errorf("error requesting nodeguard to send the swap amount to the invoice address: %v on node: %v", err, info.nodeInfo.GetIdentityPubkey())
				log.WithField("span", span).Errorf("error performing swap: %v", err)

				return err
			}

			//Log the swap request

			if withdrawalResponse.IsHotWallet {
				log.WithField("span", span).Infof("Swap request sent to nodeguard hot wallet with id: %d for swap id: %v for node: %v", rule.GetWalletId(), resp.SwapId, info.nodeInfo.GetIdentityPubkey())
			} else {
				log.WithField("span", span).Infof("Swap request sent to nodeguard cold wallet with id: %d for swap id: %v for node: %v ", rule.GetWalletId(), resp.SwapId, info.nodeInfo.GetIdentityPubkey())
			}

			//Monitor the swap
			swapStatus, err := info.loopProvider.MonitorSwap(loopdCtx, resp.SwapId, info.swapClientClient)
			if err != nil {
				return err
			}

			if swapStatus.State == looprpc.SwapState_FAILED {
				//Error log: The swap was failed
				err := fmt.Errorf("failed swap failure reason: %v, channel: %v on node: %v", swapStatus.GetFailureReason(), channel.GetChanId(), info.nodeInfo.GetIdentityPubkey())
				log.WithField("span", span).Error(err)
				return err
			} else if swapStatus.State == looprpc.SwapState_SUCCESS {
				//Success log: The swap was successful
				log.WithField("span", span).Infof("swap performed for channel %v, swap id: %v on node %v", channel.GetChanId(), resp.SwapId, info.nodeInfo.GetIdentityPubkey())

				swapType := "swap"

				//Add fees to counter
				recordSwapFees(swapStatus, info, swapType, channel)
			}

		}

	}

	return nil

}

// Add fees to the prometheus counter for swaps
func recordSwapFees(swapStatus looprpc.SwapStatus, info ManageChannelLiquidityInfo, swapType string, channel *lnrpc.Channel) {

	offChainFees := float64(swapStatus.GetCostOffchain())
	onChainFees := float64(swapStatus.GetCostOnchain())
	providerFees := float64(swapStatus.GetCostServer())

	if offChainFees > 0 {

		prometheusMetrics.offchainFees.With(prometheus.Labels{"node_pubkey": info.nodeInfo.GetIdentityPubkey(), "node_alias": info.nodeInfo.GetAlias(), "swap_type": swapType, "chan_id": fmt.Sprintf("%v", channel.ChanId), "provider": "loop"}).Add(offChainFees)

	}

	if onChainFees > 0 {
		prometheusMetrics.onchainFees.With(prometheus.Labels{"node_pubkey": info.nodeInfo.GetIdentityPubkey(), "node_alias": info.nodeInfo.GetAlias(), "swap_type": swapType, "chan_id": fmt.Sprintf("%v", channel.ChanId), "provider": "loop"}).Add(onChainFees)
	}

	if providerFees > 0 {

		prometheusMetrics.providerFees.With(prometheus.Labels{"node_pubkey": info.nodeInfo.GetIdentityPubkey(), "node_alias": info.nodeInfo.GetAlias(), "swap_type": swapType, "chan_id": fmt.Sprintf("%v", channel.ChanId), "provider": "loop"}).Add(providerFees)
	}

}

// Record the channel balance in a prometheus gauge
func recordChannelBalanceMetric(nodeHost string, channel *lnrpc.Channel, channelBalanceRatio float64, lightningClient lnrpc.LightningClient, context context.Context) {
	//Start span
	_, span := otel.Tracer("monitorChannel").Start(context, "recordChannelBalanceMetric")
	defer span.End()

	channelId := fmt.Sprint(channel.GetChanId())

	localNodeInfo, err := getInfo(lightningClient, context)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		log.WithField("span", span).Errorf("error getting local node info: %v", err)
	}

	remoteNodeInfo, err := getNodeInfo(channel.RemotePubkey, lightningClient, context)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		log.WithField("span", span).Errorf("error getting remote node info: %v", err)

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
func getInfo(lightningClient lnrpc.LightningClient, context context.Context) (lnrpc.GetInfoResponse, error) {

	//Call GetInfo method of lightning client
	response, err := lightningClient.GetInfo(context, &lnrpc.GetInfoRequest{})

	if err != nil {
		log.Errorf("error getting info: %v", err)
		return lnrpc.GetInfoResponse{}, err
	}

	return *response, nil
}

// Gets the info from the node which we have the macaroon
func getLocalNodeInfo(lightningClient lnrpc.LightningClient, context context.Context) (lnrpc.GetInfoResponse, error) {

	//Call GetInfo method of lightning client
	response, err := lightningClient.GetInfo(context, &lnrpc.GetInfoRequest{})

	if err != nil {
		log.Errorf("error getting info: %v", err)
		return lnrpc.GetInfoResponse{}, err
	}

	return *response, nil
}

// Gets the info from the node with the given pubkey
func getNodeInfo(pubkey string, lightningClient lnrpc.LightningClient, context context.Context) (*lnrpc.NodeInfo, error) {

	if pubkey == "" {
		err := fmt.Errorf("pubkey is empty")
		log.Error(err, "pubkey is empty")
		return nil, err

	}

	//Call GetNodeInfo method of lightning client with metadata headers
	response, err := (lightningClient).GetNodeInfo(context, &lnrpc.NodeInfoRequest{
		PubKey: pubkey,
	})

	if err != nil {
		log.Errorf("error getting node info: %v", err)
		return nil, err
	}

	return response, nil
}
