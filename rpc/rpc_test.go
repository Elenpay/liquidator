package rpc

import (
	"testing"
)

func getTlsCertEncoded() string {
	return "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUMyVENDQW42Z0F3SUJBZ0lSQUxTdHNXanBQdUplajVTYURMNHRlTTB3Q2dZSUtvWkl6ajBFQXdJd09qRWcKTUI0R0ExVUVDaE1YYkc5dmNDQmhkWFJ2WjJWdVpYSmhkR1ZrSUdObGNuUXhGakFVQmdOVkJBTVREVVJ2YjIxdApZV011Ykc5allXd3dIaGNOTWpNd01qQXhNVGd5TWpJeFdoY05NalF3TXpJNE1UZ3lNakl4V2pBNk1TQXdIZ1lEClZRUUtFeGRzYjI5d0lHRjFkRzluWlc1bGNtRjBaV1FnWTJWeWRERVdNQlFHQTFVRUF4TU5SRzl2YlcxaFl5NXMKYjJOaGJEQlpNQk1HQnlxR1NNNDlBZ0VHQ0NxR1NNNDlBd0VIQTBJQUJQSzFZL21uRjN4V1FYbUtuTWcwTkNFRgpncDlhclQwOXR0c3N6cFZYTENRYXU1WWNyMmI4T0xKNEtCMytCMlhkTXNneW5Tc1pabU5tejNaV1k4MXNIWFNqCmdnRmpNSUlCWHpBT0JnTlZIUThCQWY4RUJBTUNBcVF3RXdZRFZSMGxCQXd3Q2dZSUt3WUJCUVVIQXdFd0R3WUQKVlIwVEFRSC9CQVV3QXdFQi96QWRCZ05WSFE0RUZnUVU4YVVrTzd5U3lRSE9VeXNlbWppRlFDNUZobVF3Z2dFRwpCZ05WSFJFRWdmNHdnZnVDRFVSdmIyMXRZV011Ykc5allXeUNDV3h2WTJGc2FHOXpkSUlFZFc1cGVJSUtkVzVwCmVIQmhZMnRsZElJSFluVm1ZMjl1Ym9jRWZ3QUFBWWNRQUFBQUFBQUFBQUFBQUFBQUFBQUFBWWNRL29BQUFBQUEKQUFBQUFBQUFBQUFBQVljUS9vQUFBQUFBQUFCVWdsci8vdm9TV29jUS9vQUFBQUFBQUFCVWdsci8vdm9TVzRjUQovb0FBQUFBQUFBQlVnbHIvL3ZvU1hJY1Evb0FBQUFBQUFBQVlUY0hDcWZiMWFJY0V3S2dCMVljUS9vQUFBQUFBCkFBQ1lNQUQvL3MrNVhZY0VDcEN0TzRjUS9vQUFBQUFBQUFDQnZ1NUFkcS9LMzRjUS9vQUFBQUFBQUFET2dRc2MKdlN3R25vY1Evb0FBQUFBQUFBQldOeGNhbXk5a1F6QUtCZ2dxaGtqT1BRUURBZ05KQURCR0FpRUFvcFB6MnpuSgpleXFoWlhHc2VRbEZrblg5Qis5a2xyU0pUOVdpYzR0bDloMENJUUNSd1diMzNITnNlNnI3R1dPckhlNGZ4YWlVCkE4MUU4WlM1NHNIM25OU1YzQT09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
}

func TestCreateNodeGuardClient(t *testing.T) {

	type args struct {
		nodeguardEndpoint string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "CreateNodeguardClient_Success",
			args: args{
				nodeguardEndpoint: "test:9001",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := CreateNodeGuardClient(tt.args.nodeguardEndpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateNodeguardClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				t.Errorf("NodeguardClient is nil")
			}

			if got1 == nil {
				t.Errorf("Connection is nil")
			}

		})
	}
}

func TestCreateLightningClient(t *testing.T) {

	tlsCertEncoded := getTlsCertEncoded()

	type args struct {
		nodeEndpoint   string
		tlsCertEncoded string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "CreateLightningClient_Success",
			args: args{
				nodeEndpoint:   "test:10001",
				tlsCertEncoded: tlsCertEncoded,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := CreateLightningClient(tt.args.nodeEndpoint, tt.args.tlsCertEncoded)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateLightningClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				t.Errorf("LightningClient is nil")
			}

			if got1 == nil {
				t.Errorf("Connection is nil")
			}

		})
	}
}

func TestCreateSwapClientClient(t *testing.T) {
	tlsCertEncoded := getTlsCertEncoded()

	type args struct {
		loopdEndpoint  string
		tlsCertEncoded string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "CreateSwapClientClient_Success",
			args: args{
				loopdEndpoint:  "test:11001",
				tlsCertEncoded: tlsCertEncoded,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := CreateSwapClientClient(tt.args.loopdEndpoint, tt.args.tlsCertEncoded)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSwapClientClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil {
				t.Errorf("SwapClientClient is nil")
			}

			if got1 == nil {
				t.Errorf("Connection is nil")
			}
		})
	}
}

func Test_generateCredentials(t *testing.T) {
	type args struct {
		tlsCertEncoded string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "GenerateCredentials_Success",
			args: args{
				tlsCertEncoded: getTlsCertEncoded(),
			},
			wantErr: false,
		},
		{
			name: "GenerateCredentials_InvalidTlsCert",
			args: args{
				tlsCertEncoded: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := generateCredentials(tt.args.tlsCertEncoded)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
