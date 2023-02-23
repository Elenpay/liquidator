package main

import (
	context "context"

	"github.com/Elenpay/liquidator/nodeguard"
	"github.com/Elenpay/liquidator/provider"
	"github.com/lightninglabs/loop/looprpc"
	"github.com/lightningnetwork/lnd/lnrpc"
)

type MonitorChannelsInfo struct {
	nodeHost        string
	nodeInfo        lnrpc.GetInfoResponse
	nodeMacaroon    string
	loopdMacaroon   string
	lightningClient lnrpc.LightningClient
	nodeguardClient nodeguard.NodeGuardServiceClient
	swapClient      looprpc.SwapClientClient
	nodeCtx         context.Context
}

type MonitorChannelInfo struct {
	channel          *lnrpc.Channel
	nodeHost         string
	lightningClient  lnrpc.LightningClient
	context          context.Context
	liquidationRules map[uint64][]nodeguard.LiquidityRule
	swapClient       looprpc.SwapClientClient
	nodeguardClient  nodeguard.NodeGuardServiceClient
	loopProvider     provider.Provider
	loopdMacaroon    string
	nodeInfo         lnrpc.GetInfoResponse
}

type ManageChannelLiquidityInfo struct {
	channel             *lnrpc.Channel
	channelBalanceRatio float64
	channelRules        *[]nodeguard.LiquidityRule
	swapClientClient    looprpc.SwapClientClient
	nodeguardClient     nodeguard.NodeGuardServiceClient
	loopProvider        provider.Provider
	loopdMacaroon       string
	nodeInfo            lnrpc.GetInfoResponse
	ctx                 context.Context
}
