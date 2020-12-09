# Vega Wallet

`vegawallet` is the command line interface for running a Wallet service, implemented in Go. It is used to sign transactions for use on [Vega](#about-vega). Vega Wallet creates and manages ed25519 keypairs for one or more wallets.

## How to install and run Vega Wallet
These instructions are written to be used in command line. Below, in the snippets, you'll see commands in blue text. Copy those instructions and paste them into your command line interface.

### Download
Download and save the zip file from Releases to the location on your computer you want to run it from: 
https://github.com/vegaprotocol/go-wallet/releases

**If you’re using MacOS:**

Download `vegawallet-darwin-amd64.zip`

When you open the file, you may need to change your system preferences for this specific instance, in order to run Vega Wallet. If you open the file from downloads, you may get a message saying “vegawallet-darwin-amd64” cannot be opened because it is from an unidentified developer.

Click on the (?) help button, which will open a window that links you to the System Preferences, and instructs you how to allow this software to run.

You’ll need to go to System Preferences > Security & Privacy > General, and choose “Open Anyway”.

**If you’re using Windows:**

Download `vegawallet-windows-amd64.zip`

You may need to change your system preferences for this specific instance, in order to run Vega Wallet. If you open the file from downloads, you may get a message from Windows Defender saying “vegawallet-darwin-amd64” cannot be opened because it is from an unidentified developer.

Click on the (More info) text, which will reveal a button to "Run anyway".

**If you’re using Linux:**

Download `vegawallet-linux-amd64.zip`

## Generate key pair and credentials

### Execute the program

> Tip: You can use the tab key to auto-fill the name of the file, after you type the first few characters.

**MacOS & Linux**
Open a new terminal. Type

```console
./vegawallet
```
to execute the program.

**Windows**

Open a new command prompt. Type

```console
./vegawallet.exe
```
to execute the program.

### Create name and passphrase
Next, **create a user name and passphrase** for your Wallet, and **create a public and private key** (genkey). 

Insert your chosen username where it says create-username below:

```console
./vegawallet genkey -n create-username
```
It will then prompt you to **input a passphrase**, and then **confirm that passphrase**. You'll use this username and passphrase to login to Vega Console. (Instructions on connecting to Console are below.)

The genkey command in that instruction will generate public and private keys for the wallet, at the same time as creating a user name.

You’ll see an output with your public and private key. DO NOT SHARE YOUR PRIVATE KEY. You don’t need to save this information anywhere, as you’ll be able to retrieve it from your Wallet in the future.

> Tip: You can see a list of available commands by running
    ```console
    ./vegawallet -h
    ```

## Run the Wallet service
Now, **connect your Wallet to the Testnet nodes**. The `init` command (below) will initialise the configuration. A configuration file will be stored in your home folder, in a folder called `.vega`.

```console
./vegawallet service init
```

> Tip: If you want to specify a root-path, it will not go into the default path, but a folder you choose to create. If you want to create a new config for a new wallet, or test or isolate it, you should specify the root path.

Connect to a Console proxy so you can trade via the UI. (You'll need collateral to trade, and you can deposit it through Vega Console, once you're connected.)

To trade, **start the Vega Console proxy** and open Console in the default browser with the command

```console
./vegawallet service run -p
```

Otherwise, if you prefer to **trade using the APIs**, use the command below. You'll need to deposit assets to trade, find a link to instructions below.

```console
./vegawallet service run
```

> Tip: To terminate the process, such as if you want to run other commands in Wallet, use ctrl+c.

### Create and deposit testnet tokens
Now you'll need to **deposit Ropsten Ethereum-based tokens** to start trading.

You can create and deposit assets directly through the proxy Console via Wallet.

If you'd like more information or guidance, there are instructions in the [Vega documentation](https://docs.testnet.vega.xyz/docs/wallet/).

If you'd prefer to request tokens from the contracts directly, there are instructions in the [testnet bridge tools repo readme](https://github.com/vegaprotocol/Public_Test_Bridge_Tools/blob/master/docs/mew.md).

### Use the wallet API
Using the API is documented [here](./wallet/README.md).

## Support

**[Documentation](https://docs.testnet.vega.xyz)**

Get API reference documentation, learn more about how Vega works, and explore sample scripts for API trading

**[Nolt](https://vega-testnet.nolt.io/)**

Raise issues and see what others have raised.

**[Discord](https://vega.xyz/discord)**

Ask us for help, find out about scheduled open office hours, and keep up with Vega generally.

## Building
```console
cd go-wallet && make
```

# About Vega
[Vega](https://vega.xyz) is a protocol for creating and trading derivatives on a fully decentralised network. The network, secured with proof-of-stake, will facilitate fully automated, end-to-end margin trading and execution of complex financial products. Anyone will be able to build decentralised markets using the protocol.

Read more at [https://vega.xyz](https://vega.xyz).
