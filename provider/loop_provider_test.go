package provider

import (
	"context"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/lightninglabs/loop/looprpc"
	"github.com/spf13/viper"
	gomock "go.uber.org/mock/gomock"
)

func TestMain(m *testing.M) {

	setLimitFees()

	m.Run()
}

func setLimitFees() {
	viper.Set("limitQuoteFees", "0.005")
	viper.Set("limitFeesL2", "0.002")
	viper.Set("swapLockTimeout", "30s") // Set a short timeout for testing
}

func TestLoopProvider_RequestSubmarineSwap(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
					ChannelId:          12345,
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
			want: nil,

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
		want    *looprpc.SwapStatus
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
			want:    swapStatusSuccess,
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
			want:    swapStatusFailure,
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

// TestLoopProvider_LockFunctionality tests the submarine swap and reverse swap locking mechanisms
func TestLoopProvider_LockFunctionality(t *testing.T) {

	// Test submarine swap lock functionality
	t.Run("SubmarineSwapLock", func(t *testing.T) {
		l := &LoopProvider{}
		testChannelId := uint64(12345)

		// Test acquiring lock for the first time
		err := l.acquireSubmarineSwapLock(testChannelId)
		if err != nil {
			t.Errorf("Expected no error when acquiring submarine swap lock for the first time, got: %v", err)
		}

		// Test that lock is active - should fail to acquire again
		err = l.acquireSubmarineSwapLock(testChannelId)
		if err == nil {
			t.Error("Expected error when trying to acquire submarine swap lock while already locked")
		}

		// Verify the error message contains rate limit info
		if err != nil && !contains(err.Error(), "submarine swap is locked") {
			t.Errorf("Expected rate limit error message, got: %v", err)
		}

		// Note: In production, lock only expires via timeout, no manual release for anti-spam
	})

	// Test reverse swap lock functionality
	t.Run("ReverseSwapLock", func(t *testing.T) {
		l := &LoopProvider{}
		testChannelId := uint64(67890)

		// Test acquiring lock for the first time
		err := l.acquireReverseSwapLock(testChannelId)
		if err != nil {
			t.Errorf("Expected no error when acquiring reverse swap lock for the first time, got: %v", err)
		}

		// Test that lock is active - should fail to acquire again
		err = l.acquireReverseSwapLock(testChannelId)
		if err == nil {
			t.Error("Expected error when trying to acquire reverse swap lock while already locked")
		}

		// Verify the error message contains rate limit info
		if err != nil && !contains(err.Error(), "reverse submarine swap is locked") {
			t.Errorf("Expected rate limit error message, got: %v", err)
		}

		// Note: In production, lock only expires via timeout, no manual release for anti-spam
	})

	// Test that submarine and reverse swap locks are independent
	t.Run("IndependentLocks", func(t *testing.T) {
		l := &LoopProvider{}
		testChannelId1 := uint64(11111)
		testChannelId2 := uint64(22222)

		// Acquire submarine swap lock
		err := l.acquireSubmarineSwapLock(testChannelId1)
		if err != nil {
			t.Errorf("Expected no error when acquiring submarine swap lock, got: %v", err)
		}

		// Should still be able to acquire reverse swap lock on different channel
		err = l.acquireReverseSwapLock(testChannelId2)
		if err != nil {
			t.Errorf("Expected no error when acquiring reverse swap lock while submarine swap is locked, got: %v", err)
		}

		// Note: In production, locks only expire via timeout, no manual release for anti-spam
	})
}

// TestLoopProvider_LockTimeout tests the automatic timeout functionality
func TestLoopProvider_LockTimeout(t *testing.T) {
	// Set a very short timeout for this test
	originalTimeout := viper.GetDuration("swapLockTimeout")
	viper.Set("swapLockTimeout", "100ms")
	defer viper.Set("swapLockTimeout", originalTimeout)

	t.Run("SubmarineSwapLockTimeout", func(t *testing.T) {
		l := &LoopProvider{}
		testChannelId := uint64(33333)

		// Acquire lock
		err := l.acquireSubmarineSwapLock(testChannelId)
		if err != nil {
			t.Errorf("Expected no error when acquiring submarine swap lock, got: %v", err)
		}

		// Immediately try to acquire again - should fail
		err = l.acquireSubmarineSwapLock(testChannelId)
		if err == nil {
			t.Error("Expected error when trying to acquire submarine swap lock while already locked")
		}

		// Wait for timeout
		time.Sleep(150 * time.Millisecond)

		// Should be able to acquire again after timeout
		err = l.acquireSubmarineSwapLock(testChannelId)
		if err != nil {
			t.Errorf("Expected no error when acquiring submarine swap lock after timeout, got: %v", err)
		}

		// Note: In production, locks only expire via timeout, no manual cleanup needed
	})

	t.Run("ReverseSwapLockTimeout", func(t *testing.T) {
		l := &LoopProvider{}
		testChannelId := uint64(44444)

		// Acquire lock
		err := l.acquireReverseSwapLock(testChannelId)
		if err != nil {
			t.Errorf("Expected no error when acquiring reverse swap lock, got: %v", err)
		}

		// Immediately try to acquire again - should fail
		err = l.acquireReverseSwapLock(testChannelId)
		if err == nil {
			t.Error("Expected error when trying to acquire reverse swap lock while already locked")
		}

		// Wait for timeout
		time.Sleep(150 * time.Millisecond)

		// Should be able to acquire again after timeout
		err = l.acquireReverseSwapLock(testChannelId)
		if err != nil {
			t.Errorf("Expected no error when acquiring reverse swap lock after timeout, got: %v", err)
		}

		// Note: In production, locks only expire via timeout, no manual cleanup needed
	})
}

// TestLoopProvider_ConcurrentLockAccess tests thread safety of the lock mechanism
func TestLoopProvider_ConcurrentLockAccess(t *testing.T) {
	l := &LoopProvider{}

	t.Run("ConcurrentSubmarineSwapLock", func(t *testing.T) {
		testChannelId := uint64(55555)
		successCount := 0
		errorCount := 0
		var mu sync.Mutex
		var wg sync.WaitGroup

		// Try to acquire the lock from multiple goroutines
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := l.acquireSubmarineSwapLock(testChannelId)
				mu.Lock()
				if err == nil {
					successCount++
				} else {
					errorCount++
				}
				mu.Unlock()
			}()
		}

		wg.Wait()

		// Only one should succeed, the rest should fail
		if successCount != 1 {
			t.Errorf("Expected exactly 1 successful lock acquisition, got: %d", successCount)
		}
		if errorCount != 9 {
			t.Errorf("Expected exactly 9 failed lock acquisitions, got: %d", errorCount)
		}

		// Note: In production, locks only expire via timeout, no manual cleanup needed
	})
}

// TestLoopProvider_PerChannelLocking tests the per-channel locking behavior
func TestLoopProvider_PerChannelLocking(t *testing.T) {

	t.Run("DifferentChannelsCanSwapSimultaneously", func(t *testing.T) {
		l := &LoopProvider{}
		channel1 := uint64(12345)
		channel2 := uint64(67890)
		channel3 := uint64(11111)
		channel4 := uint64(22222)

		// Different channels should be able to acquire submarine swap locks
		err1 := l.acquireSubmarineSwapLock(channel1)
		if err1 != nil {
			t.Errorf("Expected no error for channel1 submarine swap lock, got: %v", err1)
		}

		err2 := l.acquireSubmarineSwapLock(channel2)
		if err2 != nil {
			t.Errorf("Expected no error for channel2 submarine swap lock, got: %v", err2)
		}

		// Different channels should be able to acquire reverse swap locks
		err3 := l.acquireReverseSwapLock(channel3)
		if err3 != nil {
			t.Errorf("Expected no error for channel3 reverse swap lock, got: %v", err3)
		}

		err4 := l.acquireReverseSwapLock(channel4)
		if err4 != nil {
			t.Errorf("Expected no error for channel4 reverse swap lock, got: %v", err4)
		}
	})

	t.Run("SameChannelCannotHaveDuplicateSubmarineSwaps", func(t *testing.T) {
		l := &LoopProvider{}
		channelId := uint64(99999)

		// First submarine swap should succeed
		err1 := l.acquireSubmarineSwapLock(channelId)
		if err1 != nil {
			t.Errorf("Expected no error for first submarine swap lock, got: %v", err1)
		}

		// Second submarine swap on same channel should fail
		err2 := l.acquireSubmarineSwapLock(channelId)
		if err2 == nil {
			t.Error("Expected error for second submarine swap lock on same channel")
		}

		// Error should mention the channel ID
		if !contains(err2.Error(), fmt.Sprintf("channel %d", channelId)) {
			t.Errorf("Expected error message to contain channel ID %d, got: %s", channelId, err2.Error())
		}
	})

	t.Run("SameChannelCannotHaveDuplicateReverseSwaps", func(t *testing.T) {
		l := &LoopProvider{}
		channelId := uint64(88888)

		// First reverse swap should succeed
		err1 := l.acquireReverseSwapLock(channelId)
		if err1 != nil {
			t.Errorf("Expected no error for first reverse swap lock, got: %v", err1)
		}

		// Second reverse swap on same channel should fail
		err2 := l.acquireReverseSwapLock(channelId)
		if err2 == nil {
			t.Error("Expected error for second reverse swap lock on same channel")
		}

		// Error should mention the channel ID
		if !contains(err2.Error(), fmt.Sprintf("channel %d", channelId)) {
			t.Errorf("Expected error message to contain channel ID %d, got: %s", channelId, err2.Error())
		}
	})

	t.Run("SameChannelCannotHaveBothSubmarineAndReverseSwaps", func(t *testing.T) {
		l := &LoopProvider{}
		channelId := uint64(77777)

		// First submarine swap should succeed
		err1 := l.acquireSubmarineSwapLock(channelId)
		if err1 != nil {
			t.Errorf("Expected no error for submarine swap lock, got: %v", err1)
		}

		// Reverse swap on same channel should fail
		err2 := l.acquireReverseSwapLock(channelId)
		if err2 == nil {
			t.Error("Expected error for reverse swap lock on same channel as submarine swap")
		}

		// Error should mention conflicting swap
		if !contains(err2.Error(), "submarine swap") {
			t.Errorf("Expected error message to mention submarine swap conflict, got: %s", err2.Error())
		}
	})

	t.Run("ReverseSwapBlocksSubmarineSwap", func(t *testing.T) {
		l := &LoopProvider{}
		channelId := uint64(66666)

		// First reverse swap should succeed
		err1 := l.acquireReverseSwapLock(channelId)
		if err1 != nil {
			t.Errorf("Expected no error for reverse swap lock, got: %v", err1)
		}

		// Submarine swap on same channel should fail
		err2 := l.acquireSubmarineSwapLock(channelId)
		if err2 == nil {
			t.Error("Expected error for submarine swap lock on same channel as reverse swap")
		}

		// Error should mention conflicting swap
		if !contains(err2.Error(), "reverse submarine swap") {
			t.Errorf("Expected error message to mention reverse submarine swap conflict, got: %s", err2.Error())
		}
	})

	t.Run("ChannelLockTimeout", func(t *testing.T) {
		// Set a short timeout for testing
		originalTimeout := viper.GetString("swapLockTimeout")
		viper.Set("swapLockTimeout", "50ms")
		defer viper.Set("swapLockTimeout", originalTimeout)

		l := &LoopProvider{}
		channelId := uint64(55555)

		// Acquire lock
		err1 := l.acquireSubmarineSwapLock(channelId)
		if err1 != nil {
			t.Errorf("Expected no error for initial lock, got: %v", err1)
		}

		// Immediate retry should fail
		err2 := l.acquireSubmarineSwapLock(channelId)
		if err2 == nil {
			t.Error("Expected error for immediate retry")
		}

		// Wait for timeout
		time.Sleep(100 * time.Millisecond)

		// Should be able to acquire again after timeout
		err3 := l.acquireSubmarineSwapLock(channelId)
		if err3 != nil {
			t.Errorf("Expected no error after timeout, got: %v", err3)
		}
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
