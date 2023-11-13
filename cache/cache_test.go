package cache

import (
	"reflect"
	"testing"

	"github.com/Elenpay/liquidator/nodeguard"
)

func TestNewCache(t *testing.T) {
	tests := []struct {
		name    string
		want    Cache
		wantErr bool
	}{
		{
			name:    "TestNewCache",
			want:    Cache(&BigCache{}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCache()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//Check that got is not nil and set up
			if got == nil {
				t.Errorf("NewCache() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigCache_SetGetLiquidityRules(t *testing.T) {

	cache, err := NewCache()
	if err != nil {
		t.Fatalf("NewCache() error = %v, wantErr %v", err, false)
	}

	type args struct {
		rules    []nodeguard.LiquidityRule
		mapRules map[uint64][]nodeguard.LiquidityRule
	}
	pubKey := "010101"
	x := nodeguard.LiquidityRule{
		ChannelId:            1,
		NodePubkey:           "",
		SwapWalletId:         0,
		ReverseSwapWalletId:  nil,
		MinimumLocalBalance:  0,
		MinimumRemoteBalance: 0,
	}

	y := nodeguard.LiquidityRule{
		ChannelId:            2,
		NodePubkey:           "",
		SwapWalletId:         0,
		ReverseSwapWalletId:  nil,
		MinimumLocalBalance:  0,
		MinimumRemoteBalance: 0,
	}
	tests := []struct {
		name    string
		c       *BigCache
		args    args
		wantErr bool
	}{
		{
			name: "TestSetLiquidityRules_Positive1",
			c:    cache.(*BigCache),
			args: args{
				rules: []nodeguard.LiquidityRule{x},
				mapRules: map[uint64][]nodeguard.LiquidityRule{
					1: []nodeguard.LiquidityRule{x},
				},
			},
			wantErr: false,
		},
		{
			name: "TestSetLiquidityRules_Positive2",
			c:    cache.(*BigCache),
			args: args{
				rules: []nodeguard.LiquidityRule{
					x,
					y},
				mapRules: map[uint64][]nodeguard.LiquidityRule{
					1: []nodeguard.LiquidityRule{x},
					2: []nodeguard.LiquidityRule{y}},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.SetLiquidityRules(pubKey, tt.args.rules); (err != nil) != tt.wantErr {
				t.Errorf("BigCache.SetLiquidityRules() error = %v, wantErr %v", err, tt.wantErr)
			}
			//Check that the channelId is 1 and liquidity rule is set
			rules, err := tt.c.GetLiquidityRules(pubKey)

			if err != nil {
				t.Errorf("BigCache.SetLiquidityRules() error = %v, wantErr %v", err, tt.wantErr)
			}

			//Check that the rules slice is the same as the args
			if !reflect.DeepEqual(rules, tt.args.mapRules) {
				t.Errorf("BigCache.SetLiquidityRules() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
