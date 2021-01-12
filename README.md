# Vega Wallet

`vegawallet` is the command line interface for running a Wallet service, implemented in Go. It is used to sign transactions for use on [Vega](#about-vega). Vega Wallet creates and manages ed25519 keypairs for one or more wallets.

## How to install and run Vega Wallet
These instructions are written to be used in command line. Below, in the snippets, you'll see commands in blue text. Copy those instructions and paste them into your command line interface.

### Download
Download and save the zip file from Releases. We suggest you keep track of where you've saved the file, because that's where the command line interface will look for it. 
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

> Tip: You'll need to run the commands from the directory you've saved the wallet file in. Use the command `pwd` to find out where your terminal is looking in the file system. Use the command `cd` and the path/to/wallet/directory to tell the command line where to find the file. 

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

### Create wallet ID and passphrase
Next, **create a wallet ID and passphrase** for your Wallet, and **create a public and private key** (genkey). 

Replace "YOUR_CUSTOM_ID" (below) with your chosen wallet ID (username):

```console
./vegawallet genkey -n "YOUR_CUSTOM_ID"
```
It will then prompt you to **input a passphrase**, and then **confirm that passphrase**. You'll use this wallet ID and passphrase to login to Vega Console. (Instructions on connecting to Console are below.)

The genkey command in that instruction will generate public and private keys for the wallet, at the same time as creating a wallet ID.

You’ll see an output with your public and private key. DO NOT SHARE YOUR PRIVATE KEY. You don’t need to save this information anywhere, as you’ll be able to retrieve it from your Wallet in the future.

> Tip: You can see a list of available commands by running  ```./vegawallet -h```

## Run the Wallet service
Now, **connect your Wallet to the Testnet nodes**. The `init` command (below) will initialise the configuration. A configuration file will be stored in your home folder, in a folder called `.vega`.

```console
./vegawallet service init
```

> Tip: If you want to specify a root-path, it will not go into the default path, but a folder you choose to create. If you want to create a new config for a new wallet, or test or isolate it, you should specify the root path.

Next: To trade, run the wallet and **start the Vega Console** with the command below. (You'll need collateral to trade, and you can deposit it through Vega Console, once you're connected.)

```console
./vegawallet service run -p
```

> Tip: If you're running an ad/tracker blocker, and you're getting errors, it may be blocking the node from connecting. Try allowlisting lb.testnet.vega.xyz and refreshing.

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
