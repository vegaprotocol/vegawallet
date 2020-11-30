# Vega Wallet

`vegawallet` is the command line interface for running a Wallet service, implemented in Go. It is used to sign transactions for use on [Vega](#about-vega). Vega Wallet creates and manages ed25519 keypairs for one or more wallets.

## How to install and run Vega Wallet 
These instructions are written to be used in command line. Below, in the snippets, you'll see commands in black text. You can copy those command line instructions, and paste them into your command line interface. 

### Download 
Download and save the zip file from Releases to the location on your computer you want to run it from: https://github.com/vegaprotocol/go-wallet/releases 

**If you’re using MacOS:**

Download `vegawallet-darwin-amd64.zip`

When you open the file, it’s likely that you’ll need to change your system preferences for this specific instance, in order to run Vega Wallet. If you open the file from downloads, you may get a message saying “vegawallet-darwin-amd64” cannot be opened because it is from an unidentified developer.

Click on the (?) help button, which will open a window that links you to the System Preferences, and instructs you how to allow this software to run. 

You’ll need to go to System Preferences > Security & Privacy > General, and choose “Open Anyway”. 

**If you’re using Windows:**

Download `vegawallet-windows-amd64.zip`

It’s likely that you’ll need to change your system preferences for this specific instance, in order to run Vega Wallet. If you open the file from downloads, you may get a message from Windows Defender saying “vegawallet-darwin-amd64” cannot be opened because it is from an unidentified developer.

Click on the (More info) text, which will reveal a button to "Run anyway".  

**If you’re using Linux:** 

Download `vegawallet-linux-amd64.zip`

## Generate key pair and credentials

### Execute the program

*Tip:* You can use the tab key to auto-fill the name of the file, after you type the first few characters. 

**MacOS & Linux**
Open a new terminal. Type

```console
wallet@vega:~$ ./vegawallet
```
to execute the program. 

**Windows**

Open a new command prompt. Type

```console
wallet@vega:~$ ./vegawallet.exe
```
to execute the program. 

### Create name and passphrase
Next, create a user name and passphrase for your Wallet, and create a public and private key (genkey):

```console
wallet@vega:~$ ./vegawallet genkey -n [choose-a-username]

please enter passphrase:
``` 

It will then prompt you to input a passphrase, and then confirm that passphrase. 

The genkey command in that instruction will generate public and private keys for the wallet, at the same time as creating a user name. 

You’ll see an output with your public and private key. DO NOT SHARE YOUR PRIVATE KEY. You don’t need to save this information anywhere, as you’ll be able to retrieve it from your Wallet in the future. 

*Tip:* You can see a list of available commands by running
```console
wallet@vega:~$ ./vegawallet -h
```

## Run the Wallet service
Now, connect your Wallet to the Testnet nodes and UI. The `init` command (below) will initialise the configuration. A configuration file will be stored in your home folder, in a folder called `.vega`.

```console
wallet@vega:~$ ./vegawallet service init

{"level":"info","ts":1605554344.734188,"caller":"wallet/config.go:125","msg":"wallet service configuration generated successfully","path":"path/wallet-service-config.toml"}
{"level":"info","ts":1605554347.727988,"caller":"wallet/config.go:173","msg":"wallet rsa key generated successfully","path":"path/wallet_rsa"}
```

*Tip:* If you want to specify a root-path, it will not go into the default path, but a folder you choose to create. If you want to create a new config for a new wallet, or test or isolate it, you should specify the root path.

You'll need collateral to trade, but once you want to trade using the APIs, use the command 

```console
wallet@vega:~$ ./vegawallet service run

{"level":"info","ts":1587317545.61634,"logger":"wallet","caller":"wallet/service.go:147","msg":"starting wallet http server","address":"0.0.0.0:1789"}
```

Otherwise, you can connect to a Console proxy so you can trade via the UI.

Start the Vega Console proxy and open Console in the default browser:

```console
wallet@vega:~$ ./vegawallet service run -p

{"level":"info","ts":1605554589.694528,"caller":"cmd/console.go:41","msg":"starting console proxy","proxy.address":"localhost:8080","address":"dev.vega.trading"}
{"level":"info","ts":1587317545.61634,"logger":"wallet","caller":"wallet/service.go:147","msg":"starting wallet http server","address":"0.0.0.0:1789"}
```

*Tip:* To terminate the process, such as if you want to run other commands in Wallet, use ctrl+c. 

### Create and deposit testnet tokens
Now you'll need to deposit Ropsten Ethereum-based tokens to start trading. 

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
wallet@vega:~$ cd go-wallet && make
```

# About Vega
[Vega](https://vega.xyz) is a protocol for creating and trading derivatives on a fully decentralised network. The network, secured with proof-of-stake, will facilitate fully automated, end-to-end margin trading and execution of complex financial products. Anyone will be able to build decentralised markets using the protocol.

Read more at [https://vega.xyz](https://vega.xyz).
