# Vega Wallet

`vegawallet` is the command line interface for running a Wallet service, implemented in Go. It is uses to sign transactions for use on [Vega](#about-vega).

# How to install and run Vega Wallet 

Download the executable file from Releases: https://github.com/vegaprotocol/go-wallet/releases 

After you've set up the Vega wallet, you can access it by opening the executable file, and running the commands you need. 

**If you’re using a Mac:**

Download `vegawallet-darwin-amd64`

It’s likely that you’ll need to change your system preferences for this specific instance, in order to run Vega Wallet. If you open the file from downloads, you may get a message saying “vegawallet-darwin-amd64” cannot be opened because it is from an unidentified developer.

Click on the (?) help button, which will open a window that links you to the System Preferences, and instructs you how to allow this software to run. 

You’ll then need to try to open it again, and go to System Preferences > Security & Privacy > General, and choose “Open Anyway”. 

**If you’re using Windows:**

Download `vegawallet-windows-amd64`

**If you’re using Linux:** 

Download `vegawallet-linux-amd64`

You need to have the `gtk+-3.0` and `webkit2gtk-4.0` dependencies installed.

And if you want to add an example (this would be a different command in different linux distribution). For Ubuntu specifically, the command would be `sudo apt-get install gtk+-3.0 webkit2gtk-4.0`. 

These instructions are written to be used in a terminal.  (this will probably be worded more elegantly)

Find the wallet file on your system using the command 

**Mac:** `./vegawallet-darwin-amd64`
**Linux:** `./vegawallet-linux-amd64`
**Windows:** `./vegawallet-windows-amd64`

To be able to run the file, you'll need to make it executable. 

Tip: On a Mac, you might then get the message that this is a security issue. The above instructions (and what your Mac will guide you to do) will only change the permissions for this file.  It will guide you to system preferences > security and privacy > general (but the help section should give you info you need). 

Make the file executable

**Mac:** 

```
chmod +x vegawallet-darwin-amd64
```

**Linux:** 

```
chmod +x vegawallet-linux-amd64
```

**Windows:** 

```
chmod +x vegawallet-windows-amd64
```

*Tip:* You can use the tab key to auto-fill the name of the file, after you type the first few characters. 

Rename the file. (This is optional, but all instructions from here will assume you’ve renamed it to vegawallet.) 

**Mac:** 
```
mv vegawallet-darwin-amd64 vegawallet
```

**Linux:** 
```
mv vegawallet-linux-amd64 vegawallet
```

**Windows:**
```
mv vegawallet-windows-amd64 vegawallet
```

*Tip:* You can see a list of available commands by running
```
./vegawallet -h
```

Next, create a user name and password for your Wallet, and create a public and private key (genkey):

```
./vegawallet genkey --name pick-a-username
``` 

It will then ask you to input a password, and then confirm that password. 

The genkey command in that instruction will generate public and private keys for the wallet, at the same time as creating a user name. 

You’ll see an output with your public and private key. DO NOT SHARE YOUR PRIVATE KEY. You don’t need to save this information anywhere, as you’ll be able to retrieve it from your Wallet in the future. 

Now, connect your Wallet to the Testnet nodes and UI. The ‘init’ command (below) will initialise the configuration. A configuration file will be stored in your home folder, in a folder called `.vega`.

```
./vegawallet service init
```

*Tip:* If you want to specify a root-path, it will not go into the default path, but a folder you choose to create. If you want to create a new config for a new wallet, or test or isolate it, you should specify the root path.

If you want to trade using the APIs, use the command 

```
./vegawallet service run
```

Otherwise, you can connect to a Console proxy so you can trade via the UI.

Start the Vega console proxy and open the console in the default browser:

```
./vegawallet service run -p
```

Then deposit funds! (instructions tbd) 


## Support

[Documentation](https://docs.testnet.vega.xyz) 
Get API reference documentation, learn more about how Vega works, and explore sample scripts for API trading

[Nolt](https://vega-testnet.nolt.io/)
Raise issues, see what others have raised. 

[Discord](https://discord.gg/bkAF3Tu) 
Ask us for help, find out about scheduled open office hours, and keep up with Vega generally 










------------------ 

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
The service can be run in 3 different ways.

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
{"level":"info","ts":1605554589.694528,"caller":"cmd/console.go:41","msg":"starting console proxy","proxy.address":"localhost:8080","address":"testnet.vega.trading"}
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
