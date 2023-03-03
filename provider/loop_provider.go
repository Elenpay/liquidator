package provider

import (
	"context"
	"encoding/hex"
	"fmt"
	"reflect"
	"time"

	"github.com/Elenpay/liquidator/errors"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/lightninglabs/loop/looprpc"
	"github.com/lightningnetwork/lnd/routing/route"
	log "github.com/sirupsen/logrus"
)

type LoopProvider struct {
}

// Submarine Swap L1->L2 based on loop (Loop In)
func (l *LoopProvider) RequestSubmarineSwap(ctx context.Context, request SubmarineSwapRequest, client looprpc.SwapClientClient) (SubmarineSwapResponse, error) {

	log.Infof("requesting submarine swap with amount: %d sats", request.SatsAmount)

	//Check that no sub swap is already in progress
	err := checkSubmarineSwapNotInProgress(ctx, client)
	if err != nil {
		log.Error(err)
		return SubmarineSwapResponse{}, err
	}

	if request.SatsAmount <= 0 {
		//Create error
		err := fmt.Errorf("swap amount is <= 0")
		//Log error
		log.Error(err)

		return SubmarineSwapResponse{}, err

	}

	//Do a quote for loop in
	quote, err := client.GetLoopInQuote(ctx, &looprpc.QuoteRequest{
		Amt: request.SatsAmount,
		//ConfTarget:   1, //TODO Make this configurable
		ExternalHtlc: true,
		Private:      false,
	})

	if err != nil {
		log.Error(err)
		return SubmarineSwapResponse{}, err
	}

	//Get limits
	limits := getInLimits(quote)

	log.Debugf("loop in quote: %+v", quote)
	log.Debugf("loop in limits: %+v", limits)

	//Use the client to request the swap
	lastHopVertex, err := route.NewVertexFromStr(request.LastHopPubkey)
	if err != nil {
		return SubmarineSwapResponse{}, err
	}

	lastHopBytes := lastHopVertex[:]

	resp, err := client.LoopIn(ctx, &looprpc.LoopInRequest{
		Amt:            request.SatsAmount,
		MaxSwapFee:     int64(limits.maxSwapFee),
		MaxMinerFee:    int64(limits.maxMinerFee),
		LastHop:        lastHopBytes,
		ExternalHtlc:   true,
		HtlcConfTarget: 0,
		Label:          fmt.Sprintf("Submarine swap %d sats on date %s", request.SatsAmount, time.Now().Format(time.RFC3339)),
		Initiator:      "Liquidator",
		Private:        false,
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

	response := SubmarineSwapResponse{
		SwapId:            swapId,
		InvoiceBTCAddress: resp.GetHtlcAddressP2Wsh(),
	}

	log.Infof("submarine swap request successful. SwapId: %s Server message: %s", swapId, resp.GetServerMessage())

	return response, nil

}

// Check that a Submarine Swap is not already in progress, by now the limit is one swap at a time for Swaps L1->L2
func checkSubmarineSwapNotInProgress(ctx context.Context, client looprpc.SwapClientClient) error {

	//Invoking ListSwaps check that a swap is not already in progress based on channelId of the request
	swaps, err := client.ListSwaps(ctx, &looprpc.ListSwapsRequest{})
	if err != nil {
		log.Error(err)
		return err
	}

	//Filter swaps of Loop In type
	var loopInSwaps []*looprpc.SwapStatus
	for _, swap := range swaps.Swaps {
		if swap.Type == looprpc.SwapType_LOOP_IN {
			loopInSwaps = append(loopInSwaps, swap)
		}
	}

	//CHeck that all the swaps status are either SUCCESS or FAILED, meaning that they are not in progress
	for _, swap := range loopInSwaps {
		if swap.State != looprpc.SwapState_SUCCESS && swap.State != looprpc.SwapState_FAILED {
			//Create error  of Swap already in progress

			id := hex.EncodeToString(swap.GetIdBytes())
			errMessage := fmt.Sprintf("another submarine swap is already in progress, swap id: %s", id)

			swapInProgressErr := errors.SwapInProgressError{
				Message: errMessage,
			}

			return &swapInProgressErr
		}
	}

	return nil
}

// Function that checks that a reverse submarine swap is only one per this channelid
func checkReverseSubmarineSwapNotInProgress(ctx context.Context, client looprpc.SwapClientClient, request ReverseSubmarineSwapRequest) error {

	//Invoking ListSwaps check that a swap is not already in progress based on channelId of the request
	swaps, err := client.ListSwaps(ctx, &looprpc.ListSwapsRequest{})
	if err != nil {
		log.Error(err)
		return err
	}

	//If there are no swaps, return nil
	if len(swaps.Swaps) == 0 {
		return nil
	}

	//Filter swaps of Loop Out type
	var loopOutSwaps []*looprpc.SwapStatus
	for _, swap := range swaps.Swaps {
		if swap.Type == looprpc.SwapType_LOOP_OUT && reflect.DeepEqual(swap.OutgoingChanSet, request.ChannelSet) {
			loopOutSwaps = append(loopOutSwaps, swap)
		}
	}

	//CHeck that all the swaps status are either SUCCESS or FAILED, meaning that they are not in progress
	for _, swap := range loopOutSwaps {
		if swap.State != looprpc.SwapState_SUCCESS && swap.State != looprpc.SwapState_FAILED {
			//Create error
			id := hex.EncodeToString(swap.GetIdBytes())
			errMessage := fmt.Sprintf("another submarine swap is already in progress, swap id: %s", id)

			swapInProgressErr := errors.SwapInProgressError{
				Message: errMessage,
			}

			return &swapInProgressErr
		}

	}
	return nil
}

// Reverse Submarine Swap L2->L1 based on loop (Loop Out)
func (l *LoopProvider) RequestReverseSubmarineSwap(ctx context.Context, request ReverseSubmarineSwapRequest, client looprpc.SwapClientClient) (ReverseSubmarineSwapResponse, error) {

	log.Infof("requesting reverse submarine swap with amount: %d sats to BTC Address %s", request.SatsAmount, request.ReceiverBTCAddress)

	//Check that no other swap is in progress
	err := checkReverseSubmarineSwapNotInProgress(ctx, client, request)
	if err != nil {
		log.Error(err)
		return ReverseSubmarineSwapResponse{}, err
	}

	if request.SatsAmount <= 0 {
		//Create error
		err := fmt.Errorf("swap amount is <= 0")
		//Log error
		log.Error(err)

		return ReverseSubmarineSwapResponse{}, err

	}

	//Do a quote for loop out
	quote, err := client.LoopOutQuote(ctx, &looprpc.QuoteRequest{
		Amt: request.SatsAmount,
		//ConfTarget:   1, //TODO Make this configurable
		ExternalHtlc: true,
		Private:      false,
	})

	if err != nil {

		log.Errorf("error getting quote for reverse submarine swap: %s", err)
		return ReverseSubmarineSwapResponse{}, err
	}

	//Get limits
	//Amt using btcutil
	amt := btcutil.Amount(request.SatsAmount)
	limits := getOutLimits(amt, quote)

	log.Debugf("loop out quote: %+v", quote)
	log.Debugf("loop out limits: %+v", limits)

	//Use the client to request the swap
	resp, err := client.LoopOut(ctx, &looprpc.LoopOutRequest{
		Amt:                     request.SatsAmount,
		Dest:                    request.ReceiverBTCAddress,
		MaxMinerFee:             int64(limits.maxMinerFee),
		MaxPrepayAmt:            int64(limits.maxPrepayAmt),
		MaxSwapFee:              int64(limits.maxSwapFee),
		MaxPrepayRoutingFee:     int64(limits.maxPrepayRoutingFee),
		MaxSwapRoutingFee:       int64(limits.maxSwapRoutingFee),
		OutgoingChanSet:         request.ChannelSet,
		SweepConfTarget:         2, //TODO Make this configurable
		HtlcConfirmations:       2,
		SwapPublicationDeadline: uint64(time.Now().Unix()),
		Label:                   fmt.Sprintf("Reverse submarine swap %d sats on date %s", request.SatsAmount, time.Now().Format(time.RFC3339)),
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

	log.Infof("reverse submarine swap request successful. SwapId: %s Server message: %s", swapId, resp.GetServerMessage())

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

	var response looprpc.SwapStatus

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
			response = *swapInfo
			break
		}

		time.Sleep(1 * time.Second)

	}

	return response, nil

}
