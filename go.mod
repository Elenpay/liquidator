module liquidator

go 1.19

require (
	google.golang.org/grpc v1.52.0
	github.com/lightningnetwork/lnd/lnrpc v0.0.2
)
replace (
	github.com/lightningnetwork/lnd/lnrpc => ./github.com/lightningnetwork/lnd/lnrpc
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20230117162540-28d6b9783ac4 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
