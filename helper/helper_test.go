package helper

import (
	"testing"

	"google.golang.org/grpc/metadata"
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

func TestAbsInt64(t *testing.T) {
	type args struct {
		x int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "positive",
			args: args{
				x: 1,
			},
			want: 1,
		},
		{
			name: "negative",
			args: args{
				x: -1,
			},
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AbsInt64(tt.args.x); got != tt.want {
				t.Errorf("AbsInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}
