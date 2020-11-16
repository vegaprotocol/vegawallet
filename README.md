# Vega Wallet

`vegawallet` is the command line interface for running a Wallet service, implemented in Go. It is uses to sign transactions for use on [Vega](#about-vega).

## Running
### Configuration
Vega Wallet creates and manages ed25519 keypairs for one or more wallets. To get started, create yourself a wallet and a keypair by running:

```console
wallet@vega:~$ vegawallet genkey -n walletname -p password
new generated keys:
public: 473…5e
private: 1…e
```

And store them in `~/.vega/wallets/walletname` by default.

You can now use this wallet and public key to sign arbitrary messages:
 ```console
wallet@vega:~$ vegawallet sign -n walletname -p password --pubkey 473…5e --message $(echo 'test' | base64)
YN6GvBONF…==
```

### Web accessible API
Vega Wallet can also host a REST interface to make itself available for local and remote clients. This can be used to make keys available to web-based Vega clients, for example Vega Console.

#### Service configuration
Before starting the wallet service, some configuration is required, this can be done by running the following command:

```console
wallet@vega:~$ vegawallet service init
{"level":"info","ts":1605554344.734188,"caller":"wallet/config.go:125","msg":"wallet service configuration generated successfully","path":"path/wallet-service-config.toml"}
{"level":"info","ts":1605554347.727988,"caller":"wallet/config.go:173","msg":"wallet rsa key generated successfully","path":"path/wallet_rsa"}
```

This would have generated a new default configuration into `path/wallet-service-config.toml`.


#### Running the service
The service can be run in 3 different ways

Just as an API, to do so run the following command:
```console
wallet@vega:~$ vegawallet service run
{"level":"info","ts":1587317545.61634,"logger":"wallet","caller":"wallet/service.go:147","msg":"starting wallet http server","address":"0.0.0.0:1789"}
```

As an API, but also proxying the Vega Console locally:
```console
wallet@vega:~$ vegawallet service run --console-proxy
{"level":"info","ts":1605554589.694528,"caller":"cmd/console.go:41","msg":"starting console proxy","proxy.address":"localhost:8080","address":"dev.vega.trading"}
{"level":"info","ts":1587317545.61634,"logger":"wallet","caller":"wallet/service.go:147","msg":"starting wallet http server","address":"0.0.0.0:1789"}
```
This will proxy the Vega Console to localhost:8080 (you can edit this in the service configuration file).

As an API, but also starting the Vega Console in a native browser window:
```console
wallet@vega:~$ vegawallet service run --console-ui
{"level":"info","ts":1605554589.694528,"caller":"cmd/console.go:41","msg":"starting console proxy","proxy.address":"localhost:8080","address":"dev.vega.trading"}
{"level":"info","ts":1587317545.61634,"logger":"wallet","caller":"wallet/service.go:147","msg":"starting wallet http server","address":"0.0.0.0:1789"}
```


Using the API is documented [here](./wallet/README.md).

## Building
```console
wallet@vega:~$ cd go-wallet && make
```

# About Vega
 [Vega](https://vega.xyz) is a protocol for creating and trading derivatives on a fully decentralised network. The network, secured with proof-of-stake, will facilitate fully automated, end-to-end margin trading and execution of complex financial products. Anyone will be able to build decentralised markets using the protocol.

Read more at [https://vega.xyz](https://vega.xyz).
