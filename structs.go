package main

import (
	context "context"

	"github.com/Elenpay/liquidator/nodeguard"
	"github.com/Elenpay/liquidator/provider"
	"github.com/lightninglabs/loop/looprpc"
	"github.com/lightningnetwork/lnd/lnrpc"
)

type BaseInfo struct {
	nodeHost        string
	nodeInfo        lnrpc.GetInfoResponse
	nodeMacaroon    string
	loopdMacaroon   string
	lightningClients map[string]lnrpc.LightningClient
	nodeguardClient nodeguard.NodeGuardServiceClient
	swapClient      looprpc.SwapClientClient
	nodeCtx         context.Context
	provider        provider.Provider
}

type MonitorChannelsInfo struct {
	BaseInfo
}

type MonitorChannelInfo struct {
	BaseInfo
	channel          *lnrpc.Channel
	context          context.Context
	liquidationRules map[uint64][]nodeguard.LiquidityRule
}

type ManageChannelLiquidityInfo struct {
	BaseInfo
	channel             *lnrpc.Channel
	channelBalanceRatio float64
	channelRules        *[]nodeguard.LiquidityRule
	ctx                 context.Context
}
