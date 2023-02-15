package provider

import (
	"context"

	"github.com/lightninglabs/loop/looprpc"
)

type Provider interface {

	//Submarine Swap L1->L2
	RequestSubmarineSwap(context.Context, SubmarineSwapRequest, looprpc.SwapClientClient) (SubmarineSwapResponse, error)

	//Reverse Submarine Swap L2->L1
	RequestReverseSubmarineSwap(context.Context, ReverseSubmarineSwapRequest, looprpc.SwapClientClient) (ReverseSubmarineSwapResponse, error)

	//Monitor Swap
	MonitorSwap(context.Context, string, looprpc.SwapClientClient) (*looprpc.SwapStatus, error)
}

// Provider-agnostic request for a submarine swap
type SubmarineSwapRequest struct {
	SatsAmount int64
	//Last hop node to identify which channel to use, if multiple channels are with this node then there is no way to know which one will be used
	LastHopPubkey string
}

// Provider-agnostic response for a submarine swap
type SubmarineSwapResponse struct {
	SwapId string
	//L1 address to send the funds to pay for the submarine swap
	InvoiceBTCAddress string
}

// Provider-agnostic request for a reverse submarine swap
type ReverseSubmarineSwapRequest struct {
	//L1 Address to receive L2 funds
	ReceiverBTCAddress string
	SatsAmount         int64
	ChannelSet         []uint64
}

// Provider-agnostic response for a reverse submarine swap
type ReverseSubmarineSwapResponse struct {
	SwapId string
}
