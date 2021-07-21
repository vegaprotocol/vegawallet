# Vega Wallet

`vegawallet` is the command line interface for running a Wallet service,
implemented in Go. It is used to sign transactions for use
on [Vega](#about-vega). Vega Wallet creates and manages HD wallets with ed25519
key pairs.

### Use the [getting started instructions](/GETTING_STARTED.md) to install and run the wallet if you are new to command line (CLI).

## Building

```sh
cd go-wallet && make
```

Note: Building and compiling locally will install the wallet under the name
go-wallet, and not vegawallet. Thus, when reading the documentation,
replace `vegawallet` with `go-wallet`.

## Download

Download and save the zip file
from [Releases](https://github.com/vegaprotocol/go-wallet/releases). Keep track
of where you've saved the file, because that's where the CLI will look for it.

### MacOS

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

### Windows

Download `vegawallet-windows-amd64.zip`

You may need to change your system preferences for this specific instance, in
order to run Vega Wallet. If you open the file from downloads, you may get a
message from Windows Defender saying:

> "vegawallet-windows-amd64" cannot be opened because it is from an unidentified
> developer.

Click on the `(More info)` text, which will reveal a button to `Run anyway`.

### Linux

Download `vegawallet-linux-amd64.zip`

## Usage

**Important:** Before using the API and the commands, you will have to
initialise the program using the `init` command as
documented [here](cmd/README.md#initialise-the-program).

### Using the wallet commands

See a list of commands available in the wallet [here](cmd/README.md)

### Using the API

Using the API is documented [here](service/README.md).

## Support

**[Documentation](https://docs.fairground.vega.xyz)**

Get API reference documentation, learn more about how Vega works, and explore
sample scripts for API trading

**[Wallet documentation](https://docs.fairground.vega.xyz/docs/wallet/)**

Learn about how Vega interacts with wallets.

**[Nolt](https://vega-testnet.nolt.io/)**

Raise issues and see what others have raised.

**[Discord](https://vega.xyz/discord)**

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
