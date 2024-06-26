package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Elenpay/liquidator/cache"
	"github.com/Elenpay/liquidator/customerrors"
	"github.com/Elenpay/liquidator/helper"
	"github.com/Elenpay/liquidator/lndconnect"
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
	"go.opentelemetry.io/otel/trace"
)

var (
	prometheusMetrics  *metrics
	rulesCache         cache.Cache
	retries            int
	backoffCoefficient float64
	backoffLimit       float64
)

// Entrypoint of liquidator main logic
func startLiquidator() {
	//Init opentelemetry tracer
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {

		initTracer(context.TODO())
	}

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

	lightningClients := make(map[string]lnrpc.LightningClient)
	nodeCtxs := make(map[string]context.Context)

	for i, lndconnectURI := range lndconnectURIs {

		//parse lndconnectURI
		lndConnectParams, err := lndconnect.Parse(lndconnectURI)
		if err != nil {
			log.Fatalf("failed to parse lndconnectURI: %v", err)
		}

		//Macaroon of the loopd of this lnd node
		//TODO Make this optional when loop is not the provider
		loopdConnectParams, err := lndconnect.Parse(loopdconnectURIs[i])
		if err != nil {
			log.Fatalf("failed to parse loopdconnectURI: %v", err)
		}

		loopdHost := fmt.Sprintf("%v:%v", loopdConnectParams.Host, loopdConnectParams.Port)
		loopdMacaroon := loopdConnectParams.Macaroon

		if loopdMacaroon == "" {
			log.Fatalf("no macaroon provided for loopd %v", loopdHost)
		}

		log.Infof("starting monitoring for node: %v:%v loopd on %v", lndConnectParams.Host, lndConnectParams.Port, loopdHost)

		//Create a lightning client to connect to the node
		lightningClient, conn, err := rpc.CreateLightningClient(lndConnectParams)
		if err != nil {
			log.Fatalf("failed to create lightning client: %v", err)
		}
		defer conn.Close()

		//Create SwapClient to communicate with loopd

		swapClient, swapConn, err := rpc.CreateSwapClientClient(loopdConnectParams)
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

		nodeMacaroon := lndConnectParams.Macaroon

		if nodeMacaroon == "" {
			log.Fatalf("no macaroon provided for node %v", lndconnectURI)
		}

		nodeCtx, err := helper.GenerateContextWithMacaroon(nodeMacaroon, context.Background())
		if err != nil {
			log.Fatal("failed to generate context with macaroon")
		}

		//Get the local node info
		nodeInfo, err := getLocalNodeInfo(lightningClient, nodeCtx)
		if err != nil {
			log.Fatalf("failed to get local node info: %v", err)
		}

		nodeCtxs[nodeInfo.IdentityPubkey] = nodeCtx

		//Store the lightning client in a map so that it can be used later
		lightningClients[nodeInfo.IdentityPubkey] = lightningClient

		wg.Add(1)

		//Start a goroutine to poll nodeguard for liquidation rules for this node
		go startNodeGuardPolling(nodeInfo, nodeguardClient, nodeCtx)

		//Start a goroutine to monitor the channels of the node
		lndconnectParams, err := lndconnect.Parse(lndconnectURI)

		if err != nil {
			log.Fatalf("failed to parse lndconnectURI: %v", err)
		}

		go monitorChannels(MonitorChannelsInfo{
			BaseInfo: BaseInfo{
				nodeHost:         lndconnectParams.Host + ":" + lndconnectParams.Port,
				nodeInfo:         nodeInfo,
				nodeMacaroon:     nodeMacaroon,
				loopdMacaroon:    loopdMacaroon,
				lightningClients: lightningClients,
				nodeguardClient:  nodeguardClient,
				swapClient:       swapClient,
				nodeCtxs:         nodeCtxs,
				provider:         &provider.LoopProvider{},
			},
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

		if err != nil {
			log.Errorf("startNodeguardPolling: didn't get liquitadtionRules due to the following error: %v", err)
		}

		if liquidationRules == nil || len(liquidationRules.LiquidityRules) == 0 {
			log.Debugf("startNodeguardPolling: no liquidation rules found for node %v, retrying in %s...", pubkey, pollingInterval.String())
			time.Sleep(pollingInterval)
			continue
		}

		//TODO Maybe there is a better way to do this
		var derefRules []nodeguard.LiquidityRule

		for _, rule := range liquidationRules.LiquidityRules {
			derefRules = append(derefRules, *rule)
		}

		//Store liquidation rules in cache
		x := nodeInfo.GetIdentityPubkey()
		rulesCache.SetLiquidityRules(x, derefRules)

		//Sleep for 10 seconds
		time.Sleep(pollingInterval)

	}

}

// Calculate the ratio of remote balance to capacity of a channel
func getChannelBalanceRatio(channel *lnrpc.Channel, spanCtx context.Context) (float64, error) {

	//Start span
	_, span := otel.Tracer("monitorChannel").Start(spanCtx, "getChannelBalanceRatio")
	defer span.End()

	capacity := float64(channel.GetLocalBalance() + channel.GetRemoteBalance() + channel.GetUnsettledBalance())

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

	prevChannels := []*lnrpc.Channel{}

	//Infinite loop to monitor channels
	for {
		//Call ListChannels method of lightning client with metadata headers
		response, err := info.lightningClients[info.nodeInfo.IdentityPubkey].ListChannels(info.nodeCtxs[info.nodeInfo.IdentityPubkey], &lnrpc.ListChannelsRequest{
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

		// remove channel metrics that are not in list of channels anymore
		for _, channel := range prevChannels {
			found := false
			for _, newChannel := range response.Channels {
				if channel.GetChanId() == newChannel.GetChanId() {
					found = true
					break
				}
			}
			if !found {
				//Remove channel balance metric
				deleteChannelBalanceMetric(info.nodeHost, channel, info.lightningClients[info.nodeInfo.IdentityPubkey], info.nodeCtxs[info.nodeInfo.IdentityPubkey])
			}
		}

		prevChannels = response.Channels
		//TODO Support rules without nodeguard in the future

		//Get liquidation rule from cache
		nodePubKey := info.nodeInfo.GetIdentityPubkey()
		liquidationRules, err := rulesCache.GetLiquidityRules(nodePubKey)
		if err != nil {
			log.Debugf("failed to get liquidation rules from cache: %v for node %s", err, nodePubKey)
		}

		//Iterate over response channels
		for _, channel := range response.Channels {
			log.Debugf("monitoring Channel Id: %v", channel.ChanId)

			go monitorChannel(MonitorChannelInfo{
				BaseInfo:         info.BaseInfo,
				channel:          channel,
				context:          info.nodeCtxs[info.nodeInfo.IdentityPubkey],
				liquidationRules: liquidationRules,
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

	//Record the channel balance in a prometheus gauge
	channelBalanceRatio, err := getChannelBalanceRatio(info.channel, spanCtx)
	//Add atrributes to span
	span.SetAttributes(attribute.String("nodeHost", info.nodeHost), attribute.String("chanId", fmt.Sprintf("%v", info.channel.GetChanId())), attribute.Float64("channelBalanceRatio", channelBalanceRatio))

	spanId := span.SpanContext().SpanID().String()
	traceId := span.SpanContext().TraceID().String()

	log.WithField("span", span).Debugf("monitorChannel SpanId: %v TraceId: %v", spanId, traceId)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.WithField("span", span).Errorf("error calculating channel balance: %v", err)
	}
	log.Debugf("channel balance ratio for node %v channel %v is %v", info.nodeHost, info.channel.GetChanId(), channelBalanceRatio)

	recordChannelBalanceMetric(info.nodeHost, info.channel, channelBalanceRatio, info.lightningClients[info.nodeInfo.IdentityPubkey], spanCtx)

	if info.liquidationRules == nil {
		return
	}

	channelRules := info.liquidationRules[info.channel.GetChanId()]

	//Manage liquidity if there are no pending htlcs
	err = manageChannelLiquidity(ManageChannelLiquidityInfo{
		BaseInfo:            info.BaseInfo,
		channel:             info.channel,
		channelBalanceRatio: channelBalanceRatio,
		channelRules:        &channelRules,
		ctx:                 spanCtx,
	})

	if err != nil {

		//if err is not of type SwapInProgressError, record it

		switch err.(type) {
		case *customerrors.SwapInProgressError:

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

	channelRealCapacity := channel.GetLocalBalance() + channel.GetRemoteBalance() + channel.GetUnsettledBalance()
	var swapAmountTarget int64 = int64(float64(channelRealCapacity) * float64(swapTargetRatio))

	loopdCtx, err := helper.GenerateContextWithMacaroon(info.loopdMacaroon, info.ctx)
	if err != nil {
		return err
	}

	switch {
	//Both nodes are managed, then simply create 1 invoice and send it to the other node to pay it
	case info.lightningClients[rule.NodePubkey] != nil && info.lightningClients[rule.RemoteNodePubkey] != nil && rule.MinimumLocalBalance != 0 && info.channelBalanceRatio < float64(rule.MinimumLocalBalance):
		{
			//Add attribute to the span of swap requested
			span.SetAttributes(attribute.String("swapRequestedType", "invoiceRebalance"))
			span.SetAttributes(attribute.String("chanId", fmt.Sprintf("%v", channel.GetChanId())))
			span.SetAttributes(attribute.String("nodePubkey", info.nodeInfo.GetIdentityPubkey()))
			span.SetAttributes(attribute.String("nodeAlias", info.nodeInfo.GetAlias()))

			log.WithField("span", span).Infof("rebalancing via invoice on channel %v on node %v", channel.GetChanId(), info.nodeInfo.GetAlias())

			//Create an invoice for the swap amount from the remote node and pay with the rule's node

			swapAmount := helper.AbsInt64((swapAmountTarget - channel.LocalBalance))

			err := invoiceRebalance(info, swapAmount, rule.NodePubkey, rule.RemoteNodePubkey)
			if err != nil {
				return err
			}

		}
	case info.lightningClients[rule.NodePubkey] != nil && info.lightningClients[rule.RemoteNodePubkey] != nil && rule.MinimumRemoteBalance != 0 && info.channelBalanceRatio > float64(rule.MinimumRemoteBalance):
		{
			//Add attribute to the span of swap requested
			span.SetAttributes(attribute.String("swapRequestedType", "invoiceRebalance"))
			span.SetAttributes(attribute.String("chanId", fmt.Sprintf("%v", channel.GetChanId())))
			span.SetAttributes(attribute.String("nodePubkey", info.nodeInfo.GetIdentityPubkey()))
			span.SetAttributes(attribute.String("nodeAlias", info.nodeInfo.GetAlias()))

			log.WithField("span", span).Infof("rebalancing via invoice on channel %v on node %v", channel.GetChanId(), info.nodeInfo.GetAlias())
			swapAmount := helper.AbsInt64((channel.RemoteBalance - swapAmountTarget))

			//Create an invoice for the swap amount from the rule's node and pay with the remote node
			err := invoiceRebalance(info, swapAmount, rule.RemoteNodePubkey, rule.NodePubkey)
			if err != nil {
				return err
			}
		}
	case rule.MinimumLocalBalance != 0 && info.channelBalanceRatio < float64(rule.MinimumLocalBalance):
		{
			swapAmount := helper.AbsInt64((swapAmountTarget - channel.LocalBalance))
			//If the balance ratio is below the minimum threshold, perform a reverse swap to increase the local balance

			//Add attribute to the span of swap requested
			span.SetAttributes(attribute.String("swapRequestedType", "reverse"))
			span.SetAttributes(attribute.String("chanId", fmt.Sprintf("%v", channel.GetChanId())))
			span.SetAttributes(attribute.String("nodePubkey", info.nodeInfo.GetIdentityPubkey()))
			span.SetAttributes(attribute.String("nodeAlias", info.nodeInfo.GetAlias()))

			retryCounter := 1

			err := performReverseSwap(info, channel, swapAmount, rule, span, loopdCtx, retryCounter, swapAmountTarget)
			if err != nil {
				return err
			}
		}
	case rule.MinimumRemoteBalance != 0 && info.channelBalanceRatio > float64(rule.MinimumRemoteBalance):
		{
			swapAmount := helper.AbsInt64((channel.RemoteBalance - swapAmountTarget))
			//The balance ratio is above the maximum threshold, perform a swap to increase the remote balance

			//Add attribute to the span of swap requested
			span.SetAttributes(attribute.String("swapRequestedType", "swap"))
			span.SetAttributes(attribute.String("chanId", fmt.Sprintf("%v", channel.GetChanId())))
			span.SetAttributes(attribute.String("nodePubkey", info.nodeInfo.GetIdentityPubkey()))
			span.SetAttributes(attribute.String("nodeAlias", info.nodeInfo.GetAlias()))

			retryCounter := 1

			err := performSwap(info, channel, swapAmount, rule, span, loopdCtx, retryCounter, swapAmountTarget)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

// Creates one invoice in payee node and pays it with the payer node
func invoiceRebalance(info ManageChannelLiquidityInfo, swapAmount int64, payerPubKey string, payeePubKey string) error {

	if swapAmount <= 0 {
		return errors.New("swap amount is <= 0")
	}

	//Create an invoice for the swap amount
	invoiceRequest := &lnrpc.Invoice{
		Memo:  fmt.Sprintf("Invoice rebalance for channel %v on date %v", info.channel.GetChanId(), time.Now().Format("2006-01-02 15:04:05")),
		Value: swapAmount,
	}

	payeeLnClient := info.lightningClients[payeePubKey]
	payerLnClient := info.lightningClients[payerPubKey]

	invoiceResponse, err := payeeLnClient.AddInvoice(info.nodeCtxs[payeePubKey], invoiceRequest)
	if err != nil {
		return err
	}

	sendResponse, err := payerLnClient.SendPaymentSync(info.nodeCtxs[payerPubKey], &lnrpc.SendRequest{
		PaymentRequest: invoiceResponse.PaymentRequest,
	})
	if err != nil {
		return err
	}

	if sendResponse.PaymentError != "" {
		return errors.New(sendResponse.PaymentError)
	}

	log.Infof("invoice rebalance successful for channel %v on node %v", info.channel.GetChanId(), info.nodeInfo.GetIdentityPubkey())

	return nil

}

func performSwap(info ManageChannelLiquidityInfo, channel *lnrpc.Channel, swapAmount int64, rule nodeguard.LiquidityRule, span trace.Span, loopdCtx context.Context, retryCounter int, swapAmountTarget int64) error {
	//Perform the swap
	swapRequest := provider.SubmarineSwapRequest{
		SatsAmount:    swapAmount,
		LastHopPubkey: channel.RemotePubkey,
	}

	//Log including channel.GetChanId() and swapAmount
	log.WithField("span", span).Infof("requesting submarine swap with amount: %d sats to node %s, channel: %v", swapRequest.SatsAmount, info.nodeInfo.GetIdentityPubkey(), channel.GetChanId())

	resp, err := info.provider.RequestSubmarineSwap(loopdCtx, swapRequest, info.swapClient)
	if err != nil {
		log.WithField("span", span).Errorf("error requesting swap: %v on node: %v", err, info.nodeInfo.GetIdentityPubkey())
		return err
	}

	log.WithField("span", span).Infof("submarine swap requested for channel %v, swap id: %v on node %v", channel.GetChanId(), resp.SwapId, info.nodeInfo.GetIdentityPubkey())

	invoiceAddress := resp.InvoiceBTCAddress

	//Retry for 10 times to try to get the invoice btc address if for some reason it was empty
	if invoiceAddress == "" && resp.SwapId != "" {

		//While the swap htlc address is not set, retry 10 times, each time exponentially increasing the time between retries
		for i := 0; i < 10; i++ {
			//Get the swap status

			swapStatus, err := info.provider.GetSwapStatus(loopdCtx, resp.SwapId, info.swapClient)

			if err != nil {
				return err
			}

			if swapStatus.HtlcAddressP2Wsh != "" || swapStatus.HtlcAddressP2Tr != "" {
				//The swap htlc address is set, break the loop

				if swapStatus.HtlcAddressP2Tr != "" {
					invoiceAddress = swapStatus.HtlcAddressP2Tr
				} else {
					invoiceAddress = swapStatus.HtlcAddressP2Wsh
				}

				break
			}

			//Sleep for 2^i seconds
			time.Sleep(time.Duration(math.Pow(2, float64(i))) * time.Second)

		}

	}

	if invoiceAddress == "" {
		err := fmt.Errorf("invoice BTC address is empty for swap id: %v on node: %v", resp.SwapId, info.nodeInfo.GetIdentityPubkey())
		log.WithField("span", span).Errorf("error performing swap: %v", err)
		return err
	}

	//Request nodeguard to send the swap amount to the invoice address

	withdrawalRequest := nodeguard.RequestWithdrawalRequest{
		WalletId:       rule.SwapWalletId,
		Address:        resp.InvoiceBTCAddress,
		Amount:         swapAmount,
		Description:    fmt.Sprintf("Swap %v", resp.SwapId),
		MempoolFeeRate: nodeguard.FEES_TYPE_ECONOMY_FEE,
	}

	withdrawalResponse, err := info.nodeguardClient.RequestWithdrawal(info.ctx, &withdrawalRequest)
	if err != nil {
		err = fmt.Errorf("error requesting nodeguard to send the swap amount to the invoice address: %v on node: %v", err, info.nodeInfo.GetIdentityPubkey())
		log.WithField("span", span).Errorf("error performing swap: %v", err)

		return err
	}

	//Log the swap request

	if withdrawalResponse.IsHotWallet {
		log.WithField("span", span).Infof("Swap request sent to nodeguard hot wallet with id: %d for swap id: %v for node: %v", rule.GetSwapWalletId(), resp.SwapId, info.nodeInfo.GetIdentityPubkey())
	} else {
		log.WithField("span", span).Infof("Swap request sent to nodeguard cold wallet with id: %d for swap id: %v for node: %v ", rule.GetSwapWalletId(), resp.SwapId, info.nodeInfo.GetIdentityPubkey())
	}

	//Monitor the swap
	swapStatus, err := info.provider.MonitorSwap(loopdCtx, resp.SwapId, info.swapClient)
	if err != nil {
		return err
	}

	if swapStatus.State == looprpc.SwapState_FAILED {
		//Error log: The swap was failed
		err := fmt.Errorf("failed swap failure reason: %v, channel: %v on node: %v", swapStatus.GetFailureReason(), channel.GetChanId(), info.nodeInfo.GetIdentityPubkey())
		log.WithField("span", span).Error(err)

		if retryCounter < retries {
			//Retry the swap
			retryCounter++
			log.WithField("span", span).Infof("retrying swap for channel %v, retry number: %v/%v on node %v", channel.GetChanId(), retryCounter, retries, info.nodeInfo.GetIdentityPubkey())
			err := performSwap(info, channel, swapAmount, rule, span, loopdCtx, retryCounter, swapAmountTarget)
			if err != nil {
				return err
			}
		} else {
			limitSwapAmount := float64(helper.AbsInt64((channel.RemoteBalance - swapAmountTarget))) * backoffLimit
			if limitSwapAmount < float64(swapAmount)*backoffCoefficient {
				newSwapAmount := int64(float64(swapAmount) * backoffCoefficient)
				err := performSwap(info, channel, newSwapAmount, rule, span, loopdCtx, 0, swapAmountTarget)
				if err != nil {
					return err
				}
			}
			err := fmt.Errorf("reached max retries for reverse swap, failure reason: %v channel: %v on node: %v", swapStatus.GetFailureReason(), channel.GetChanId(), info.nodeInfo.GetIdentityPubkey())
			return err
		}

		return err
	} else if swapStatus.State == looprpc.SwapState_SUCCESS {
		//Success log: The swap was successful
		log.WithField("span", span).Infof("swap performed for channel %v, swap id: %v on node %v", channel.GetChanId(), resp.SwapId, info.nodeInfo.GetIdentityPubkey())

		swapType := "swap"

		//Add fees to counter
		recordSwapFees(swapStatus, info, swapType, channel)
	}
	return nil
}

func performReverseSwap(info ManageChannelLiquidityInfo, channel *lnrpc.Channel, swapAmount int64, rule nodeguard.LiquidityRule, span trace.Span, loopdCtx context.Context, retryCounter int, swapAmountTarget int64) error {
	// Check if it is a reverse swap to a wallet or to an address
	var address string
	if rule.IsReverseSwapWalletRule {
		//Request nodeguard a new destination address for the reverse swap
		walletRequest := &nodeguard.GetNewWalletAddressRequest{
			WalletId: *rule.ReverseSwapWalletId,
		}

		addrResponse, err := info.nodeguardClient.GetNewWalletAddress(info.ctx, walletRequest)
		if err != nil || addrResponse.GetAddress() == "" {
			log.WithField("span", span).Errorf("error requesting nodeguard a new wallet address: %v on node: %v", err, info.nodeInfo.GetIdentityPubkey())
			return err
		}

		address = addrResponse.Address
	} else {
		address = *rule.ReverseSwapAddress
	}

	//Perform the swap
	swapRequest := provider.ReverseSubmarineSwapRequest{
		SatsAmount:         swapAmount,
		ChannelSet:         []uint64{channel.GetChanId()},
		ReceiverBTCAddress: address,
	}

	log.WithField("span", span).Infof("requesting reverse submarine swap with amount: %d sats to BTC Address %s", swapRequest.SatsAmount, swapRequest.ReceiverBTCAddress)

	resp, err := info.provider.RequestReverseSubmarineSwap(loopdCtx, swapRequest, info.swapClient)
	if err != nil {
		log.WithField("span", span).Errorf("error requesting reverse swap: %v on node: %v", err, info.nodeInfo.GetIdentityPubkey())
		if strings.Contains(err.Error(), "channel balance too low for loop out amount") {
			limitSwapAmount := float64(helper.AbsInt64((channel.RemoteBalance - swapAmountTarget))) * backoffLimit
			if limitSwapAmount < float64(swapAmount)*backoffCoefficient {
				newSwapAmount := int64(float64(swapAmount) * backoffCoefficient)
				err = performReverseSwap(info, channel, newSwapAmount, rule, span, loopdCtx, 0, swapAmountTarget)
				if err != nil {
					return err
				}
			}
		}
		return err
	}

	log.WithField("span", span).Infof("reverse swap requested for channel %v, swap id: %v on node %v", channel.GetChanId(), resp.SwapId, info.nodeInfo.GetIdentityPubkey())

	//Monitor the swap
	swapStatus, err := info.provider.MonitorSwap(loopdCtx, resp.SwapId, info.swapClient)
	if err != nil {
		return err
	}

	if swapStatus.State == looprpc.SwapState_FAILED {
		//Error log: The swap was failed
		err := fmt.Errorf("failed reverse swap, failure reason: %v channel: %v on node: %v", swapStatus.GetFailureReason(), channel.GetChanId(), info.nodeInfo.GetIdentityPubkey())
		log.WithField("span", span).Error(err)

		if retryCounter < retries {
			//Retry the swap
			retryCounter++
			log.WithField("span", span).Infof("retrying reverse swap for channel %v, retry number: %v/%v on node %v", channel.GetChanId(), retryCounter, retries, info.nodeInfo.GetIdentityPubkey())
			err := performReverseSwap(info, channel, swapAmount, rule, span, loopdCtx, retryCounter, swapAmountTarget)
			if err != nil {
				return err
			}
		} else {
			limitSwapAmount := float64(helper.AbsInt64((channel.RemoteBalance - swapAmountTarget))) * backoffLimit
			if limitSwapAmount < float64(swapAmount)*backoffCoefficient {
				newSwapAmount := int64(float64(swapAmount) * backoffCoefficient)
				err := performReverseSwap(info, channel, newSwapAmount, rule, span, loopdCtx, 0, swapAmountTarget)
				if err != nil {
					return err
				}
			}
			err := fmt.Errorf("reached max retries for reverse swap, failure reason: %v channel: %v on node: %v", swapStatus.GetFailureReason(), channel.GetChanId(), info.nodeInfo.GetIdentityPubkey())
			return err
		}
	} else if swapStatus.State == looprpc.SwapState_SUCCESS {
		//Success log: The swap was successful
		log.WithField("span", span).Infof("reverse swap performed for channel %v, swap id: %v on node %v", channel.GetChanId(), resp.SwapId, info.nodeInfo.GetIdentityPubkey())

		swapType := "rswap"

		//Add fees to counter
		recordSwapFees(swapStatus, info, swapType, channel)
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

// Delete the channel balance metric in a prometheus gauge for a specific channel
func deleteChannelBalanceMetric(nodeHost string, channel *lnrpc.Channel, lightningClient lnrpc.LightningClient, context context.Context) {
	//Start span
	_, span := otel.Tracer("monitorChannel").Start(context, "deleteChannelBalanceMetric")
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

	prometheusMetrics.channelBalanceGauge.Delete(prometheus.Labels{
		"chan_id":            channelId,
		"local_node_pubkey":  localPubKey,
		"remote_node_pubkey": remotePubKey,
		"local_node_alias":   localAlias,
		"remote_node_alias":  remoteAlias,
		"active":             active,
		"initiator":          initiator,
	})
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
