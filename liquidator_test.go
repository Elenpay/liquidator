package main

import (
	//"context"
	"testing"

	//gomock "github.com/golang/mock/gomock"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/prometheus/client_golang/prometheus"
)

// Tear up method
func TestMain(m *testing.M) {
	//Tear up
	InitMetrics(prometheus.NewRegistry())

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
			0.1, false}},
		{"Test 2 positive", args{
			&lnrpc.Channel{
				ChanId:        1,
				LocalBalance:  1000,
				RemoteBalance: 900,
				Capacity:      1000,
			},
			1, false}},
		{"Test 3 negative", args{
			&lnrpc.Channel{
				ChanId:        1,
				LocalBalance:  1000,
				RemoteBalance: 900,
				Capacity:      0,
			},
			-1, true}},
		{"Test 4 negative", args{
			&lnrpc.Channel{
				ChanId:        1,
				LocalBalance:  1100,
				RemoteBalance: 100,
				Capacity:      1000,
			},
			-1, true},
		},
	}

	for _, tt := range tests {
		t.Logf("Running test: %v", tt.name)

		t.Run(tt.name, func(t *testing.T) {
			actualBalance, err := recordChannelBalance(tt.args.channel)
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
