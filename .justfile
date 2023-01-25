init-submodules:
    git submodule init && git submodule update
build:
    go build
test:
    go test
run:
    go run .
compile-proto:
    rm -rf ./github.com && protoc -I lnd/lnrpc --go_out=. --go-grpc_out=.  lnd/lnrpc/*.proto && cd ./github.com/lightningnetwork/lnd/lnrpc && go mod init lnrpc
cover-test:
    go test -coverprofile=coverage.out && go tool cover -html=coverage.out



