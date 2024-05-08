# NodeGuard Liquidator
[![Go Report Card](https://goreportcard.com/badge/github.com/Elenpay/liquidator)](https://goreportcard.com/report/github.com/Elenpay/liquidator)
<p align="center">
  <img src="liquidator.png" width="100px" />
</p>

# Description
A CLI tool to monitor and automate the liquidity of your LND channels
```
Usage:
  liquidator [flags]

Flags:
      --backoffCoefficient float       Coefficient to apply to the backoff (default 0.95)
      --backoffLimit float             Limit coefficient of the backoff (default 0.1)
  -h, --help                           help for liquidator
      --limitFeesL2 float              Limit fee ratio for swaps max routing fee e.g. 0.01 = 1% (default 0.002)
      --limitQuoteFees float           Limit fee ratio for swaps quotes (e.g. onchain+service fee estimation) e.g. 0.01 = 1% (default 0.005)
      --lndconnecturis string          CSV of lndconnect strings to connect to lnd(s)
      --logFormat string               Log format from: {text, json} (default "text")
      --logLevel string                Log level from values: {trace, debug, info, warn, error, fatal, panic} (default "info")
      --loopdconnecturis string        CSV of loopdconnect strings to connect to loopd(s)
      --nodeguardHost string           Hostname:port to connect to nodeguard
      --pollingInterval string         Interval to poll data (default "15s")
      --retriesBeforeBackoff int       Number of retries before applying backoff to the swap (default 3)
      --swapPublicationOffset string   Swap publication deadline offset (Maximum time for the swap provider to publish the swap) (default "60m")
      --sweepConfTarget string         Target number of confirmations for swaps, this uses bitcoin core broken estimator, procced with caution (default "400")
```
# Requirements
This project uses [just](https://github.com/casey/just) with the following recipes
```
Available recipes:
    build
    build-loopserver arg=''
    compile-lnrpc-proto
    compile-loop-proto
    compile-nodeguard-proto
    compile-provider-mocks
    cover-test
    fmt
    init-submodules
    install-loopd-loop
    loop *args=''
    loopin sats='1000000'
    loopout chanid sats='500000'
    mine
    run *args=''
    start-all
    start-loopd-carol
    start-loopserver
    test
    unzip-loopd-datadir
```
# Environment Variables / Flags

All the flags can be set as environment variables, with the following format, except stated, they are all mandatory:

- LNDCONNECTURIS : CSV of lndconnect strings to connect to lnd(s)\
- LIMITQUOTEFEES (optional) : Limit to total Max Quote Fees (L1+service fee) (default 0.007 -> 0.7% of the Swap size)
- LIMITFEESL2 (optional) : Limit to total Max Routing Fees (L2) (default 0.002 -> 0.2% of the Swap size)
- LOOPDCONNECTURIS : CSV of loopdconnect strings to connect to loopd(s)
- POLLINGINTERVAL (optional) : Interval to poll data(default 15s)
- LOGLEVEL (optional) : Log level (default info) from: {trace, debug, info, warn, error, fatal, panic}
- LOGFORMAT (optional) : Log format (default json) from: {json, text}
- SWAPPUBLICATIONOFFSET (optional) : Swap publication deadline offset (Maximum time for the swap provider to publish the swap) (default 30m)
- RETRIESBEFOREBACKOFF (optional) : Number of retries before applying backoff to the swap (default: 3)
- BACKOFFCOEFFICIENT (optional) : Coefficient to apply to the backoff (default: 0.95)
- BACKOFFLIMIT (optional) : Limit coefficient of the backoff (default: 0.1)

# Build & test
## Build
```
just build
```

## Testing
```
just test
```
# Dev environment
1. Launch a local regtest with polar from regtest.polar.zip
2. Lauch with VS Code pre-defined configuration

## Setup loop/loopd/loopserver regtest environment
1. Make sure your regtest.polar.zip network in polar is running
2. Unzip loopd datadir
````
just unzip-loopd-datadir
````
3. Get git submodules
````
just init-submodules
````
4. Compile and install loopd/loop binaries (Make sure your golang install bin dir is reachable from your PATH)
````
just install-loopd-loop
````
5. Using just, run the following command:
````
just start-all
````
6. This comand should have build a `loopserver `docker image, and started a `loopserver` container along with loopd as a native binary.

## Loop just recipes
There are a few recipes using `just -l` to interact with loopd for loop in, loop out and calling loop CLI with args (`just loop <args>`).




# Metrics
The following metrics are exposed in the `/metrics` endpoint on port `9000` (e.g. `localhost:9000/metrics`)):
- `liquidator_channel_balance`: Channel balance ratio in 0 to 1 range, 0 means all the balance is on the local side of the channel, 1 means all the balance is on the remote side of the channel
 
Example:
 ```
liquidator_channel_balance{active="false",chan_id="118747255865345",local_node_alias="alice",local_node_pubkey="03b48034270e522e4033afdbe43383d66d426638927b940d09a8a7a0de4d96e807",remote_node_alias="",remote_node_pubkey="02f97d034c6c8f5ad95b1fe6abfe68ab154e85b1f5bb909815baeb5c8a46cdf622",initiator="false"} 0.99

liquidator_channel_balance{active="false",chan_id="125344325632000",local_node_alias="alice",local_node_pubkey="03b48034270e522e4033afdbe43383d66d426638927b940d09a8a7a0de4d96e807",remote_node_alias="",remote_node_pubkey="02f97d034c6c8f5ad95b1fe6abfe68ab154e85b1f5bb909815baeb5c8a46cdf622",initiator="false"} 0

liquidator_channel_balance{active="false",chan_id="131941395398656",local_node_alias="carol",local_node_pubkey="03485d8dcdd149c87553eeb80586eb2bece874d412e9f117304446ce189955d375",remote_node_alias="",remote_node_pubkey="02f97d034c6c8f5ad95b1fe6abfe68ab154e85b1f5bb909815baeb5c8a46cdf622",initiator="false"} 0

```

Other go-based metrics of the application are exposed for performance monitoring.
