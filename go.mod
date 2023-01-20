module liquidator

go 1.19

require (
	github.com/lightningnetwork/lnd/lnrpc v0.0.2
	google.golang.org/grpc v1.52.0
)

replace github.com/lightningnetwork/lnd/lnrpc => ./github.com/lightningnetwork/lnd/lnrpc

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20230117162540-28d6b9783ac4 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
