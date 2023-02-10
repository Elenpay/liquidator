package main

import (
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

// func Test_monitorChannels(t *testing.T) {

// 	//Arrange
// 	mockCtrl := gomock.NewController(t)
// 	defer mockCtrl.Finish()
// 	mockLightningClient := NewMockLightningClient(mockCtrl)
// 	mockLightningClient.EXPECT().ListChannels(gomock.Any(), gomock.Any()).Return(&lnrpc.ListChannelsResponse{
// 		Channels: []*lnrpc.Channel{
// 			{
// 				ChanId:        1,
// 				LocalBalance:  100,
// 				RemoteBalance: 900,
// 				Capacity:      1000,
// 			},
// 		},
// 	}, nil)

// 	//Act

// 	monitorChannels("localhost:5001", "macaroon", mockLightningClient, context.TODO())

// 	//Assert

// 	//Assert that the local balance is 100/900

// 	// metric := prometheusMetrics.channelBalanceGauge.With(prometheus.Labels{"channel_id": "1"})

// 	// t.Log(metric)

// }

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
			actualBalance, err := getChannelBalanceRatio(tt.args.channel)
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
	type args struct {
		channel             *lnrpc.Channel
		channelBalanceRatio float64
		channelRules        *[]nodeguard.LiquidityRule
		swapClientClient    looprpc.SwapClientClient
		nodeguardClient     nodeguard.NodeGuardServiceClient
		loopProvider        *provider.LoopProvider
	}

	//gomock controller
	mockCtrl := gomock.NewController(t)
	

	//Mock swap client
	swapClient := provider.NewMockSwapClientClient(mockCtrl)

	//Mock nodeguard client
	

	//Active channel
	channelActive := &lnrpc.Channel{
		Active:                true,
		ChanId:                123,
		Capacity:              1000,
		LocalBalance:          100,
		RemoteBalance:         900,
	}

	//Unactive channel
	channelUnactive := &lnrpc.Channel{
		Active:                false,
		ChanId:                123,	
		Capacity:              1000,
		LocalBalance:          100,
		RemoteBalance:         900,
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Manage channel liquidity test valid",
			args:    args{
				channel:             channelActive,
				channelBalanceRatio: 0.1,
				channelRules:        &[]nodeguard.LiquidityRule{
					{
						ChannelId:            123,
						NodePubkey:           "",
						WalletId:             1,
						MinimumLocalBalance:  0.2,
						MinimumRemoteBalance: 0.8,
					},
				},
				swapClientClient:    mockSwapClient,
				nodeguardClient:     mockNodeGuardClient,
				loopProvider:        &provider.LoopProvider{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := manageChannelLiquidity(tt.args.channel, tt.args.channelBalanceRatio, tt.args.channelRules, tt.args.swapClientClient, tt.args.nodeguardClient, tt.args.loopProvider); (err != nil) != tt.wantErr {
				t.Errorf("manageChannelLiquidity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

