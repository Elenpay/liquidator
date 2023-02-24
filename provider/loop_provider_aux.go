package provider

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/lightninglabs/loop/looprpc"
)

const (
	// FeeRateTotalParts defines the granularity of the fee rate.
	// Throughout the codebase, we'll use fix based arithmetic to compute
	// fees.
	FeeRateTotalParts = 1e6
)

var (
	// Define route independent max routing fees. We have currently no way
	// to get a reliable estimate of the routing fees. Best we can do is
	// the minimum routing fees, which is not very indicative.
	maxRoutingFeeBase = btcutil.Amount(10)

	maxRoutingFeeRate = int64(20000)
)

// Took from loopd/cmd/loop/main.go

func getMaxRoutingFee(amt btcutil.Amount) btcutil.Amount {
	return calcFee(amt, maxRoutingFeeBase, maxRoutingFeeRate)
}

// CalcFee returns the swap fee for a given swap amount.
func calcFee(amount, feeBase btcutil.Amount, feeRate int64) btcutil.Amount {
	return feeBase + amount*btcutil.Amount(feeRate)/
		btcutil.Amount(FeeRateTotalParts)
}

// Took from loopd/cmd/loop/main.go

type inLimits struct {
	maxMinerFee btcutil.Amount
	maxSwapFee  btcutil.Amount
}

// Took from loopd/cmd/loop/main.go

func getInLimits(quote *looprpc.InQuoteResponse) *inLimits {
	return &inLimits{
		// Apply a multiplier to the estimated miner fee, to not get
		// the swap canceled because fees increased in the mean time.
		maxMinerFee: btcutil.Amount(quote.HtlcPublishFeeSat) * 3,
		maxSwapFee:  btcutil.Amount(quote.SwapFeeSat),
	}
}

// Took from loopd/cmd/loop/main.go
type outLimits struct {
	maxSwapRoutingFee   btcutil.Amount
	maxPrepayRoutingFee btcutil.Amount
	maxMinerFee         btcutil.Amount
	maxSwapFee          btcutil.Amount
	maxPrepayAmt        btcutil.Amount
}

// Took from loopd/cmd/loop/main.go
func getOutLimits(amt btcutil.Amount,
	quote *looprpc.OutQuoteResponse) *outLimits {

	maxSwapRoutingFee := getMaxRoutingFee(amt)
	maxPrepayRoutingFee := getMaxRoutingFee(btcutil.Amount(
		quote.PrepayAmtSat,
	))
	maxPrepayAmt := btcutil.Amount(quote.PrepayAmtSat)

	return &outLimits{
		maxSwapRoutingFee:   maxSwapRoutingFee,
		maxPrepayRoutingFee: maxPrepayRoutingFee,

		// Apply a multiplier to the estimated miner fee, to not get
		// the swap canceled because fees increased in the mean time.
		maxMinerFee: btcutil.Amount(quote.HtlcSweepFeeSat) * 250,

		maxSwapFee:   btcutil.Amount(quote.SwapFeeSat),
		maxPrepayAmt: maxPrepayAmt,
	}
}
