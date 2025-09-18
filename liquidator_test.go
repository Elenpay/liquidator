package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Elenpay/liquidator/nodeguard"
	"github.com/Elenpay/liquidator/provider"
	"github.com/lightninglabs/loop/looprpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"go.uber.org/goleak"
	gomock "go.uber.org/mock/gomock"
)

// Tear up method
func TestMain(m *testing.M) {

	log.SetLevel(log.DebugLevel)

	//Tear up
	initMetrics(prometheus.NewRegistry())

	//Run tests and verify goroutine leaks
	goleak.VerifyTestMain(m)

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
				LocalBalance:  0,
				RemoteBalance: 0,
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
	mockNodeGuardClient.EXPECT().GetNewWalletAddress(gomock.Any(), gomock.Any()).Return(&nodeguard.GetNewWalletAddressResponse{
		Address: "bcrt1q6zszlnxhlq0lsmfc42nkwgqedy9kvmvmxhkvme",
	}, nil).AnyTimes()

	mockNodeGuardClient.EXPECT().RequestWithdrawal(gomock.Any(), gomock.Any()).Return(&nodeguard.RequestWithdrawalResponse{
		Txid:        "bd0d500cc43b8c60769fd480170ace6660f2881d69bef475e03210f7f8e80c6f",
		IsHotWallet: true,
	}, nil).AnyTimes()

	//Mock provider valid swaps
	mockProvider := createMockProviderValidSwap(mockCtrl)

	//Mock provider invalid swap
	mockProviderInvalid := createMockProviderInvalidSwap(mockCtrl)

	//Mock provider invalid reverse swap loop error
	mockProviderInvalidLoopError := createMockProviderInvalidSwapLoopError(mockCtrl)

	// Wallet and address for reverse swaps
	var walletId int32 = 1
	var address string = "bcrt1q6zszlnxhlq0lsmfc42nkwgqedy9kvmvmxhkvme"

	//Active channel
	channelActive := &lnrpc.Channel{
		Active:        true,
		ChanId:        123,
		Capacity:      1000,
		LocalBalance:  100,
		RemoteBalance: 900,
		RemotePubkey:  "03485d8dcdd149c87553eeb80586eb2bece874d412e9f117304446ce189955d375",
		LocalConstraints: &lnrpc.ChannelConstraints{
			CsvDelay:          144,
			ChanReserveSat:    250,
			DustLimitSat:      300,
			MaxPendingAmtMsat: 300 * 1000,
			MinHtlcMsat:       350 * 1000,
			MaxAcceptedHtlcs:  30,
		},
	}

	//Mock for lightning clients for invoice rebalance
	lightningClient := NewMockLightningClient(mockCtrl)
	lightningClient.EXPECT().AddInvoice(gomock.Any(), gomock.Any()).Return(&lnrpc.AddInvoiceResponse{}, nil).AnyTimes()
	lightningClient.EXPECT().SendPaymentSync(gomock.Any(), gomock.Any()).Return(&lnrpc.SendResponse{}, nil).AnyTimes()

	nodeInfo := lnrpc.GetInfoResponse{
		Alias:          "Test",
		IdentityPubkey: "1",
	}

	tests := []struct {
		name    string
		args    ManageChannelLiquidityInfo
		wantErr bool
	}{

		{
			name: "Manage channel liquidity test valid reverse swap",
			args: ManageChannelLiquidityInfo{
				BaseInfo: BaseInfo{
					nodeHost:         "",
					nodeInfo:         nodeInfo,
					nodeMacaroon:     "",
					loopdMacaroon:    "1",
					lightningClients: map[string]lnrpc.LightningClient{},
					nodeguardClient:  mockNodeGuardClient,
					swapClient:       nil,
					provider:         mockProvider,
				},
				channel:             channelActive,
				channelBalanceRatio: 0.1,
				channelRules:        &[]nodeguard.LiquidityRule{{ChannelId: 123, NodePubkey: "", IsReverseSwapWalletRule: true, ReverseSwapWalletId: &walletId, MinimumLocalBalance: 20, MinimumRemoteBalance: 80, RebalanceTarget: 60}},
				ctx:                 context.TODO(),
			},
			wantErr: false,
		},
		{
			name: "Manage channel liquidity test valid reverse swap",
			args: ManageChannelLiquidityInfo{
				BaseInfo: BaseInfo{
					nodeHost:         "",
					nodeInfo:         nodeInfo,
					nodeMacaroon:     "",
					loopdMacaroon:    "1",
					lightningClients: map[string]lnrpc.LightningClient{},
					nodeguardClient:  mockNodeGuardClient,
					swapClient:       nil,
					provider:         mockProvider,
				},
				channel:             channelActive,
				channelBalanceRatio: 0.1,
				channelRules:        &[]nodeguard.LiquidityRule{{ChannelId: 123, NodePubkey: "", IsReverseSwapWalletRule: false, ReverseSwapAddress: &address, MinimumLocalBalance: 20, MinimumRemoteBalance: 80, RebalanceTarget: 60}},
				ctx:                 context.TODO(),
			},
			wantErr: false,
		},
		{
			name: "Manage channel liquidity test valid reverse swap bypassing max pending amt",
			args: ManageChannelLiquidityInfo{
				BaseInfo: BaseInfo{
					nodeHost:         "",
					nodeInfo:         nodeInfo,
					nodeMacaroon:     "",
					loopdMacaroon:    "1",
					lightningClients: map[string]lnrpc.LightningClient{},
					nodeguardClient:  mockNodeGuardClient,
					swapClient:       nil,
					provider:         mockProvider,
				},
				channel:             channelActive,
				channelBalanceRatio: 0.1,
				channelRules:        &[]nodeguard.LiquidityRule{{ChannelId: 123, NodePubkey: "", IsReverseSwapWalletRule: true, ReverseSwapWalletId: &walletId, MinimumLocalBalance: 20, MinimumRemoteBalance: 80, RebalanceTarget: 40}},
				ctx:                 context.TODO(),
			},
			wantErr: false,
		},
		{
			name: "Manage channel liquidity test valid swap",
			args: ManageChannelLiquidityInfo{
				BaseInfo: BaseInfo{
					nodeHost:         "",
					nodeInfo:         nodeInfo,
					nodeMacaroon:     "",
					loopdMacaroon:    "1",
					lightningClients: map[string]lnrpc.LightningClient{},
					nodeguardClient:  mockNodeGuardClient,
					swapClient:       nil,
					provider:         mockProvider,
				},
				channel:             channelActive,
				channelBalanceRatio: 0.9,
				ctx:                 context.TODO(),
				channelRules: &[]nodeguard.LiquidityRule{
					{
						ChannelId:            123,
						NodePubkey:           "03485d8dcdd149c87553eeb80586eb2bece874d412e9f117304446ce189955d375",
						SwapWalletId:         1,
						MinimumLocalBalance:  20,
						MinimumRemoteBalance: 80,
						RebalanceTarget:      60,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Manage channel liquidity test failed reverse swap",
			args: ManageChannelLiquidityInfo{
				BaseInfo: BaseInfo{
					nodeHost:         "",
					nodeInfo:         nodeInfo,
					nodeMacaroon:     "",
					loopdMacaroon:    "1",
					lightningClients: map[string]lnrpc.LightningClient{},
					nodeguardClient:  mockNodeGuardClient,
					swapClient:       nil,
					provider:         mockProviderInvalid,
				},
				channel:             channelActive,
				channelBalanceRatio: 0.1,
				channelRules:        &[]nodeguard.LiquidityRule{{ChannelId: 123, NodePubkey: "", IsReverseSwapWalletRule: true, ReverseSwapWalletId: &walletId, MinimumLocalBalance: 20, MinimumRemoteBalance: 80, RebalanceTarget: 60}},
				ctx:                 context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "Manage channel liquidity test failed swap",
			args: ManageChannelLiquidityInfo{
				BaseInfo: BaseInfo{
					nodeHost:         "",
					nodeInfo:         nodeInfo,
					nodeMacaroon:     "",
					loopdMacaroon:    "1",
					lightningClients: map[string]lnrpc.LightningClient{},
					nodeguardClient:  mockNodeGuardClient,
					swapClient:       nil,
					provider:         mockProviderInvalid,
				},
				channel:             channelActive,
				channelBalanceRatio: 0.9,
				channelRules:        &[]nodeguard.LiquidityRule{{ChannelId: 123, NodePubkey: "", SwapWalletId: 1, MinimumLocalBalance: 20, MinimumRemoteBalance: 80, RebalanceTarget: 60}},
				ctx:                 context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "Manage channel liquidity test failed reverse swap channel outbound capacity",
			args: ManageChannelLiquidityInfo{
				BaseInfo: BaseInfo{
					nodeHost:         "",
					nodeInfo:         nodeInfo,
					nodeMacaroon:     "",
					loopdMacaroon:    "1",
					lightningClients: map[string]lnrpc.LightningClient{},
					nodeguardClient:  mockNodeGuardClient,
					swapClient:       nil,
					provider:         mockProviderInvalidLoopError,
				},
				channel:             channelActive,
				channelBalanceRatio: 0.1,
				channelRules:        &[]nodeguard.LiquidityRule{{ChannelId: 123, NodePubkey: "", IsReverseSwapWalletRule: true, ReverseSwapWalletId: &walletId, MinimumLocalBalance: 20, MinimumRemoteBalance: 80, RebalanceTarget: 60}},
				ctx:                 context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "Manage channel liquidity test valid invoice rebalance minimum local balance threshold",
			args: ManageChannelLiquidityInfo{
				BaseInfo: BaseInfo{
					nodeInfo:         nodeInfo,
					lightningClients: map[string]lnrpc.LightningClient{"1": lightningClient, "2": lightningClient},
					nodeMacaroon:     "1",
					loopdMacaroon:    "2",
				},
				channel:             channelActive,
				channelBalanceRatio: 0.1,
				channelRules: &[]nodeguard.LiquidityRule{{
					NodePubkey:           "1",
					MinimumLocalBalance:  20,
					MinimumRemoteBalance: 80,
					RebalanceTarget:      50,
					RemoteNodePubkey:     "2",
				}},
				ctx: context.TODO(),
			},
			wantErr: false,
		},
		{
			name: "Manage channel liquidity test valid invoice rebalance remote local balance threshold",
			args: ManageChannelLiquidityInfo{
				BaseInfo: BaseInfo{
					nodeInfo:         nodeInfo,
					lightningClients: map[string]lnrpc.LightningClient{"1": lightningClient, "2": lightningClient},
					nodeMacaroon:     "1",
					loopdMacaroon:    "2",
				},
				channel:             channelActive,
				channelBalanceRatio: 0.9,
				channelRules: &[]nodeguard.LiquidityRule{{
					NodePubkey:           "1",
					MinimumLocalBalance:  20,
					MinimumRemoteBalance: 80,
					RebalanceTarget:      50,
					RemoteNodePubkey:     "2",
				}},
				ctx: context.TODO(),
			},
			wantErr: false,
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

func createMockProviderValidSwap(mockCtrl *gomock.Controller) *provider.MockProvider {
	mockProvider := provider.NewMockProvider(mockCtrl)

	mockProvider.EXPECT().RequestReverseSubmarineSwap(gomock.Any(), gomock.Any(), gomock.Any()).Return(provider.ReverseSubmarineSwapResponse{
		SwapId: "1234",
	}, nil).AnyTimes()

	mockProvider.EXPECT().RequestSubmarineSwap(gomock.Any(), gomock.Any(), gomock.Any()).Return(provider.SubmarineSwapResponse{
		SwapId:            "1234",
		InvoiceBTCAddress: "bcrt1q6zszlnxhlq0lsmfc42nkwgqedy9kvmvmxhkvme",
	}, nil).AnyTimes()

	mockProvider.EXPECT().MonitorSwap(gomock.Any(), gomock.Any(), gomock.Any()).Return(&looprpc.SwapStatus{
		Amt:              0,
		Id:               "",
		IdBytes:          []byte{},
		Type:             0,
		State:            looprpc.SwapState_SUCCESS,
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
	}, nil).AnyTimes()
	return mockProvider
}

func createMockProviderInvalidSwap(mockCtrl *gomock.Controller) *provider.MockProvider {
	mockProviderInvalid := provider.NewMockProvider(mockCtrl)

	mockProviderInvalid.EXPECT().RequestReverseSubmarineSwap(gomock.Any(), gomock.Any(), gomock.Any()).Return(provider.ReverseSubmarineSwapResponse{
		SwapId: "1234",
	}, nil).AnyTimes()

	mockProviderInvalid.EXPECT().RequestSubmarineSwap(gomock.Any(), gomock.Any(), gomock.Any()).Return(provider.SubmarineSwapResponse{
		SwapId:            "1234",
		InvoiceBTCAddress: "bcrt1q6zszlnxhlq0lsmfc42nkwgqedy9kvmvmxhkvme",
	}, nil).AnyTimes()

	mockProviderInvalid.EXPECT().MonitorSwap(gomock.Any(), gomock.Any(), gomock.Any()).Return(&looprpc.SwapStatus{
		Amt:           0,
		Id:            "",
		IdBytes:       []byte{},
		Type:          0,
		State:         looprpc.SwapState_FAILED,
		FailureReason: looprpc.FailureReason_FAILURE_REASON_OFFCHAIN,
	}, nil).AnyTimes()

	return mockProviderInvalid

}

func createMockProviderInvalidSwapLoopError(mockCtrl *gomock.Controller) *provider.MockProvider {
	mockProviderInvalid := provider.NewMockProvider(mockCtrl)

	mockProviderInvalid.EXPECT().RequestReverseSubmarineSwap(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(
			provider.ReverseSubmarineSwapResponse{},
			errors.New("code = Unknown desc = channel balance too low for loop out amount: Requested swap amount of 2471564 sats along with the maximum routing fee of 49441 sats is more than what can be routed given current state of the channel set"),
		).AnyTimes()

	return mockProviderInvalid
}

func Test_monitorChannel(t *testing.T) {

	mockCtrl := gomock.NewController(t)

	type args struct {
		info       MonitorChannelInfo
		iterations int
	}

	channel := &lnrpc.Channel{
		Active:        true,
		RemotePubkey:  "1",
		ChannelPoint:  "",
		ChanId:        1,
		Capacity:      1000,
		LocalBalance:  900,
		RemoteBalance: 100,
	}

	//channel with htlcs
	channelHtlcs := &lnrpc.Channel{
		Active:        true,
		RemotePubkey:  "1",
		ChannelPoint:  "",
		ChanId:        1,
		Capacity:      1000,
		LocalBalance:  900,
		RemoteBalance: 100,
		PendingHtlcs: []*lnrpc.HTLC{
			{
				Incoming:            false,
				Amount:              44,
				HashLock:            []byte{},
				ExpirationHeight:    0,
				HtlcIndex:           0,
				ForwardingChannel:   0,
				ForwardingHtlcIndex: 0,
			},
		},
	}

	mockLightningClient := NewMockLightningClient(mockCtrl)

	mockLightningClient.EXPECT().GetInfo(gomock.Any(), gomock.Any()).Return(&lnrpc.GetInfoResponse{
		IdentityPubkey: "1",
	}, nil).AnyTimes()

	mockLightningClient.EXPECT().GetNodeInfo(gomock.Any(), gomock.Any()).Return(&lnrpc.NodeInfo{
		Node: &lnrpc.LightningNode{
			LastUpdate: 0,
			PubKey:     "1",
		},
	}, nil).AnyTimes()

	mockLightningClient.EXPECT().ListChannels(gomock.Any(), gomock.Any()).Return(&lnrpc.ListChannelsResponse{
		Channels: []*lnrpc.Channel{
			channel,
		},
	}, nil).AnyTimes()

	// Wallet id for reverse swaps
	var walletId *int32 = new(int32)
	*walletId = 1

	// Liquidity rules for the channel
	liquidityRules := map[uint64][]nodeguard.LiquidityRule{
		channel.ChanId: {
			{
				ChannelId:               channel.ChanId,
				NodePubkey:              "",
				SwapWalletId:            1,
				IsReverseSwapWalletRule: true,
				ReverseSwapWalletId:     walletId,
				MinimumLocalBalance:     20,
				MinimumRemoteBalance:    80,
				RebalanceTarget:         50,
			},
		},
	}

	//Mock nodeguard client
	nodeguardClient := nodeguard.NewMockNodeGuardServiceClient(mockCtrl)

	//Mock GetNewWalletAddress
	nodeguardClient.EXPECT().GetNewWalletAddress(gomock.Any(), gomock.Any()).Return(&nodeguard.GetNewWalletAddressResponse{
		Address: "bcrt1q6zszlnxhlq0lsmfc42nkwgqedy9kvmvmxhkvme",
	}, nil).AnyTimes()

	//lightning client mock
	mockLightningClientWithHTLCs := NewMockLightningClient(mockCtrl)

	//Mock list channels with htlcs
	mockLightningClientWithHTLCs.EXPECT().ListChannels(gomock.Any(), gomock.Any()).Return(&lnrpc.ListChannelsResponse{
		Channels: []*lnrpc.Channel{
			channelHtlcs,
		},
	}, nil).AnyTimes()

	mockLightningClientWithHTLCs.EXPECT().GetInfo(gomock.Any(), gomock.Any()).Return(&lnrpc.GetInfoResponse{
		IdentityPubkey: "1",
	}, nil).AnyTimes()

	mockLightningClientWithHTLCs.EXPECT().GetNodeInfo(gomock.Any(), gomock.Any()).Return(&lnrpc.NodeInfo{
		Node: &lnrpc.LightningNode{
			LastUpdate: 0,
			PubKey:     "1",
		},
	}, nil).AnyTimes()

	nodeInfo := lnrpc.GetInfoResponse{IdentityPubkey: "1"}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Monitor channel goroutine leak test invalid swap",
			args: args{
				info: MonitorChannelInfo{
					BaseInfo: BaseInfo{
						nodeHost:         "",
						lightningClients: map[string]lnrpc.LightningClient{"1": mockLightningClient},
						swapClient:       provider.NewMockSwapClientClient(mockCtrl),
						nodeguardClient:  nodeguardClient,
						provider:         createMockProviderInvalidSwap(mockCtrl),
						nodeInfo:         nodeInfo},
					channel:          channel,
					context:          context.TODO(),
					liquidationRules: liquidityRules,
				},
				iterations: 4,
			},
		},
		{
			name: "Monitor channel with ongoing htlc",
			args: args{
				info: MonitorChannelInfo{
					BaseInfo: BaseInfo{
						nodeHost:         "",
						lightningClients: map[string]lnrpc.LightningClient{"1": mockLightningClientWithHTLCs},
						swapClient:       provider.NewMockSwapClientClient(mockCtrl),
						nodeguardClient:  nodeguardClient,
						provider:         &provider.LoopProvider{},
						nodeInfo:         nodeInfo},
					channel:          channelHtlcs,
					context:          context.TODO(),
					liquidationRules: map[uint64][]nodeguard.LiquidityRule{},
				},
				iterations: 4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for i := 0; i < tt.args.iterations; i++ {
				go monitorChannel(tt.args.info)

			}

		})
	}

	// Wait for the goroutine to finish
	time.Sleep(1 * time.Second)
}
