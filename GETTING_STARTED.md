# Getting started

## How to install and run Vega Wallet
These instructions are written to be used in command line. Below, in the snippets, you'll see commands in blue text. Copy those instructions and paste them into your command line interface.

### Download
Download and save the zip file from Releases. We suggest you keep track of where you've saved the file, because that's where the command line interface will look for it.
https://github.com/vegaprotocol/go-wallet/releases

**For MacOS:**

Download `vegawallet-darwin-amd64.zip`

When you open the file, you may need to change your system preferences for this specific instance, in order to run Vega Wallet. If you open the file from downloads, you may get a message saying “vegawallet-darwin-amd64” cannot be opened because it is from an unidentified developer.

Click on the `(?)` help button, which will open a window that links you to the `System Preferences`, and instructs you how to allow this software to run.

You’ll need to go to `System Preferences` > `Security & Privacy` > `General`, and choose `Open Anyway`.

**For Windows:**

Download `vegawallet-windows-amd64.zip`

You may need to change your system preferences for this specific instance, in order to run Vega Wallet. If you open the file from downloads, you may get a message from Windows Defender saying “vegawallet-windows-amd64” cannot be opened because it is from an unidentified developer.

Click on the (More info) text, which will reveal a button to "Run anyway".

**For Linux:**

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
vegawallet
```
to execute the program.

> Tip: You can see a list of available commands by running  `./vegawallet -h` on MacOS and Linux, or `vegawallet -h` on Windows.

### Initialise the program

The `init` command (below) will initialise the configuration. A configuration file will be stored in your home folder, in a folder called `.vega`.

**MacOS & Linux**

```console
./vegawallet init
```
**Windows**

```console
vegawallet init
```

> Tip: If you want to specify a custom Vega home folder, it will not go into the default path, but a folder you choose to create. If you want to create a new config for a new wallet, or test or isolate it, you should specify the `--vega-home` flag.


### Create name and passphrase
Next, **create a user name and passphrase** for your Wallet, and **create a public and private key** (key generate).

Replace "YOUR_CUSTOM_USERNAME" (below) with your chosen username:

**MacOS & Linux**

```console
./vegawallet key generate --name "YOUR_CUSTOM_USERNAME"
```

**Windows**

```console
vegawallet key generate --name "YOUR_CUSTOM_USERNAME"
```

It will then prompt you to **input a passphrase**, and then **confirm that passphrase**. You'll use this username and passphrase to login to Vega Console. (Instructions on connecting to Console are below.)

The key generate command in that instruction will generate public and private keys for the wallet, at the same time as creating a user name.

You’ll see an output with a "mnemonic" and a public and private key. DO NOT SHARE YOUR MNEMONIC OR YOUR PRIVATE KEY.

**The mnemonic acts as a backup key, from which the wallet can restore all your keys.** Keep it safe and secret. If lost, you won't be able to retrieve your keys. Anyone who has this mnemonic will be able to use your keys.

You don’t need to save your private key, as you’ll be able to retrieve it from your Wallet in the future.

#### Give each new key a nickname/alias

When creating a key, you can give an alias by adding a metadata named `name`.

**MacOS & Linux**

```sh
./vegawallet key generate --name "YOUR_CUSTOM_USERNAME" --metas "name:CHOOSE_CUSTOM_ALIAS_FOR_KEY"`
```

**Windows**

```sh
vegawallet key generate --name "YOUR_CUSTOM_USERNAME" --metas "name:CHOOSE_CUSTOM_ALIAS_FOR_KEY"
```

#### Give an existing key a nickname/alias

**MacOS & Linux**

```sh
./vegawallet key annotate --metas="name:CHOOSE_CUSTOM_ALIAS_FOR_KEY" --name="YOUR_CUSTOM_USERNAME" --pubkey="REPLACE_THIS_WITH_YOUR_PUBLIC_KEY"
```

**Windows**

```sh
vegawallet key annotate --metas="name:CHOOSE_CUSTOM_ALIAS_FOR_KEY" --name="YOUR_CUSTOM_USERNAME" --pubkey="REPLACE_THIS_WITH_YOUR_PUBLIC_KEY"
```

> Tip: You can also use the annotate command to tag a key with other data you might want, using a property name and a value. This will be useful for developing with Vega Wallet in the future.

## Run the Wallet service
Now, **connect your wallet to the Vega testnet (Fairground) nodes**. To trade, run the wallet and **start the Vega Console** with the command below. (You'll need collateral to trade, and you can deposit it through Vega Console, once you're connected.)

**MacOS & Linux**

```console
./vegawallet service run --console-proxy
```
**Windows**

```console
vegawallet service run --console-proxy
```

> Tip: If you're running an ad/tracker blocker, and you're getting errors, it may be blocking the node from connecting. Try allowlisting lb.testnet.vega.xyz and refreshing.

> Tip: To terminate the process, such as if you want to run other commands in Wallet, use ctrl+c.

> Tip: See a full list of available commands in [the cmd readme](/cmd/README.md).

### Create and deposit testnet tokens
Now you'll need to **deposit Ropsten Ethereum-based tokens** to start trading.

You can create and deposit assets directly through the proxy Console via Wallet.

If you'd like more information or guidance, there are instructions in the [Vega documentation](https://docs.fairground.vega.xyz/docs/wallet/).

If you'd prefer to request tokens from the contracts directly, there are instructions in the [testnet bridge tools repo readme](https://github.com/vegaprotocol/Public_Test_Bridge_Tools/blob/master/docs/mew.md).
