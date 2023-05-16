package rpc

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"github.com/Elenpay/liquidator/lndconnect"
	"github.com/Elenpay/liquidator/nodeguard"
	"github.com/lightninglabs/loop/looprpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// Generates the gRPC lightning client∏
func CreateLightningClient(lndConnectParams lndconnect.LndConnectParams) (lnrpc.LightningClient, *grpc.ClientConn, error) {
	creds, err := generateCredentials(lndConnectParams.Cert)
	if err != nil {
		return nil, nil, err
	}

	endpoint := fmt.Sprintf("%s:%s", lndConnectParams.Host, lndConnectParams.Port)

	conn, err := getConn(endpoint, creds)
	if err != nil {
		return nil, nil, err
	}

	lightningClient := lnrpc.NewLightningClient(conn)

	return lightningClient, conn, nil
}

// Creates the SwapClient client similar to CreateLightningClient function
func CreateSwapClientClient(lndConnectParams lndconnect.LndConnectParams) (looprpc.SwapClientClient, *grpc.ClientConn, error) {

	creds, err := generateCredentials(lndConnectParams.Cert)
	if err != nil {
		return nil, nil, err
	}

	endpoint := fmt.Sprintf("%s:%s", lndConnectParams.Host, lndConnectParams.Port)

	conn, err := getConn(endpoint, creds)
	if err != nil {
		return nil, nil, err
	}

	swapClient := looprpc.NewSwapClientClient(conn)

	return swapClient, conn, nil
}

// Creates the NodeGuard grpc client
func CreateNodeGuardClient(nodeGuardEndpoint string) (nodeguard.NodeGuardServiceClient, *grpc.ClientConn, error) {

	//TODO ADD TLS to NodeGuard API

	conn, err := getConn(nodeGuardEndpoint, insecure.NewCredentials())
	if err != nil {
		return nil, nil, err
	}

	client := nodeguard.NewNodeGuardServiceClient(conn)

	return client, conn, nil
}

// generates the gRPC connection based on the node endpoint and the credentials
func getConn(nodeEndpoint string, creds credentials.TransportCredentials) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(nodeEndpoint, grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))

	if err != nil {
		log.Errorf("did not connect: %v", err)
		return nil, err
	}

	return conn, nil
}

// Generates gRPC credentials for the clients
func generateCredentials(certDer string) (credentials.TransportCredentials, error) {

	base64decoded, err := base64.RawURLEncoding.DecodeString(certDer)
	if err != nil {
		log.Errorf("Failed to decode base64 string")
		return nil, fmt.Errorf("failed to decode base64 string")
	}
	cp := x509.NewCertPool()
	cert, err := x509.ParseCertificate(base64decoded)
	if err != nil {
		log.Errorf("Failed to parse certificate")
		return nil, fmt.Errorf("failed to parse certificate")
	}
	cp.AddCert(cert)

	creds := credentials.NewClientTLSFromCert(cp, "")

	return creds, nil

}
