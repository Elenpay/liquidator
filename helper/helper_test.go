package helper

import (
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestGenerateContextWithMacaroon_Positive(t *testing.T) {

	//Arrange
	macaroon := "test"
	//Act
	ctx, err := GenerateContextWithMacaroon(macaroon)

	//Assert

	//Assert that there is no error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	//Assert that the context is not nil
	if ctx == nil {
		t.Errorf("Expected context to not be nil")
	}

	//Assert that the context has a macaroon
	//Cast context to metadata
	md, ok := metadata.FromOutgoingContext(ctx)

	if !ok {
		t.Errorf("Expected context to have metadata")
	}

	//Assert that the metadata has a macaroon
	x := md["macaroon"][0]
	if x != macaroon {
		t.Errorf("Expected metadata to have macaroon %v, got %v", macaroon, x)
	}

}
func TestGenerateContextWithMacaroon_Negative(t *testing.T) {

	//Arrange
	macaroon := ""
	//Act
	_, err := GenerateContextWithMacaroon(macaroon)

	//Assert

	//Assert that there is no error
	if err == nil {
		t.Errorf("Expected no error, got %v", err)
	}

}
