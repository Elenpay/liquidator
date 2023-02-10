package rpc

import (
	"crypto/x509"
	"encoding/base64"

	"github.com/Elenpay/liquidator/nodeguard"
	"github.com/lightninglabs/loop/looprpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Generates the gRPC lightning client∏
func CreateLightningClient(nodeEndpoint string, tlsCertEncoded string) (lnrpc.LightningClient, *grpc.ClientConn, error) {
	creds := generateCredentials(tlsCertEncoded)

	conn, err := getConn(nodeEndpoint, creds)
	if err != nil {
		return nil, nil, err
	}

	lightningClient := lnrpc.NewLightningClient(conn)

	return lightningClient, conn, nil
}

// Creates the SwapClient client similar to CreateLightningClient function
func CreateSwapClientClient(nodeEndpoint string, tlsCertEncoded string) (looprpc.SwapClientClient, *grpc.ClientConn, error) {

	creds := generateCredentials(tlsCertEncoded)

	conn, err := getConn(nodeEndpoint, creds)
	if err != nil {
		return nil, nil, err
	}

	swapClient := looprpc.NewSwapClientClient(conn)

	return swapClient, conn, nil
}

// Creates the NodeGuard grpc client
func CreateNodeGuardClient(nodeGuardEndpoint string) (nodeguard.NodeGuardServiceClient, *grpc.ClientConn, error) {

	//TODO ADD TLS to NodeGuard API

	conn, err := getConn(nodeGuardEndpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	client := nodeguard.NewNodeGuardServiceClient(conn)

	return client, conn, nil
}

// generates the gRPC connection based on the node endpoint and the credentials
func getConn(nodeEndpoint string, creds credentials.TransportCredentials) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(nodeEndpoint, grpc.WithTransportCredentials(creds))

	if err != nil {
		log.Errorf("did not connect: %v", err)
		return nil, err
	}

	return conn, nil
}

// Generates gRPC credentials for the clients
func generateCredentials(tlsCertEncoded string) credentials.TransportCredentials {
	//Generate TLS credentials from directory
	tlsCertDecoded, err := base64.StdEncoding.DecodeString(tlsCertEncoded)
	if err != nil {
		log.Fatalf("Failed to decode TLS cert: %v", err)
	}

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(tlsCertDecoded) {
		log.Fatalf("credentials: failed to append certificates")
	}

	creds := credentials.NewClientTLSFromCert(cp, "")

	if err != nil {
		log.Fatalf("Failed to load credentials: %v", err)
	}

	return creds

}
