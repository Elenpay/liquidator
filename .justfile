set positional-arguments

fmt:
    go fmt github.com/Elenpay/...
init-submodules:
    git submodule init && git submodule update
build:
    go build
test:
    go test ./...
run: build
    go run .
install-loopd-loop:
    cd loop && go install . 
compile-lnrpc-proto:
    rm -rf ./github.com/lightningnetwork/lnd && protoc -I lnd/lnrpc --go_out=. --go-grpc_out=.  lnd/lnrpc/*.proto && cd ./github.com/lightningnetwork/lnd/lnrpc && go mod init lnrpc
compile-nodeguard-proto:
    rm -rf nodeguard && protoc -I rpc --go_out=. --go-grpc_out=.  rpc/*.proto && mockgen -destination ./nodeguard/nodeguard_mock.go -source nodeguard/nodeguard_grpc.pb.go  -package nodeguard && mockgen -destination ./nodeguard/nodeguard_mock.go -source nodeguard/nodeguard_grpc.pb.go  -package nodeguard
compile-loop-proto:
    rm -rf ./github.com/lightninglabs/loop && cd loop && protoc -I. -I looprpc -I swapserverrpc --go_out=.. --go-grpc_out=..  looprpc/*.proto  && cd ../github.com/lightninglabs/loop/looprpc && go mod init looprpc
cover-test:
    go test ./... -coverprofile=coverage.out; go tool cover -html=coverage.out
start-loopserver: build-loopserver
    docker-compose up -d
start-loopd-carol: 
    loopd   --network=regtest --debuglevel=debug --server.host=localhost:11009 --server.notls --lnd.host=localhost:10003 --lnd.macaroonpath=regtest.polar/volumes/lnd/carol/data/chain/bitcoin/regtest/admin.macaroon --lnd.tlspath=regtest.polar/volumes/lnd/carol/tls.cert --debuglevel=debug
start-all:  start-loopserver && start-loopd-carol
    @echo "started all (lopd-carol, loopserver)"
build-loopserver arg='':
    rm -rf regtest.polar && mkdir regtest.polar && tar -xf regtest.polar.zip -C regtest.polar && docker build {{arg}} -t loopserver -f Dockerfile.loopserver .
loop *args='':
    loop -n regtest $@
loopin sats='1000000':
    just loop in --amt {{sats}} -v
loopout chanid sats='500000' :
    just loop out --amt {{sats}} --channel {{chanid}} -v --fast

