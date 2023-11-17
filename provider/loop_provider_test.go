package provider

import (
	"context"
	"encoding/hex"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lightninglabs/loop/looprpc"
)

func TestLoopProvider_RequestSubmarineSwap(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	os.Setenv("LIMITFEES", "0.1")

	//Mock lightning swapClient GetLoopInQuote and LoopIn methods to return fake data
	swapClient := NewMockSwapClientClient(ctrl)

	swapClient.EXPECT().GetLoopInQuote(gomock.Any(), gomock.Any()).Return(&looprpc.InQuoteResponse{
		SwapFeeSat:        1,
		HtlcPublishFeeSat: 1,
		CltvDelta:         0,
		ConfTarget:        1,
	}, nil).Times(1)

	idBytes, err := hex.DecodeString("1234")
	if err != nil {
		t.Errorf("Error decoding hex string: %v", err)
	}

	swapClient.EXPECT().LoopIn(gomock.Any(), gomock.Any()).Return(&looprpc.SwapResponse{
		Id:               "",
		IdBytes:          idBytes,
		HtlcAddress:      "",
		HtlcAddressP2Wsh: "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
		HtlcAddressP2Tr:  "",
		ServerMessage:    "Test",
	}, nil).Times(1)

	//Mock ListSwaps to return fake data
	swapClient.EXPECT().ListSwaps(gomock.Any(), gomock.Any()).Return(&looprpc.ListSwapsResponse{
		Swaps: []*looprpc.SwapStatus{
			{
				Amt:              0,
				Id:               "",
				IdBytes:          idBytes,
				Type:             0,
				State:            0,
				FailureReason:    0,
				InitiationTime:   0,
				LastUpdateTime:   0,
				HtlcAddress:      "",
				HtlcAddressP2Wsh: "",
				HtlcAddressP2Tr:  "",
				CostServer:       0,
				CostOnchain:      0,
				CostOffchain:     0,
				LastHop:          []byte{},
				OutgoingChanSet:  []uint64{},
				Label:            "",
			},
		},
	}, nil).AnyTimes()

	type args struct {
		ctx     context.Context
		request SubmarineSwapRequest
		client  looprpc.SwapClientClient
	}
	tests := []struct {
		name    string
		l       *LoopProvider
		args    args
		want    SubmarineSwapResponse
		wantErr bool
	}{
		{
			name: "Test RequestSubmarineSwap_InvalidAmt",
			l:    &LoopProvider{},
			args: args{
				ctx: context.Background(),
				request: SubmarineSwapRequest{
					SatsAmount:    0,
					LastHopPubkey: "03485d8dcdd149c87553eeb80586eb2bece874d412e9f117304446ce189955d375",
				},
				client: swapClient,
			},
			want:    SubmarineSwapResponse{},
			wantErr: true,
		},
		{
			name: "Test RequestSubmarineSwap_InvalidAmt2",
			l:    &LoopProvider{},
			args: args{
				ctx: context.Background(),
				request: SubmarineSwapRequest{
					SatsAmount:    -1,
					LastHopPubkey: "03485d8dcdd149c87553eeb80586eb2bece874d412e9f117304446ce189955d375",
				},
				client: swapClient,
			},
			want:    SubmarineSwapResponse{},
			wantErr: true,
		},
		{
			name: "Test RequestSubmarineSwap_Valid",
			l:    &LoopProvider{},
			args: args{
				ctx: context.Background(),
				request: SubmarineSwapRequest{
					SatsAmount:    100000000,
					LastHopPubkey: "03485d8dcdd149c87553eeb80586eb2bece874d412e9f117304446ce189955d375",
				},
				client: swapClient,
			},
			want: SubmarineSwapResponse{
				SwapId:            hex.EncodeToString(idBytes),
				InvoiceBTCAddress: "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.l.RequestSubmarineSwap(tt.args.ctx, tt.args.request, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoopProvider.RequestSubmarineSwap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoopProvider.RequestSubmarineSwap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoopProvider_RequestReverseSubmarineSwap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//Mock lightning swapClient GetLoopInQuote and LoopIn methods to return fake data
	swapClient := NewMockSwapClientClient(ctrl)

	swapClient.EXPECT().LoopOutQuote(gomock.Any(), gomock.Any()).Return(&looprpc.OutQuoteResponse{
		SwapFeeSat:      1,
		PrepayAmtSat:    0,
		HtlcSweepFeeSat: 0,
		SwapPaymentDest: []byte{},
		CltvDelta:       0,
		ConfTarget:      1,
	}, nil).AnyTimes()

	idBytes, err := hex.DecodeString("1234")
	if err != nil {
		t.Errorf("Error decoding hex string: %v", err)
	}

	swapClient.EXPECT().LoopOut(gomock.Any(), gomock.Any()).Return(&looprpc.SwapResponse{
		Id:               "",
		IdBytes:          idBytes,
		HtlcAddress:      "",
		HtlcAddressP2Wsh: "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
		HtlcAddressP2Tr:  "",
		ServerMessage:    "Test",
	}, nil).AnyTimes()

	//Mock ListSwaps to return fake data

	swapClient.EXPECT().ListSwaps(gomock.Any(), gomock.Any()).Return(&looprpc.ListSwapsResponse{
		Swaps: []*looprpc.SwapStatus{
			{
				Amt:              0,
				Id:               "",
				IdBytes:          idBytes,
				Type:             0,
				State:            looprpc.SwapState_FAILED,
				FailureReason:    0,
				InitiationTime:   0,
				LastUpdateTime:   0,
				HtlcAddress:      "",
				HtlcAddressP2Wsh: "",
				HtlcAddressP2Tr:  "",
				CostServer:       0,
				CostOnchain:      0,
				CostOffchain:     0,
				LastHop:          []byte{},
				OutgoingChanSet:  []uint64{},
				Label:            "",
			},
		},
	}, nil).AnyTimes()

	type args struct {
		ctx     context.Context
		request ReverseSubmarineSwapRequest
		client  looprpc.SwapClientClient
	}
	tests := []struct {
		name    string
		l       *LoopProvider
		args    args
		want    ReverseSubmarineSwapResponse
		wantErr bool
	}{
		{
			name: "Test RequestReverseSubmarineSwap_InvalidAmt",
			l:    &LoopProvider{},
			args: args{
				ctx: context.Background(),
				request: ReverseSubmarineSwapRequest{
					SatsAmount: 0,
				},
				client: swapClient,
			},
			want:    ReverseSubmarineSwapResponse{},
			wantErr: true,
		},
		{
			name: "Test RequestReverseSubmarineSwap_InvalidAmt2",
			l:    &LoopProvider{},
			args: args{
				ctx: context.Background(),
				request: ReverseSubmarineSwapRequest{
					SatsAmount: -1,
				},
				client: swapClient,
			},
			want:    ReverseSubmarineSwapResponse{},
			wantErr: true,
		},
		{
			name: "Test RequestReverseSubmarineSwap_Valid",
			l:    &LoopProvider{},
			args: args{
				ctx: context.Background(),
				request: ReverseSubmarineSwapRequest{
					ReceiverBTCAddress: "",
					SatsAmount:         100000000,
					ChannelSet:         []uint64{},
				},
				client: swapClient,
			},
			want: ReverseSubmarineSwapResponse{
				SwapId: hex.EncodeToString(idBytes),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.l.RequestReverseSubmarineSwap(tt.args.ctx, tt.args.request, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoopProvider.RequestSubmarineSwap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoopProvider.RequestSubmarineSwap() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestLoopProvider_GetSwapStatus(t *testing.T) {
	//Mock swap client to return fake data
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//Mock SwapInfo
	client := NewMockSwapClientClient(ctrl)
	idBytes, err := hex.DecodeString("1234")

	if err != nil {
		t.Errorf("Error decoding hex string: %v", err)
	}

	status := looprpc.SwapStatus{
		Amt:              1000,
		Id:               "1234",
		IdBytes:          idBytes,
		Type:             0,
		State:            looprpc.SwapState_INITIATED,
		FailureReason:    0,
		InitiationTime:   0,
		LastUpdateTime:   0,
		HtlcAddress:      "",
		HtlcAddressP2Wsh: "",
		HtlcAddressP2Tr:  "",
		CostServer:       0,
		CostOnchain:      0,
		CostOffchain:     0,
		LastHop:          []byte{},
		OutgoingChanSet:  []uint64{},
		Label:            "",
	}

	client.EXPECT().SwapInfo(gomock.Any(), gomock.Any()).Return(&status, nil).AnyTimes()

	type args struct {
		ctx     context.Context
		request string
		client  looprpc.SwapClientClient
	}

	tests := []struct {
		name    string
		l       *LoopProvider
		args    args
		want    looprpc.SwapStatus
		wantErr bool
	}{
		{
			name: "Test GetSwapStatus_InvalidId",
			l:    &LoopProvider{},
			args: args{
				ctx:     context.Background(),
				request: "",
				client:  nil,
			},
			want: looprpc.SwapStatus{},

			wantErr: true,
		},
		{
			name: "Test GetSwapStatus_Valid",
			l:    &LoopProvider{},
			args: args{
				ctx:     context.Background(),
				request: "1234",
				client:  client,
			},
			want:    status,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.l.GetSwapStatus(tt.args.ctx, tt.args.request, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoopProvider.GetSwapStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoopProvider.GetSwapStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkSubmarineSwapNotInProgress(t *testing.T) {

	//Swap client with ListSwaps returning a swap in progress
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//Mock SwapInfo
	swapClientWithOngoingSwaps := NewMockSwapClientClient(ctrl)
	idBytes, err := hex.DecodeString("1234")
	if err != nil {
		t.Errorf("Error decoding hex string: %v", err)
	}

	swapClientWithOngoingSwaps.EXPECT().ListSwaps(gomock.Any(), gomock.Any()).Return(&looprpc.ListSwapsResponse{
		Swaps: []*looprpc.SwapStatus{
			{
				Amt:           0,
				Id:            "",
				IdBytes:       idBytes,
				Type:          looprpc.SwapType_LOOP_IN,
				State:         looprpc.SwapState_INITIATED,
				FailureReason: 0,
				//InitiationTime 4 hours ago
				InitiationTime:   time.Now().Add(-4 * time.Hour).UnixNano(),
				LastUpdateTime:   time.Now().Add(-4 * time.Hour).UnixNano(),
				HtlcAddress:      "",
				HtlcAddressP2Wsh: "",
				HtlcAddressP2Tr:  "",
				CostServer:       0,
				CostOnchain:      0,
				CostOffchain:     0,
				LastHop:          []byte{},
				OutgoingChanSet:  []uint64{},
				Label:            "",
			},
		},
	}, nil).AnyTimes()

	//Swap client with ListSwaps returning no swaps
	swapClientWithNoOngoingSwaps := NewMockSwapClientClient(ctrl)
	swapClientWithNoOngoingSwaps.EXPECT().ListSwaps(gomock.Any(), gomock.Any()).Return(&looprpc.ListSwapsResponse{
		Swaps: []*looprpc.SwapStatus{
			{
				Amt:           0,
				Id:            "",
				IdBytes:       idBytes,
				Type:          looprpc.SwapType_LOOP_IN,
				State:         looprpc.SwapState_INITIATED,
				FailureReason: 0,
				//InitiationTime is more than 24 hours ago, stuck but ignored
				InitiationTime:   time.Now().Add(-25 * time.Hour).UnixNano(),
				LastUpdateTime:   time.Now().Add(-25 * time.Hour).UnixNano(),
				HtlcAddress:      "",
				HtlcAddressP2Wsh: "",
				HtlcAddressP2Tr:  "",
				CostServer:       0,
				CostOnchain:      0,
				CostOffchain:     0,
				LastHop:          []byte{},
				OutgoingChanSet:  []uint64{},
				Label:            "",
			},
		},
	}, nil).AnyTimes()

	type args struct {
		ctx    context.Context
		client looprpc.SwapClientClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "checkSubmarineSwapNotInProgress_valid",
			args: args{
				ctx:    context.Background(),
				client: swapClientWithNoOngoingSwaps,
			},
			wantErr: false,
		},
		{
			name: "checkSubmarineSwapNotInProgress_ErrorOngoingSwap",
			args: args{
				ctx:    context.Background(),
				client: swapClientWithOngoingSwaps,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkSubmarineSwapNotInProgress(tt.args.ctx, tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("checkSubmarineSwapNotInProgress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// checkReverseSubmarineSwapNotInProgress
func Test_checkReverseSubmarineSwapNotInProgress(t *testing.T) {

	//Swap client with ListSwaps returning a swap in progress
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	channelSet := []uint64{1, 2, 3}

	//Mock SwapInfo
	swapClientWithOngoingSwaps := NewMockSwapClientClient(ctrl)
	idBytes, err := hex.DecodeString("1234")
	if err != nil {
		t.Errorf("Error decoding hex string: %v", err)
	}

	swapClientWithOngoingSwaps.EXPECT().ListSwaps(gomock.Any(), gomock.Any()).Return(&looprpc.ListSwapsResponse{
		Swaps: []*looprpc.SwapStatus{
			{
				Amt:              0,
				Id:               "",
				IdBytes:          idBytes,
				Type:             looprpc.SwapType_LOOP_OUT,
				State:            looprpc.SwapState_INITIATED,
				FailureReason:    0,
				InitiationTime:   0,
				LastUpdateTime:   0,
				HtlcAddress:      "",
				HtlcAddressP2Wsh: "",
				HtlcAddressP2Tr:  "",
				CostServer:       0,
				CostOnchain:      0,
				CostOffchain:     0,
				LastHop:          []byte{},
				OutgoingChanSet:  channelSet,
				Label:            "",
			},
		},
	}, nil).AnyTimes()

	//Swap client with ListSwaps returning no swaps
	swapClientWithNoOngoingSwaps := NewMockSwapClientClient(ctrl)
	swapClientWithNoOngoingSwaps.EXPECT().ListSwaps(gomock.Any(), gomock.Any()).Return(&looprpc.ListSwapsResponse{
		Swaps: []*looprpc.SwapStatus{},
	}, nil).AnyTimes()

	type args struct {
		ctx     context.Context
		client  looprpc.SwapClientClient
		request ReverseSubmarineSwapRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "checkReverseSubmarineSwapNotInProgress_valid",
			args: args{
				ctx:    context.Background(),
				client: swapClientWithNoOngoingSwaps,
				request: ReverseSubmarineSwapRequest{
					ReceiverBTCAddress: "",
					SatsAmount:         0,
					ChannelSet:         channelSet,
				},
			},
			wantErr: false,
		},
		{
			name: "checkReverseSubmarineSwapNotInProgress_ErrorOngoingSwap",
			args: args{
				ctx:    context.Background(),
				client: swapClientWithOngoingSwaps,
				request: ReverseSubmarineSwapRequest{
					ReceiverBTCAddress: "",
					SatsAmount:         0,
					ChannelSet:         channelSet,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkReverseSubmarineSwapNotInProgress(tt.args.ctx, tt.args.client, tt.args.request); (err != nil) != tt.wantErr {
				t.Errorf("checkReverseSubmarineSwapNotInProgress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

}

func TestLoopProvider_MonitorSwap(t *testing.T) {

	ctrl := gomock.NewController(t)

	//Success Swap Client
	mockSwapClientSuccess := NewMockSwapClientClient(ctrl)

	idBytes, err := hex.DecodeString("1234")
	if err != nil {
		t.Errorf("Error decoding hex string: %v", err)
	}

	//Mock SwapInfo method
	swapStatusSuccess := &looprpc.SwapStatus{
		IdBytes: idBytes,
		State:   looprpc.SwapState_SUCCESS,
	}
	mockSwapClientSuccess.EXPECT().SwapInfo(gomock.Any(), gomock.Any()).Return(swapStatusSuccess, nil).AnyTimes()

	//List Swaps

	mockSwapClientSuccess.EXPECT().ListSwaps(gomock.Any(), gomock.Any()).Return(&looprpc.ListSwapsResponse{
		Swaps: []*looprpc.SwapStatus{
			{

				IdBytes: idBytes,
			},
		},
	}, nil).AnyTimes()

	//Failure Swap Client
	mockSwapClientFailure := NewMockSwapClientClient(ctrl)

	swapStatusFailure := &looprpc.SwapStatus{
		State: looprpc.SwapState_FAILED,
	}

	//Mock SwapInfo method
	mockSwapClientFailure.EXPECT().SwapInfo(gomock.Any(), gomock.Any()).Return(swapStatusFailure, nil).AnyTimes()

	mockSwapClientFailure.EXPECT().ListSwaps(gomock.Any(), gomock.Any()).Return(&looprpc.ListSwapsResponse{
		Swaps: []*looprpc.SwapStatus{
			{
				IdBytes: idBytes,
			},
		},
	}, nil).AnyTimes()

	//Swap client for swap not found

	//Failure Swap Client
	mockSwapSwapNotFound := NewMockSwapClientClient(ctrl)

	mockSwapSwapNotFound.EXPECT().ListSwaps(gomock.Any(), gomock.Any()).Return(&looprpc.ListSwapsResponse{
		Swaps: []*looprpc.SwapStatus{},
	}, nil).AnyTimes()

	type args struct {
		ctx        context.Context
		swapId     string
		swapClient looprpc.SwapClientClient
	}
	tests := []struct {
		name    string
		l       *LoopProvider
		args    args
		want    looprpc.SwapStatus
		wantErr bool
	}{
		{
			name: "MonitorSwap_Success",
			l:    &LoopProvider{},
			args: args{
				ctx:        context.TODO(),
				swapId:     "1234",
				swapClient: mockSwapClientSuccess,
			},
			want:    *swapStatusSuccess,
			wantErr: false,
		},
		{
			name: "MonitorSwap_Failed",
			l:    &LoopProvider{},
			args: args{
				ctx:        context.TODO(),
				swapId:     "1234",
				swapClient: mockSwapClientFailure,
			},
			want:    *swapStatusFailure,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LoopProvider{}
			got, err := l.MonitorSwap(tt.args.ctx, tt.args.swapId, tt.args.swapClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoopProvider.MonitorSwap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoopProvider.MonitorSwap() = %v, want %v", got, tt.want)
			}
		})
	}
}
