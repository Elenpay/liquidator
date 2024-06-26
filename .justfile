set positional-arguments
set fallback := true

fmt:
    go fmt github.com/Elenpay/...
init-submodules:
    git submodule init && git submodule update
build:
    go build
test:
    go test ./...
run *args='': build
    go run . $@
install-loopd-loop:
    cd loop/cmd/loop && go install .
    cd loop/cmd/loopd && go install .

compile-nodeguard-proto:
    rm -rf nodeguard && protoc -I rpc --go_out=. --go-grpc_out=.  rpc/*.proto && mockgen -destination ./nodeguard/nodeguard_mock.go -source nodeguard/nodeguard_grpc.pb.go  -package nodeguard && mockgen -destination ./nodeguard/nodeguard_mock.go -source nodeguard/nodeguard_grpc.pb.go  -package nodeguard
compile-lnd-client-mock:
    mockgen -destination lightning_rpc_mock.go -source lnd/lnrpc/lightning_grpc.pb.go  -package main
compile-provider-mock:
    mockgen -destination ./provider/provider_mock.go -source provider/provider.go  -package provider && mockgen -destination ./provider/loopd_mock.go -source loop/looprpc/client_grpc.pb.go  -package provider 
cover-test:
    go test ./... -coverprofile=coverage.out; go tool cover -html=coverage.out
start-loopserver: build-loopserver
    docker-compose up -d
start-loopd-carol: 
    loopd   --network=regtest --debuglevel=debug --server.host=localhost:11009 --server.notls --lnd.host=localhost:10003 --lnd.macaroonpath=regtest.polar/volumes/lnd/carol/data/chain/bitcoin/regtest/admin.macaroon --lnd.tlspath=regtest.polar/volumes/lnd/carol/tls.cert --debuglevel=debug --loopdir .loop
start-all:  start-loopserver && start-loopd-carol
    @echo "started all (lopd-carol, loopserver)"
build-loopserver arg='':
    rm -rf regtest.polar && mkdir regtest.polar && tar -xf regtest.polar.zip -C regtest.polar && docker build {{arg}} -t loopserver -f Dockerfile.loopserver .
loop *args='':
    loop -n regtest --loopdir .loop $@
loopin sats='1000000':
    just loop in --amt {{sats}} -v
loopout chanid sats='500000' :
    just loop out --amt {{sats}} --channel {{chanid}} -v --fast
unzip-loopd-datadir:
    rm -rf .loop && mkdir .loop; unzip -d .loop/regtest loopd.zip; rm -rf .loop/regtest/loop.db
mine:
    while true; do docker exec polar-n1-backend1 bitcoin-cli -regtest -rpcuser=polaruser -rpcpassword=polarpass -generate 1; sleep 60; done


