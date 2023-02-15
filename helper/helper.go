package helper

import (
	"context"
	"errors"

	"google.golang.org/grpc/metadata"
)

// Generates the context with metadata for the grpc connection with a macaroon for a given node
func GenerateContextWithMacaroon(macaroon string) (context.Context, error) {

	if macaroon == "" {
		err := errors.New("macaroon is empty")
		return nil, err
	}

	ctx := context.Background()

	md := metadata.New(map[string]string{"macaroon": macaroon})

	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx, nil

}

// Write a function that returns the absolute value of a int64
func AbsInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
