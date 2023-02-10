package provider

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/btcsuite/btcutil"
	"github.com/lightninglabs/loop/looprpc"
	log "github.com/sirupsen/logrus"
)

type LoopProvider struct {
}

// Submarine Swap L1->L2 based on loop (Loop In)
func (l *LoopProvider) RequestSubmarineSwap(ctx context.Context, request SubmarineSwapRequest, client looprpc.SwapClientClient) (SubmarineSwapResponse, error) {

	log.Infof("requesting submarine swap with amount: %d sats", request.SatsAmount)

	if request.SatsAmount <= 0 {
		//Create error
		err := fmt.Errorf("swap amount is <= 0")
		//Log error
		log.Error(err)

		return SubmarineSwapResponse{}, err

	}

	//Do a quote for loop in
	quote, err := client.GetLoopInQuote(ctx, &looprpc.QuoteRequest{
		Amt:          request.SatsAmount,
		ConfTarget:   1, //TODO Make this configurable
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
	resp, err := client.LoopIn(ctx, &looprpc.LoopInRequest{
		Amt:            request.SatsAmount,
		MaxMinerFee:    int64(limits.maxMinerFee), //TODO Make this configurable
		MaxSwapFee:     int64(limits.maxSwapFee),  //TODO Make this configurable
		ExternalHtlc:   true,
		HtlcConfTarget: 3, //TODO Make this configurable
		Label:          fmt.Sprintf("Submarine swap %d sats on date %s", request.SatsAmount, time.Now().Format(time.RFC3339)),
		Initiator:      "Liquidator",
		Private:        false,
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

// Reverse Submarine Swap L2->L1 based on loop (Loop Out)
func (l *LoopProvider) RequestReverseSubmarineSwap(ctx context.Context, request ReverseSubmarineSwapRequest, client looprpc.SwapClientClient) (ReverseSubmarineSwapResponse, error) {

	log.Infof("requesting reverse submarine swap with amount: %d sats to BTC Address %s", request.SatsAmount, request.ReceiverBTCAddress)

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
		ConfTarget:   1, //TODO Make this configurable
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
		SweepConfTarget:         1, //TODO Make this configurable
		HtlcConfirmations:       1,
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
func (l *LoopProvider) GetSwapStatus(ctx context.Context, swapId string, client looprpc.SwapClientClient) (*looprpc.SwapStatus, error) {

	log.Infof("getting swap status for swapId: %s", swapId)

	if len(swapId) == 0 {
		//Create error
		err := fmt.Errorf("swapId is empty")
		//Log error
		log.Error(err)

		return &looprpc.SwapStatus{}, err
	}

	//Decode swapId from hex string to bytes
	swapIdBytes, err := hex.DecodeString(swapId)

	if err != nil {
		//Log error
		log.Errorf("error decoding swapId: %s", err)
		return &looprpc.SwapStatus{}, err
	}

	//Get swap info
	swapInfo, err := client.SwapInfo(ctx, &looprpc.SwapInfoRequest{
		Id: swapIdBytes,
	})

	if err != nil {
		//Log error
		log.Errorf("error getting swap info: %s", err)
		return &looprpc.SwapStatus{}, err
	}

	//Log response
	log.Debugf("swap info response: %+v", swapInfo)

	//Log success
	log.Infof("swap status for swapId: %s is %s", swapId, swapInfo.State.String())

	return swapInfo, nil

}
