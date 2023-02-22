package rpc

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"github.com/Elenpay/liquidator/nodeguard"
	"github.com/lightninglabs/loop/looprpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

// Generates the gRPC lightning client∏
func CreateLightningClient(nodeEndpoint string, tlsCertEncoded string) (lnrpc.LightningClient, *grpc.ClientConn, error) {
	creds, err := generateCredentials(tlsCertEncoded)
	if err != nil {
		return nil, nil, err
	}

	conn, err := getConn(nodeEndpoint, creds)
	if err != nil {
		return nil, nil, err
	}

	lightningClient := lnrpc.NewLightningClient(conn)

	return lightningClient, conn, nil
}

// Creates the SwapClient client similar to CreateLightningClient function
func CreateSwapClientClient(loopdEndpoint string, tlsCertEncoded string) (looprpc.SwapClientClient, *grpc.ClientConn, error) {

	creds, err := generateCredentials(tlsCertEncoded)
	if err != nil {
		return nil, nil, err
	}
	conn, err := getConn(loopdEndpoint, creds)
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
func generateCredentials(tlsCertEncoded string) (credentials.TransportCredentials, error) {
	//Generate TLS credentials from directory
	tlsCertDecoded, err := base64.StdEncoding.DecodeString(tlsCertEncoded)
	if err != nil {

		err := fmt.Errorf("Failed to decode TLS cert: %v", err)
		log.Error(err)
		return nil, err
	}

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(tlsCertDecoded) {
		err := fmt.Errorf("Failed to append certificates")
		log.Error(err)
	}

	creds := credentials.NewClientTLSFromCert(cp, "")

	if err != nil {
		log.Errorf("Failed to create ClientTLS from credentials %v", err)
		return nil, err
	}

	return creds, nil

}
