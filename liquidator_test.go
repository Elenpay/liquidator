package main

import (
	context "context"
	"testing"

	"github.com/Elenpay/liquidator/nodeguard"
	"github.com/Elenpay/liquidator/provider"
	gomock "github.com/golang/mock/gomock"
	"github.com/lightninglabs/loop/looprpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/prometheus/client_golang/prometheus"
)

// Tear up method
func TestMain(m *testing.M) {
	//Tear up
	initMetrics(prometheus.NewRegistry())

	//Run tests
	m.Run()

	//Tear down
}

func Test_recordChannelBalance(t *testing.T) {
	type args struct {
		channel       *lnrpc.Channel
		wantedBalance float64
		expectedError bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"Test 1 positive", args{
			&lnrpc.Channel{
				ChanId:        1,
				LocalBalance:  100,
				RemoteBalance: 900,
				Capacity:      1000,
			},
			0.9, false}},
		{"Test 2 positive", args{
			&lnrpc.Channel{
				ChanId:        1,
				LocalBalance:  900,
				RemoteBalance: 100,
				Capacity:      1000,
			},
			0.1, false}},
		{"Test 3 negative", args{
			&lnrpc.Channel{
				ChanId:        1,
				LocalBalance:  1000,
				RemoteBalance: 900,
				Capacity:      0,
			},
			-1, true}},
	}

	for _, tt := range tests {
		t.Logf("Running test: %v", tt.name)

		t.Run(tt.name, func(t *testing.T) {
			actualBalance, err := getChannelBalanceRatio(tt.args.channel, context.TODO())
			//If we expect an error and we don't get one, fail
			if tt.args.expectedError && err == nil {
				t.Errorf("Error: %v", err)
			}

			//If balance is not what we think it should be, fail
			if actualBalance != tt.args.wantedBalance {
				t.Errorf("Expected balance: %v, actual balance: %v", tt.args.wantedBalance, actualBalance)
			}
		})
	}
}

func Test_manageChannelLiquidity(t *testing.T) {

	//gomock controller
	mockCtrl := gomock.NewController(t)

	//Mock nodeguard client
	mockNodeGuardClient := nodeguard.NewMockNodeGuardServiceClient(mockCtrl)

	//Mock provider valid swaps
	mockProvider := provider.NewMockProvider(mockCtrl)

	//Mock RequestReverseSubmarineSwap
	mockProvider.EXPECT().RequestReverseSubmarineSwap(gomock.Any(), gomock.Any(), gomock.Any()).Return(provider.ReverseSubmarineSwapResponse{
		SwapId: "1234",
	}, nil).AnyTimes()

	//Mock RequestSubmarineSwap
	mockProvider.EXPECT().RequestSubmarineSwap(gomock.Any(), gomock.Any(), gomock.Any()).Return(provider.SubmarineSwapResponse{
		SwapId:            "1234",
		InvoiceBTCAddress: "bcrt1q6zszlnxhlq0lsmfc42nkwgqedy9kvmvmxhkvme",
	}, nil).AnyTimes()

	//MockRequestSubmarineSwap
	mockProvider.EXPECT().MonitorSwap(gomock.Any(), gomock.Any(), gomock.Any()).Return(looprpc.SwapStatus{
		Amt:              0,
		Id:               "",
		IdBytes:          []byte{},
		Type:             0,
		State:            looprpc.SwapState_SUCCESS,
		FailureReason:    looprpc.FailureReason_FAILURE_REASON_INSUFFICIENT_VALUE,
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
	}, nil).AnyTimes()

	//Mock provider invalid swaps
	mockProviderInvalid := provider.NewMockProvider(mockCtrl)

	//Mock RequestReverseSubmarineSwap
	mockProviderInvalid.EXPECT().RequestReverseSubmarineSwap(gomock.Any(), gomock.Any(), gomock.Any()).Return(provider.ReverseSubmarineSwapResponse{
		SwapId: "1234",
	}, nil).AnyTimes()

	//Mock RequestSubmarineSwap
	mockProviderInvalid.EXPECT().RequestSubmarineSwap(gomock.Any(), gomock.Any(), gomock.Any()).Return(provider.SubmarineSwapResponse{
		SwapId:            "1234",
		InvoiceBTCAddress: "bcrt1q6zszlnxhlq0lsmfc42nkwgqedy9kvmvmxhkvme",
	}, nil).AnyTimes()

	//MockRequestSubmarineSwap
	mockProviderInvalid.EXPECT().MonitorSwap(gomock.Any(), gomock.Any(), gomock.Any()).Return(looprpc.SwapStatus{
		Amt:           0,
		Id:            "",
		IdBytes:       []byte{},
		Type:          0,
		State:         looprpc.SwapState_FAILED,
		FailureReason: looprpc.FailureReason_FAILURE_REASON_INSUFFICIENT_VALUE,
	}, nil).AnyTimes()

	//Mock get new wallet address
	mockNodeGuardClient.EXPECT().GetNewWalletAddress(gomock.Any(), gomock.Any()).Return(&nodeguard.GetNewWalletAddressResponse{
		Address: "bcrt1q6zszlnxhlq0lsmfc42nkwgqedy9kvmvmxhkvme",
	}, nil).AnyTimes()
	//Mock request withdrawal
	mockNodeGuardClient.EXPECT().RequestWithdrawal(gomock.Any(), gomock.Any()).Return(&nodeguard.RequestWithdrawalResponse{
		Txid:        "bd0d500cc43b8c60769fd480170ace6660f2881d69bef475e03210f7f8e80c6f",
		IsHotWallet: true,
	}, nil).AnyTimes()

	//Active channel
	channelActive := &lnrpc.Channel{
		Active:        true,
		ChanId:        123,
		Capacity:      1000,
		LocalBalance:  100,
		RemoteBalance: 900,
		RemotePubkey:  "03485d8dcdd149c87553eeb80586eb2bece874d412e9f117304446ce189955d375",
	}

	nodeInfo := lnrpc.GetInfoResponse{
		Alias: "Test",
	}

	tests := []struct {
		name    string
		args    ManageChannelLiquidityInfo
		wantErr bool
	}{
		{
			name: "Manage channel liquidity test valid reverse swap",
			args: ManageChannelLiquidityInfo{
				channel:             channelActive,
				channelBalanceRatio: 0.1,
				channelRules:        &[]nodeguard.LiquidityRule{{ChannelId: 123, NodePubkey: "", WalletId: 1, MinimumLocalBalance: 20, MinimumRemoteBalance: 80, RebalanceTarget: 60}},
				nodeguardClient:     mockNodeGuardClient,
				loopProvider:        mockProvider,
				loopdMacaroon:       "0201036c6e6402f801030a10dc64226b045d25f090b114baebcbf04c1201301a160a0761646472657373120472656164120577726974651a130a04696e666f120472656164120577726974651a170a08696e766f69636573120472656164120577726974651a210a086d616361726f6f6e120867656e6572617465120472656164120577726974651a160a076d657373616765120472656164120577726974651a170a086f6666636861696e120472656164120577726974651a160a076f6e636861696e120472656164120577726974651a140a057065657273120472656164120577726974651a180a067369676e6572120867656e657261746512047265616400000620a21b8cc8c071aa5104b706b751aede972f642537c05da31450fb4b02c6da776e",
				nodeInfo:            nodeInfo,
				ctx:                 context.TODO(),
			},
			wantErr: false,
		},
		{
			name: "Manage channel liquidity test valid swap",
			args: ManageChannelLiquidityInfo{
				channel:             channelActive,
				channelBalanceRatio: 0.9,
				channelRules: &[]nodeguard.LiquidityRule{
					{
						ChannelId:            123,
						NodePubkey:           "03485d8dcdd149c87553eeb80586eb2bece874d412e9f117304446ce189955d375",
						WalletId:             1,
						MinimumLocalBalance:  20,
						MinimumRemoteBalance: 80,
						RebalanceTarget:      60,
					},
				},
				nodeguardClient: mockNodeGuardClient,
				loopProvider:    mockProvider,
				loopdMacaroon:   "0201036c6e6402f801030a10dc64226b045d25f090b114baebcbf04c1201301a160a0761646472657373120472656164120577726974651a130a04696e666f120472656164120577726974651a170a08696e766f69636573120472656164120577726974651a210a086d616361726f6f6e120867656e6572617465120472656164120577726974651a160a076d657373616765120472656164120577726974651a170a086f6666636861696e120472656164120577726974651a160a076f6e636861696e120472656164120577726974651a140a057065657273120472656164120577726974651a180a067369676e6572120867656e657261746512047265616400000620a21b8cc8c071aa5104b706b751aede972f642537c05da31450fb4b02c6da776e",
				nodeInfo:        nodeInfo,
				ctx:             context.TODO(),
			},
			wantErr: false,
		},
		{
			name: "Manage channel liquidity test failed reverse swap",
			args: ManageChannelLiquidityInfo{
				channel:             channelActive,
				channelBalanceRatio: 0.1,
				channelRules:        &[]nodeguard.LiquidityRule{{ChannelId: 123, NodePubkey: "", WalletId: 1, MinimumLocalBalance: 20, MinimumRemoteBalance: 80, RebalanceTarget: 60}},
				nodeguardClient:     mockNodeGuardClient,
				loopProvider:        mockProviderInvalid,
				loopdMacaroon:       "0201036c6e6402f801030a10dc64226b045d25f090b114baebcbf04c1201301a160a0761646472657373120472656164120577726974651a130a04696e666f120472656164120577726974651a170a08696e766f69636573120472656164120577726974651a210a086d616361726f6f6e120867656e6572617465120472656164120577726974651a160a076d657373616765120472656164120577726974651a170a086f6666636861696e120472656164120577726974651a160a076f6e636861696e120472656164120577726974651a140a057065657273120472656164120577726974651a180a067369676e6572120867656e657261746512047265616400000620a21b8cc8c071aa5104b706b751aede972f642537c05da31450fb4b02c6da776e",
				nodeInfo:            nodeInfo,
				ctx:                 context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "Manage channel liquidity test failed swap",
			args: ManageChannelLiquidityInfo{
				channel:             channelActive,
				channelBalanceRatio: 0.9,
				channelRules:        &[]nodeguard.LiquidityRule{{ChannelId: 123, NodePubkey: "", WalletId: 1, MinimumLocalBalance: 20, MinimumRemoteBalance: 80, RebalanceTarget: 60}},
				nodeguardClient:     mockNodeGuardClient,
				loopProvider:        mockProviderInvalid,
				loopdMacaroon:       "0201036c6e6402f801030a10dc64226b045d25f090b114baebcbf04c1201301a160a0761646472657373120472656164120577726974651a130a04696e666f120472656164120577726974651a170a08696e766f69636573120472656164120577726974651a210a086d616361726f6f6e120867656e6572617465120472656164120577726974651a160a076d657373616765120472656164120577726974651a170a086f6666636861696e120472656164120577726974651a160a076f6e636861696e120472656164120577726974651a140a057065657273120472656164120577726974651a180a067369676e6572120867656e657261746512047265616400000620a21b8cc8c071aa5104b706b751aede972f642537c05da31450fb4b02c6da776e",
				nodeInfo:            nodeInfo,
				ctx:                 context.TODO(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := manageChannelLiquidity(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("manageChannelLiquidity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
