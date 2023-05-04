package lndconnect

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
)

type LndConnectParams struct {
	Host     string // Hostname or IP address
	Port     string // Port number
	Cert     string // Base64 of DER-encoded TLS certificate
	Macaroon string // Hex-encoded macaroon
}

// Parse an lndconnect URI and returns the parameters
func Parse(uri string) (LndConnectParams, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return LndConnectParams{}, err
	}

	if parsed.Scheme != "lndconnect" {
		return LndConnectParams{}, fmt.Errorf("invalid scheme: %s", parsed.Scheme)
	}

	host := parsed.Hostname()
	port := parsed.Port()
	if port == "" {
		port = "10009"
	}

	cert := parsed.Query().Get("cert")

	macaroon := parsed.Query().Get("macaroon")
	decodedMacaroon, err := base64.RawURLEncoding.DecodeString(macaroon)
	if err != nil {
		return LndConnectParams{}, fmt.Errorf("error decoding macaroon: %w", err)
	}
	macaroon = hex.EncodeToString(decodedMacaroon)

	return LndConnectParams{
		Host:     host,
		Port:     port,
		Cert:     cert,
		Macaroon: macaroon,
	}, nil
}
