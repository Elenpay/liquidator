package provider

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Elenpay/liquidator/customerrors"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/lightninglabs/loop/looprpc"
	"github.com/lightningnetwork/lnd/routing/route"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type LoopProvider struct {
	// Timestamp of when submarine swap lock was acquired
	submarineSwapLockTime *time.Time

	// Timestamp of when reverse submarine swap lock was acquired
	reverseSwapLockTime *time.Time

	// General mutex for managing lock state
	stateMutex sync.RWMutex
}

// acquireSubmarineSwapLock tries to acquire the submarine swap lock
func (l *LoopProvider) acquireSubmarineSwapLock() error {
	l.stateMutex.Lock()
	defer l.stateMutex.Unlock()

	// Check if there's an active lock and if it has expired
	if l.submarineSwapLockTime != nil {
		swapLockTimeout := viper.GetDuration("swapLockTimeout")
		if time.Since(*l.submarineSwapLockTime) < swapLockTimeout {
			return &customerrors.SwapInProgressError{
				Message: fmt.Sprintf("submarine swap is locked, started at %s, will expire at %s",
					l.submarineSwapLockTime.Format(time.RFC3339),
					l.submarineSwapLockTime.Add(swapLockTimeout).Format(time.RFC3339)),
			}
		}
	}

	// Acquire the lock
	now := time.Now()
	l.submarineSwapLockTime = &now
	log.Infof("Acquired submarine swap lock at %s", now.Format(time.RFC3339))
	return nil
}

// releaseSubmarineSwapLock releases the submarine swap lock
func (l *LoopProvider) releaseSubmarineSwapLock() {
	l.stateMutex.Lock()
	defer l.stateMutex.Unlock()

	if l.submarineSwapLockTime != nil {
		log.Infof("Released submarine swap lock that was acquired at %s", l.submarineSwapLockTime.Format(time.RFC3339))
		l.submarineSwapLockTime = nil
	}
}

// acquireReverseSwapLock tries to acquire the reverse swap lock
func (l *LoopProvider) acquireReverseSwapLock() error {
	l.stateMutex.Lock()
	defer l.stateMutex.Unlock()

	// Check if there's an active lock and if it has expired
	if l.reverseSwapLockTime != nil {
		swapLockTimeout := viper.GetDuration("swapLockTimeout")
		if time.Since(*l.reverseSwapLockTime) < swapLockTimeout {
			return &customerrors.SwapInProgressError{
				Message: fmt.Sprintf("reverse submarine swap is locked, started at %s, will expire at %s",
					l.reverseSwapLockTime.Format(time.RFC3339),
					l.reverseSwapLockTime.Add(swapLockTimeout).Format(time.RFC3339)),
			}
		}
	}

	// Acquire the lock
	now := time.Now()
	l.reverseSwapLockTime = &now
	log.Infof("Acquired reverse swap lock at %s", now.Format(time.RFC3339))
	return nil
}

// releaseReverseSwapLock releases the reverse swap lock
func (l *LoopProvider) releaseReverseSwapLock() {
	l.stateMutex.Lock()
	defer l.stateMutex.Unlock()

	if l.reverseSwapLockTime != nil {
		log.Infof("Released reverse swap lock that was acquired at %s", l.reverseSwapLockTime.Format(time.RFC3339))
		l.reverseSwapLockTime = nil
	}
}

// Submarine Swap L1->L2 based on loop (Loop In)
func (l *LoopProvider) RequestSubmarineSwap(ctx context.Context, request SubmarineSwapRequest, client looprpc.SwapClientClient) (SubmarineSwapResponse, error) {

	//Check that no sub swap is already in progress and acquire lock (rate limiting)
	err := l.acquireSubmarineSwapLock()
	if err != nil {
		log.Error(err)
		return SubmarineSwapResponse{}, err
	}

	// Note: Lock will expire automatically via timeout - no manual release for anti-spam purposes

	if request.SatsAmount <= 0 {
		//Create error
		err := fmt.Errorf("swap amount is <= 0")
		//Log error
		log.Error(err)

		return SubmarineSwapResponse{}, err

	}

	//Use the client to request the swap
	lastHopVertex, err := route.NewVertexFromStr(request.LastHopPubkey)
	if err != nil {
		return SubmarineSwapResponse{}, err
	}

	lastHopBytes := lastHopVertex[:]

	//Do a quote for loop in
	quote, err := client.GetLoopInQuote(ctx, &looprpc.QuoteRequest{
		Amt: request.SatsAmount,
		//ConfTarget:   1, //TODO Make this configurable
		ExternalHtlc:  true,
		Private:       false,
		LoopInLastHop: lastHopBytes,
	})

	if err != nil {
		log.Error(err)
		return SubmarineSwapResponse{}, err
	}

	limitFeesStr := viper.GetString("limitQuoteFees")
	limitFees, err := strconv.ParseFloat(limitFeesStr, 64)
	if err != nil {
		return SubmarineSwapResponse{}, err
	}

	sumFees := quote.SwapFeeSat + quote.HtlcPublishFeeSat

	maximumFeesAllowed := int64(float64(request.SatsAmount) * limitFees)
	if sumFees > maximumFeesAllowed {
		err := fmt.Errorf("swap fees are greater than max limit fees, quote fees: %d, maximum fees allowed: %d", sumFees, maximumFeesAllowed)
		log.Error(err)
		return SubmarineSwapResponse{}, err
	}

	//Get limits
	limits := getInLimits(quote)

	log.Debugf("loop in quote: %+v", quote)
	log.Debugf("loop in limits: %+v", limits)

	resp, err := client.LoopIn(ctx, &looprpc.LoopInRequest{
		Amt:            request.SatsAmount,
		MaxSwapFee:     int64(limits.maxSwapFee),
		MaxMinerFee:    int64(limits.maxMinerFee),
		LastHop:        lastHopBytes,
		ExternalHtlc:   true,
		HtlcConfTarget: 0,
		Label:          fmt.Sprintf("Swap in %d sats on date %s", request.SatsAmount, time.Now().Format(time.RFC3339)),
		Initiator:      "Liquidator",
		Private:        true,

		//TODO Review if hop hints are needed
	})

	if err != nil {
		//Log error
		log.Error(err)
		return SubmarineSwapResponse{}, err
	}

	//Id bytes to hex string
	swapId := hex.EncodeToString(resp.GetIdBytes())

	//Log response
	log.Debugf("submarine swap response: %+v", resp)
	//Return the response

	htlcAddress := resp.GetHtlcAddressP2Tr()

	if htlcAddress == "" {
		htlcAddress = resp.GetHtlcAddressP2Wsh()
	}

	response := SubmarineSwapResponse{
		SwapId:            swapId,
		InvoiceBTCAddress: htlcAddress,
	}

	return response, nil

}

// Reverse Submarine Swap L2->L1 based on loop (Loop Out)
// Reverse Submarine Swap L2->L1 based on loop (Loop Out)
func (l *LoopProvider) RequestReverseSubmarineSwap(ctx context.Context, request ReverseSubmarineSwapRequest, client looprpc.SwapClientClient) (ReverseSubmarineSwapResponse, error) {

	//Check that no other swap is in progress and acquire lock (rate limiting)
	err := l.acquireReverseSwapLock()
	if err != nil {
		log.Error(err)
		return ReverseSubmarineSwapResponse{}, err
	}

	// Note: Lock will expire automatically via timeout - no manual release for anti-spam purposes

	if request.SatsAmount <= 0 {
		//Create error
		err := fmt.Errorf("swap amount is <= 0")
		//Log error
		log.Error(err)

		return ReverseSubmarineSwapResponse{}, err

	}

	//Do a quote for loop out
	quote, err := client.LoopOutQuote(ctx, &looprpc.QuoteRequest{
		Amt:          request.SatsAmount,
		ConfTarget:   viper.GetInt32("sweepConfTarget"),
		ExternalHtlc: true,
		Private:      false,
	})

	if err != nil {
		log.Errorf("error getting quote for reverse submarine swap: %s", err)
		return ReverseSubmarineSwapResponse{}, err
	}

	limitQuoteFees := viper.GetFloat64("limitQuoteFees")

	//This fees are onchain + service fees, NOT L2 fees as they are not in the quote
	sumFees := quote.SwapFeeSat + quote.HtlcSweepFeeSat + quote.PrepayAmtSat
	maximumFeesAllowed := int64(float64(request.SatsAmount) * limitQuoteFees)

	if sumFees > maximumFeesAllowed {
		log.Warnf("swap quote fees (L1+Service estimation fees) are greater than max limit fees, quote fees: %d, maximum fees theoretically allowed: %d", sumFees, maximumFeesAllowed)
	}

	//Max swap routing fee (L2 fees) is a percentage of the swap amount
	l2MaxRoutingFeeRatio := viper.GetFloat64("limitFeesL2")
	maxSwapRoutingFee := int64(float64(request.SatsAmount) * l2MaxRoutingFeeRatio)

	log.Infof("max L2 routing fees for the swap: %d", maxSwapRoutingFee)

	//Get limits
	//Amt using btcutil
	amt := btcutil.Amount(request.SatsAmount)
	limits := getOutLimits(amt, quote)

	log.Debugf("loop out quote: %+v", quote)
	log.Debugf("loop out limits: %+v", limits)

	//Use the client to request the swap
	resp, err := client.LoopOut(ctx, &looprpc.LoopOutRequest{
		Amt:                 request.SatsAmount,
		Dest:                request.ReceiverBTCAddress,
		MaxMinerFee:         int64(limits.maxMinerFee),
		MaxPrepayAmt:        int64(limits.maxPrepayAmt),
		MaxSwapFee:          int64(limits.maxSwapFee),
		MaxPrepayRoutingFee: int64(limits.maxPrepayRoutingFee),
		MaxSwapRoutingFee:   maxSwapRoutingFee,
		// OutgoingChanSet:     request.ChannelSet, Disabled, evidence indicates this is not needed
		SweepConfTarget:   viper.GetInt32("sweepConfTarget"),
		HtlcConfirmations: 3,
		//The publication deadline is maximum the offset of the swap deadline conf plus the current time
		SwapPublicationDeadline: uint64(time.Now().Add(viper.GetDuration("swapPublicationOffset") * time.Minute).Unix()),
		Label:                   fmt.Sprintf("Reverse submarine swap %d sats on date %s for channels: %v", request.SatsAmount, time.Now().Format(time.RFC3339), request.ChannelSet),
		Initiator:               "Liquidator",
	})

	if err != nil {
		//Log error
		log.Errorf("error requesting reverse submarine swap: %s", err)
		return ReverseSubmarineSwapResponse{}, err
	}

	//Id bytes to hex string
	swapId := hex.EncodeToString(resp.GetIdBytes())

	//Log response
	log.Debugf("reverse submarine swap response: %+v", resp)

	//Return the response
	response := ReverseSubmarineSwapResponse{
		SwapId: swapId,
	}

	return response, nil

}

// Get the status of a swap by invoking SwapInfo method from the client
func (l *LoopProvider) GetSwapStatus(ctx context.Context, swapId string, client looprpc.SwapClientClient) (looprpc.SwapStatus, error) {

	log.Infof("getting swap status for swapId: %s", swapId)

	if len(swapId) == 0 {
		//Create error
		err := fmt.Errorf("swapId is empty")
		//Log error
		log.Error(err)

		return looprpc.SwapStatus{}, err
	}

	//Decode swapId from hex string to bytes
	swapIdBytes, err := hex.DecodeString(swapId)

	if err != nil {
		//Log error
		log.Errorf("error decoding swapId: %s", err)
		return looprpc.SwapStatus{}, err
	}

	//Get swap info
	swapInfo, err := client.SwapInfo(ctx, &looprpc.SwapInfoRequest{
		Id: swapIdBytes,
	})

	if err != nil {
		//Log error
		log.Errorf("error getting swap info: %s", err)
		return looprpc.SwapStatus{}, err
	}

	//Log response
	log.Debugf("swap info response: %+v", swapInfo)

	//Log success
	log.Infof("swap status for swapId: %s is %s", swapId, swapInfo.State.String())

	return *swapInfo, nil

}

// Monitor a swap status changes and stops when the swap is completed or failed
func (l *LoopProvider) MonitorSwap(ctx context.Context, swapId string, swapClient looprpc.SwapClientClient) (looprpc.SwapStatus, error) {

	if swapId == "" {
		err := fmt.Errorf("swapId is empty")
		log.Error(err)
		return looprpc.SwapStatus{}, err
	}

	//Decode swapId from hex string to bytes
	monitoredSwapIdBytes, err := hex.DecodeString(swapId)
	if err != nil {
		//Log error
		log.Errorf("error decoding swapId: %s", err)
		return looprpc.SwapStatus{}, err
	}

	for {

		//Get swap status
		swapInfo, err := swapClient.SwapInfo(ctx, &looprpc.SwapInfoRequest{
			Id: monitoredSwapIdBytes,
		})

		if err != nil {
			//Log error
			log.Errorf("error getting swap info: %s", err)
			return looprpc.SwapStatus{}, err
		}

		//Log response
		log.Debugf("swap info response: %+v", swapInfo)

		//If the swap is completed or failed, return the response
		if swapInfo.State == looprpc.SwapState_SUCCESS || swapInfo.State == looprpc.SwapState_FAILED {
			return *swapInfo, nil
		}

		time.Sleep(1 * time.Second)

	}

}
