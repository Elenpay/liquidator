set positional-arguments

init-submodules:
    git submodule init && git submodule update
build:
    go build
test:
    go test
run: build
    go run .
install loopd:
    go install .
compile-proto:
    rm -rf ./github.com && protoc -I lnd/lnrpc --go_out=. --go-grpc_out=.  lnd/lnrpc/*.proto && cd ./github.com/lightningnetwork/lnd/lnrpc && go mod init lnrpc
cover-test:
    go test -coverprofile=coverage.out && go tool cover -html=coverage.out

build-loopserver arg='':
    rm -rf regtest.polar && mkdir regtest.polar && tar -xf regtest.polar.zip -C regtest.polar && docker build {{arg}} -t loopserver -f Dockerfile.loopserver .
start-loopd-carol:
    loopd   --network=regtest --debuglevel=debug --server.host=localhost:11009 --server.notls --lnd.host=localhost:10003 --lnd.macaroonpath=regtest.polar/volumes/lnd/carol/data/chain/bitcoin/regtest/admin.macaroon --lnd.tlspath=regtest.polar/volumes/lnd/carol/tls.cert
loop *args='':
    loop -n regtest $@
loopin sats='1000000':
    just loop in --amt {{sats}} -v
loopout chanid sats='500000' :
    just loop out --amt {{sats}} --channel {{chanid}} -v --fast

