package rpc

import "context"

type TokenAuth struct {
	Token      string
	HeaderName string
}

func NewTokenAuth(token string, header string) *TokenAuth {
	return &TokenAuth{Token: token, HeaderName: header}
}

func (ta *TokenAuth) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		ta.HeaderName: ta.Token,
	}, nil
}

func (ta *TokenAuth) RequireTransportSecurity() bool {
	return false
}
