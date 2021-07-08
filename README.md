# Vega Wallet

`vegawallet` is the command line interface for running a Wallet service,
implemented in Go. It is used to sign transactions for use
on [Vega](#about-vega). Vega Wallet creates and manages HD wallets with ed25519
key pairs.

## Installation

These instructions are written to be used in command line. Below, in the
snippets, you'll see commands in blue text. Copy those instructions and paste
them into your command line interface.

### Building

```console
cd go-wallet && make
```

**Note:** This will install the wallet under the name `go-wallet`, and
not `vegawallet`. Thus, when reading the documentation, replace `vegawallet`
by `go-wallet`.

### Download

Download and save the zip file
from [Releases](https://github.com/vegaprotocol/go-wallet/releases). We suggest
you keep track of where you've saved the file, because that's where the command
line interface will look for it.

#### MacOS

Download `vegawallet-darwin-amd64.zip`

When you open the file, you may need to change your system preferences for this
specific instance, in order to run Vega Wallet. If you open the file from
downloads, you may get a message saying:

> "vegawallet-darwin-amd64” cannot be opened because it is from an unidentified
> developer.

Click on the `(?)` help button, which will open a window that links you to the
`System Preferences`, and instructs you how to allow this software to run.

You’ll need to go to `System Preferences > Security & Privacy > General`, and
choose `Open Anyway`.

#### Windows

Download `vegawallet-windows-amd64.zip`

You may need to change your system preferences for this specific instance, in
order to run Vega Wallet. If you open the file from downloads, you may get a
message from Windows Defender saying:

> "vegawallet-windows-amd64" cannot be opened because it is from an unidentified
> developer.

Click on the `(More info)` text, which will reveal a button to `Run anyway`.

#### Linux

Download `vegawallet-linux-amd64.zip`

## Usage

### Using the CLI

Using the CLI is documented [here](cmd/README.md)

### Using the API

**Important:** Before using the API, you will have to initialise the service
using the CLI.

Using the API is documented [here](service/README.md).

## Deposit tokens

Once the wallet and its service have been initialised, you'll need to **deposit
Ropsten Ethereum-based tokens** to start trading.

You can create and deposit assets directly through the proxy Console via Wallet.

If you'd like more information or guidance, there are instructions in
the [Vega documentation](https://docs.testnet.vega.xyz/docs/wallet/).

If you'd prefer to request tokens from the contracts directly, there are
instructions in
the [testnet bridge tools repo README](https://github.com/vegaprotocol/Public_Test_Bridge_Tools/blob/master/docs/mew.md)
.

## Support

**[Documentation](https://docs.testnet.vega.xyz)**

Get API reference documentation, learn more about how Vega works, and explore
sample scripts for API trading

**[Walet documentation](https://docs.testnet.vega.xyz/docs/wallet/)**

**[Nolt](https://vega-testnet.nolt.io/)**

Raise issues and see what others have raised.

**[Discord](https://vega.xyz/discord)**

Ask us for help, find out about scheduled open office hours, and keep up with
Vega generally.

## About Vega

[Vega][vega-website] is a protocol for creating and trading derivatives on a
fully decentralised network. The network, secured with proof-of-stake, will
facilitate fully automated, end-to-end margin trading and execution of complex
financial products. Anyone will be able to build decentralised markets using the
protocol.

Read more at [https://vega.xyz][vega-website].

[vega-website]: https://vega.xyz
