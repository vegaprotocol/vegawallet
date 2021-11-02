# Vega Wallet

`vegawallet` is the command line interface for running a Wallet service,
implemented in Go. It is used to sign transactions for use
on [Vega](#about-vega). Vega Wallet creates and manages HD wallets with ed25519
key pairs.

## Documentation

#### [Getting started with Vega Wallet][vega-documentation-website-create-wallet]
Learn how to install and run the stable version of Vega Wallet.

#### [Vega documentation][vega-documentation-website]
Learn more about how Vega works, and explore sample scripts for API trading

## Before continuing...

### I am not familiar with Vega Wallet...

If you want to know more about Vega Wallet, how it works and how to use it, refer to the page ["Using Vega Wallet"][vega-wallet-documentation-website].

### I don't know how to set up a Vega Wallet...

If you want to know more about how to create and use a Vega Wallet, refer to the page ["Create a Vega Wallet"][vega-documentation-website-create-wallet].

### I want to use the latest stable version...

If you want to use a stable version, refer to ["Create a Vega Wallet"][vega-documentation-website-create-wallet].

### Should I use the documentation in this repository?

If you are looking for the documentation for the stable version of Vega Wallet, refer to the [documentation website][vega-wallet-documentation-website]. **Do not refer to the documentation in this repository.**

The documentation living in this repository contains information about unreleased and unstable features, and it is meant for people running a version of Vega Wallet that is built from source code.

## A word about versions

**A release does not necessarily mean it is stable.** If a version is suffixed with `-pre` (ex: `v0.9.0-pre1`), this is not stable.

If you are not sure which version you are currently running, use the following command to find out:

```sh
vegawallet version
```

All releases can be seen on the [Releases][github-releases] page.

## Installation

To install Vega Wallet, you can download a released binary, or install it using the Golang toolchain.

### Download binaries

From the [Releases][github-releases] page, download the ZIP file matching your platform and open it.

|  Platform | Associated ZIP file            |
|-----------|--------------------------------|
| Windows   | `vegawallet-windows-amd64.zip` |
|  MacOS    | `vegawallet-darwin-amd64.zip`  |
| Linux     | `vegawallet-linux-amd64.zip`   |


### Installing from repository

You can install a released version using Golang toolchain:

```sh
go install code.vegaprotocol.io/vegawallet@VERSION
```

Replace `VERSION` with the release version of your choice.

For version `v0.9.0`, it would be:

```sh
go install code.vegaprotocol.io/vegawallet@v0.9.0
```

## Building from source

To build the Vega Wallet from the source code, use the following 

```sh
cd vegawallet && go build
```

#### Using the command-line

See a list of commands available in the wallet [here](cmd/README.md).

#### Using the API

Using the API is documented [here](service/README.md).

## Support

#### [Nolt][nolt]
Raise issues and see what others have raised.

#### [Discord][discord]
Ask us for help, find out about scheduled open sessions, and keep up with Vega
generally.

## About Vega

[Vega][vega-website] is a protocol for creating and trading derivatives on a
fully decentralised network. The network, secured with proof-of-stake, will
facilitate fully automated, end-to-end margin trading and execution of complex
financial products. Anyone will be able to build decentralised markets using the
protocol.

Read more at [https://vega.xyz][vega-website].

[vega-website]: https://vega.xyz
[vega-documentation-website]: https://docs.fairground.vega.xyz
[vega-documentation-website-create-wallet]: https://docs.fairground.vega.xyz/docs/vega-wallet/create-wallet/
[vega-wallet-documentation-website]: https://docs.fairground.vega.xyz/docs/vega-wallet/
[nolt]: https://vega-testnet.nolt.io/
[discord]: https://vega.xyz/discord
[github-release]: https://github.com/vegaprotocol/vegawallet/releases
