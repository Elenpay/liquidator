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
      --nodesHosts string        Command separated list of hostname:port to connect to
      --nodesMacaroons string    Command separated list of macaroons used in nodesHosts, in the same order of NODESHOSTS
      --nodesTLSCerts string     Command separated list of tls certs from LNDS in base64, in the same order of NODESHOSTS and NODESMACAROONS
      --pollingInterval string   Interval to poll data (default "15s")
```
# Setup
This project uses [just](https://github.com/casey/just) with the following recipes
```
Available recipes:
    build
    compile-proto
    cover-test
    init-submodules
    run
    test
```


# Build
```
just build
```

# Environment Variables / Flags

All the flags can be set as environment variables, with the following format, except stated, they are all mandatory:

- NODESHOSTS: Command separated list of hostname:port to connect to
- NODESMACAROONS : Command separated list of macaroons used in nodesHosts, in the same order of NODESHOSTS
- NODESTLSCERTS : Command separated list of tls certs from LNDS in **base64**, in the same order of NODESHOSTS and NODESMACAROONS
- POLLINGINTERVAL (optional) : Interval to poll data(default 15s)
- LOGLEVEL (optional) : Log level (default info) from: {trace, debug, info, warn, error, fatal, panic}
- LOGFORMAT (optional) : Log format (default json) from: {json, text}

# Metrics
The following metrics are exposed in the `/metrics` endpoint on port `9000` (e.g. `localhost:9000/metrics`)):
- `liquidator_channel_balance`: Channel balance ratio in 0 to 1 range, 0 means all the balance is on the local side of the channel, 1 means all the balance is on the remote side of the channel
 
Example:
 ```
liquidator_channel_balance{active="false",channel_id="118747255865345",local_node_alias="alice",local_node_pubkey="03b48034270e522e4033afdbe43383d66d426638927b940d09a8a7a0de4d96e807",remote_node_alias="",remote_node_pubkey="02f97d034c6c8f5ad95b1fe6abfe68ab154e85b1f5bb909815baeb5c8a46cdf622"} 0.99

liquidator_channel_balance{active="false",channel_id="125344325632000",local_node_alias="alice",local_node_pubkey="03b48034270e522e4033afdbe43383d66d426638927b940d09a8a7a0de4d96e807",remote_node_alias="",remote_node_pubkey="02f97d034c6c8f5ad95b1fe6abfe68ab154e85b1f5bb909815baeb5c8a46cdf622"} 0

liquidator_channel_balance{active="false",channel_id="131941395398656",local_node_alias="carol",local_node_pubkey="03485d8dcdd149c87553eeb80586eb2bece874d412e9f117304446ce189955d375",remote_node_alias="",remote_node_pubkey="02f97d034c6c8f5ad95b1fe6abfe68ab154e85b1f5bb909815baeb5c8a46cdf622"} 0

```

Other go-based metrics of the application are exposed for performance monitoring.