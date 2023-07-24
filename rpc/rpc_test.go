package rpc

import (
	"testing"

	"github.com/Elenpay/liquidator/lndconnect"
)

func getTlsCertEncoded() string {
	return "MIICJTCCAcygAwIBAgIQLYfp6m1vP9wFBXOcE-UsaDAKBggqhkjOPQQDAjAxMR8wHQYDVQQKExZsbmQgYXV0b2dlbmVyYXRlZCBjZXJ0MQ4wDAYDVQQDEwVjYXJvbDAeFw0yMzAzMjkxNTM4MjBaFw0yNDA1MjMxNTM4MjBaMDExHzAdBgNVBAoTFmxuZCBhdXRvZ2VuZXJhdGVkIGNlcnQxDjAMBgNVBAMTBWNhcm9sMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEcXT4dekJnAiZWd8Pk3FgL1BSFXMRwLGSAlk7Di5hIJnIA1B_o8RWKzlPz7u3Aw5mmWHhN8B2MWMylWlWB2130KOBxTCBwjAOBgNVHQ8BAf8EBAMCAqQwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDwYDVR0TAQH_BAUwAwEB_zAdBgNVHQ4EFgQUDOS-19_0LFGf62WRyaaUSLc3j98wawYDVR0RBGQwYoIFY2Fyb2yCCWxvY2FsaG9zdIIFY2Fyb2yCDnBvbGFyLW4xLWNhcm9sggR1bml4ggp1bml4cGFja2V0ggdidWZjb25uhwR_AAABhxAAAAAAAAAAAAAAAAAAAAABhwSsFQAFMAoGCCqGSM49BAMCA0cAMEQCHxYe59PCXrTtSmGsOjfQo6V-sS8j73cqWOzTQbvgI3gCIQCj7sOxnZWBwilec7t8bBXjwPgX9frv8408JW4QhNFOUg"
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
		lndconnectParams lndconnect.LndConnectParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "CreateLightningClient_Success",
			args: args{
				lndconnectParams: lndconnect.LndConnectParams{
					Host:     "localhost",
					Port:     "10009",
					Cert:     tlsCertEncoded,
					Macaroon: "0201036c6e6402f801030a101ec5b6370c166f6c8e2853164109145a1201301a160a0761646472657373120472656164120577726974651a130a04696e666f120472656164120577726974651a170a08696e766f69636573120472656164120577726974651a210a086d616361726f6f6e120867656e6572617465120472656164120577726974651a160a076d657373616765120472656164120577726974651a170a086f6666636861696e120472656164120577726974651a160a076f6e636861696e120472656164120577726974651a140a057065657273120472656164120577726974651a180a067369676e6572120867656e6572617465120472656164000006208e957e78ec39e7810fad25cfc43850b8e9e7c079843b8ec7bb5522bba12230d6",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := CreateLightningClient(tt.args.lndconnectParams)
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
		lndconnectParams lndconnect.LndConnectParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "CreateSwapClientClient_Success",
			args: args{
				lndconnectParams: lndconnect.LndConnectParams{
					Host:     "localhost",
					Port:     "11010",
					Cert:     tlsCertEncoded,
					Macaroon: "0201036c6e6402f801030a101ec5b6370c166f6c8e2853164109145a1201301a160a0761646472657373120472656164120577726974651a130a04696e666f120472656164120577726974651a170a08696e766f69636573120472656164120577726974651a210a086d616361726f6f6e120867656e6572617465120472656164120577726974651a160a076d657373616765120472656164120577726974651a170a086f6666636861696e120472656164120577726974651a160a076f6e636861696e120472656164120577726974651a140a057065657273120472656164120577726974651a180a067369676e6572120867656e6572617465120472656164000006208e957e78ec39e7810fad25cfc43850b8e9e7c079843b8ec7bb5522bba12230d6",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := CreateSwapClientClient(tt.args.lndconnectParams)
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
