# Description
A CLI tool to monitor and automate the liquidity of your LND channels

```
Usage:
  liquidator [flags]

Flags:
  -h, --help                     help for liquidator
      --nodesHosts strings       Command separated list of hostname:port to connect to
      --nodesMacaroons strings   Command separated list of macaroons used in nodesHosts, in the same order
      --nodesTLSCerts strings    Command separated list of tls certs from LNDS in base64, in the same order of NODESHOSTS and NODESMACAROONS
      --pollingInterval int      Interval to poll data in seconds (default 10)
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

