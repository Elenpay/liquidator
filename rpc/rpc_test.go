package rpc

import (
	"reflect"
	"testing"

	"github.com/lightninglabs/loop/looprpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func TestCreateLightningClient(t *testing.T) {
	type args struct {
		nodeEndpoint   string
		tlsCertEncoded string
	}
	tests := []struct {
		name    string
		args    args
		want    lnrpc.LightningClient
		want1   *grpc.ClientConn
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := CreateLightningClient(tt.args.nodeEndpoint, tt.args.tlsCertEncoded)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateLightningClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateLightningClient() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CreateLightningClient() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCreateSwapClientClient(t *testing.T) {
	type args struct {
		nodeEndpoint   string
		tlsCertEncoded string
	}
	tests := []struct {
		name    string
		args    args
		want    looprpc.SwapClientClient
		want1   *grpc.ClientConn
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := CreateSwapClientClient(tt.args.nodeEndpoint, tt.args.tlsCertEncoded)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSwapClientClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSwapClientClient() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CreateSwapClientClient() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getConn(t *testing.T) {
	type args struct {
		nodeEndpoint string
		creds        credentials.TransportCredentials
	}
	tests := []struct {
		name    string
		args    args
		want    *grpc.ClientConn
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getConn(tt.args.nodeEndpoint, tt.args.creds)
			if (err != nil) != tt.wantErr {
				t.Errorf("getConn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getConn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateCredentials(t *testing.T) {
	type args struct {
		tlsCertEncoded string
	}
	tests := []struct {
		name string
		args args
		want credentials.TransportCredentials
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateCredentials(tt.args.tlsCertEncoded); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}
