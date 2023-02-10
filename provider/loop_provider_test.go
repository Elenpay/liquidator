package provider

import (
	"context"
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lightninglabs/loop/looprpc"
)

func TestLoopProvider_RequestSubmarineSwap(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//Mock lightning client GetLoopInQuote and LoopIn methods to return fake data
	client := NewMockSwapClientClient(ctrl)

	client.EXPECT().GetLoopInQuote(gomock.Any(), gomock.Any()).Return(&looprpc.InQuoteResponse{
		SwapFeeSat:        1,
		HtlcPublishFeeSat: 1,
		CltvDelta:         0,
		ConfTarget:        1,
	}, nil).Times(1)

	idBytes, err := hex.DecodeString("1234")
	if err != nil {
		t.Errorf("Error decoding hex string: %v", err)
	}

	client.EXPECT().LoopIn(gomock.Any(), gomock.Any()).Return(&looprpc.SwapResponse{
		Id:               "",
		IdBytes:          idBytes,
		HtlcAddress:      "",
		HtlcAddressP2Wsh: "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
		HtlcAddressP2Tr:  "",
		ServerMessage:    "Test",
	}, nil).Times(1)

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
					SatsAmount: 0,
				},
				client: nil,
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
					SatsAmount: -1,
				},
				client: nil,
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
					SatsAmount: 100000000,
				},
				client: client,
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

	//Mock lightning client GetLoopInQuote and LoopIn methods to return fake data
	client := NewMockSwapClientClient(ctrl)

	client.EXPECT().LoopOutQuote(gomock.Any(), gomock.Any()).Return(&looprpc.OutQuoteResponse{
		SwapFeeSat:      1,
		PrepayAmtSat:    0,
		HtlcSweepFeeSat: 0,
		SwapPaymentDest: []byte{},
		CltvDelta:       0,
		ConfTarget:      1,
	}, nil).Times(1)

	idBytes, err := hex.DecodeString("1234")
	if err != nil {
		t.Errorf("Error decoding hex string: %v", err)
	}

	client.EXPECT().LoopOut(gomock.Any(), gomock.Any()).Return(&looprpc.SwapResponse{
		Id:               "",
		IdBytes:          idBytes,
		HtlcAddress:      "",
		HtlcAddressP2Wsh: "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
		HtlcAddressP2Tr:  "",
		ServerMessage:    "Test",
	}, nil).Times(1)

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
				client: nil,
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
				client: nil,
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
				client: client,
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

	client.EXPECT().SwapInfo(gomock.Any(), gomock.Any()).Return(&status, nil).Times(1)

	type args struct {
		ctx     context.Context
		request string
		client  looprpc.SwapClientClient
	}

	tests := []struct {
		name    string
		l       *LoopProvider
		args    args
		want    *looprpc.SwapStatus
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
			want: &looprpc.SwapStatus{},

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
			want:    &status,
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
