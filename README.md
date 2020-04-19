# Vega Wallet

`vegawallet` is the command line interface for running a Wallet service, implemented in Go. It is uses to sign transactions for use on [Vega](#about-vega)

## Running
### Configuration
Vega Wallet creates and manages ed25519 keypairs for one or more wallets. To get started, create yourself a wallet and a keypair by running:

```console
vegawallet genkey -n walletname -p password
```

Will generate a new pair of keys:
```console
new generated keys:
public: 473…5e
private: 1…e
```
And store them in `~/.vega/wallets/walletname` by default.

You can now use this wallet and public key to sign arbitrary messages:
 ```console
vegawallet sign -n walletname -p password --pubkey 473…5e --message $(echo 'test' | base64)
```
Will return the signed message:
```console
YN6GvBONF…==
```

### Web accessible API
Vega Wallet can also host a REST-ish interface to make itself available for local and remote clients. This can be used to make keys available to web-based Vega clients, for example Vega Console.

```console
vegawallet service init --genrsakey
vegawallet service run
	{"level":"info","ts":1587317545.61634,"logger":"wallet","caller":"wallet/service.go:147","msg":"starting wallet http server","address":"0.0.0.0:1789"}
```

Using the API is documented [here](./wallet/README.md)

## Building
```console
cd go-wallet
make
```

# About Vega
 [Vega](https://vega.xyz) is a protocol for creating and trading derivatives on a fully decentralised network. The network, secured with proof-of-stake, will facilitate fully automated, end-to-end margin trading and execution of complex financial products. Anyone will be able to build decentralised markets using the protocol.

Read more at [https://vega.xyz](https://vega.xyz)
