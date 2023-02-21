<p align="center">
  <img src="liquidator.png" width="100px" />
</p>

# Description
A CLI tool to monitor and automate the liquidity of your LND channels

```
Usage:
  liquidator [flags]

Flags:
  -h, --help                     help for liquidator
      --logFormat string         Log format from: {text, json} (default "text")
      --logLevel string          Log level from values: {trace, debug, info, warn, error, fatal, panic} (default "info")
      --loopdHosts string        Command separated list of hostname:port to connect to loopd, each position corresponds to the managed node in nodesHosts
      --loopdMacaroons string    Command separated list of macaroons used in loopdHosts in hex, in the same order of loopdHosts
      --loopdTLSCerts string     Command separated list of tls certs from loopd in base64, in the same order of loopdHosts and loopdMacaroons
      --nodeguardHost string     Hostname:port to connect to nodeguard
      --nodesHosts string        Command separated list of hostname:port to connect to
      --nodesMacaroons string    Command separated list of macaroons used in nodesHosts in hex, in the same order of NODESHOSTS
      --nodesTLSCerts string     Command separated list of tls certs from LNDS in base64, in the same order of NODESHOSTS and NODESMACAROONS
      --pollingInterval string   Interval to poll data (default "15s")
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


# Environment Variables / Flags

All the flags can be set as environment variables, with the following format, except stated, they are all mandatory:

- NODESHOSTS: Command separated list of hostname:port to connect to
- NODESMACAROONS : Command separated list of macaroons in **hex** used in nodesHosts, in the same order of NODESHOSTS
- NODESTLSCERTS : Command separated list of tls certs from LNDS in **base64**, in the same order of NODESHOSTS and NODESMACAROONS
- POLLINGINTERVAL (optional) : Interval to poll data(default 15s)
- LOGLEVEL (optional) : Log level (default info) from: {trace, debug, info, warn, error, fatal, panic}
- LOGFORMAT (optional) : Log format (default json) from: {json, text}

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