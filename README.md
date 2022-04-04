# Vega Wallet

`vegawallet` is the command line interface for running a Wallet service, implemented in Go. It is used to sign transactions for use on [Vega](#about-vega). Vega Wallet creates and manages HD wallets with ed25519 key pairs.

## Documentation

#### [Getting started with the CLI app for Vega Wallet][vega-documentation-website-create-wallet]
Learn how to install and run the stable version of the Vega Wallet CLI app.

#### [Vega documentation][vega-documentation-website]
Learn more about how Vega works, and explore sample scripts for API trading

## Before continuing...

### I am not familiar with Vega Wallet...

If you want to know more about Vega Wallet, how it works and how to use it, refer to the page ["Using Vega Wallet"][vega-wallet-documentation-website].

### I don't know how to set up a Vega Wallet...

If you want to know more about how to create and use a Vega Wallet with the CLI app, refer to the page ["Create a Vega Wallet"][vega-documentation-website-create-wallet].

If you'd prefer to use a visual interface to interact with a Vega Wallet, you can use the Vega Wallet desktop app. Refer to the [desktop app documentation][vega-wallet-desktop-getting-started].  

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

|  Platform       | Associated ZIP file            |
|-----------------|--------------------------------|
| Windows         | `vegawallet-windows-amd64.zip` |
| Windows (ARM64) | `vegawallet-windows-arm64.zip` |
| MacOS           | `vegawallet-darwin-amd64.zip`  |
| MacOS (ARM64)   | `vegawallet-darwin-arm64.zip`  |
| Linux           | `vegawallet-linux-amd64.zip`   |
| Linux (ARM64)   | `vegawallet-linux-arm64.zip`   |


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

See a list of commands available using:

```sh
vegawallet --help
```

#### Using the API

Using the API is documented [here](service/README.md).

## Support

#### [Feedback][feedback]
Raise issues and see what others have raised.

#### [Discord][discord]
Ask us for help, find out about scheduled open sessions, and keep up with Vega
generally.

## About Vega

[Vega][vega-website] is a protocol for creating and trading derivatives on a fully decentralised network. The network, secured with proof-of-stake, will facilitate fully automated, end-to-end margin trading and execution of complex financial products. Anyone will be able to build decentralised markets using the protocol.

Read more at [https://vega.xyz][vega-website].

[vega-website]: https://vega.xyz
[vega-documentation-website]: https://docs.vega.xyz
[vega-documentation-website-create-wallet]: https://docs.vega.xyz/docs/tools/vega-wallet/CLI-wallet/create-wallet
[vega-wallet-documentation-website]: https://docs.vega.xyz/docs/tools/vega-wallet/
[vega-wallet-desktop-getting-started]: https://docs.vega.xyz/docs/tools/vega-wallet/desktop-app/latest/getting-started
[feedback]: https://github.com/vegaprotocol/feedback/discussions/
[discord]: https://vega.xyz/discord
[github-releases]: https://github.com/vegaprotocol/vegawallet/releases
