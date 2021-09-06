# Wallet CLI

## Root flags

By default, the wallet will be stored at a specific location. If you want to
specify a different location for test or isolation purposes, use
the `--home` flag to do so.

## List of available commands

There are 3 ways to list the available commands

```sh
vegawallet
vegawallet -h
vegawallet help
```

## Initialise the program

Before using the wallet, you need to initialise it with the following command:

```sh
vegawallet init
```

This creates the folders, the configuration file and RSA keys needed by the
wallet and its service to operate.

## Create a wallet

To create a new wallet, generate your first key pair using the following
command:

```sh
vegawallet key generate --name "YOUR_USERNAME"
```

The `--name` flag sets the name of your wallet.

It will then prompt you to input a passphrase, and then confirm that passphrase.
You'll use this username and passphrase to login to Vega Console. You can also
specify the passphrase with the `--passphrase-file` flag.

You have the opportunity to attach metadata to your key with the ``--meta``
flag (more on this below).

This command will generate a "mnemonic", along with a public and private key and
print it on the output.

### Important

**The mnemonic is very important as it acts as a backup key, from which the
wallet can restore all your keys.** As a result, it has to be kept safe and
secret. If lost, you won't be able to retrieve your keys. If stolen, the thief
will be able to use your keys.

Also, you’ll see an output with your public and private key. **Do not share your
private key.** You don’t need to save this information anywhere, as you’ll be
able to retrieve it with specific commands.

## Import a wallet

If you want to restore your wallet, use the following command:

```sh
vegawallet import --name "YOUR_USERNAME" --mnemonic-file "PATH_TO_YOUR_MNEMONIC"
```

The flag `--mnemonic-file` is used to locate the file that contains the
mnemonic.

It will then prompt you to input a passphrase, and then confirm that passphrase.
You'll use this username and passphrase to login to Vega Console. You can also
specify the passphrase with the `--passphrase-file` flag.

This command is only able to import the wallet from which you can re-generate
your key pairs.

## List registered wallets

If you want to list all the registered wallets, use the following command:

```sh
vegawallet list
```

## Generate a key pair

To generate a key pair on the given wallet, use the following command:

```sh
vegawallet key generate --name "YOUR_USERNAME"
```

It will then prompt you to input a passphrase. You can also specify the
passphrase with the `--passphrase-file` flag.

If the wallet does not exist, it will automatically create one. See
["Create a wallet"](#create-a-wallet) for more information.

You have the opportunity to attach metadata to your key with the ``--meta``
flag (more on this below).

## Add metadata to your keys

For better key management, you may want to add metadata to your key pairs. This
is done with the following command:

```sh
vegawallet key annotate --name "YOUR_USERNAME" --meta "key1:value1;key2:value2" --pubkey "YOUR_HEX_PUBLIC_KEY"
```

An item of metadata is represented as a key-value property.

### Give an alias to a key

You can give to each key pair a nickname/alias with a metadata `name`. For
example:

```sh
vegawallet key annotate --name "YOUR_USERNAME" --meta "name:my-alias" --pubkey "YOUR_HEX_PUBLIC_KEY"
```

### Important

This command does not insert the new metadata into the existing ones, **it
replaces them**. If you want to keep the previous metadata, ensure to add them
to your update.

## Tainting a key pair

You may want to prevent the use of a key by "tainting" it with the following
command:

```sh
vegawallet key taint --name "YOUR_NAME" --pubkey "YOUR_HEX_PUBIC_KEY"
```

It will then prompt you to input a passphrase. You can also specify the
passphrase with the `--passphrase-file` flag.

## Untainting a key pair

You may have tainted a key by mistake. If you want to untaint it, use the
following command:

```sh
vegawallet key untaint --name "YOUR_NAME" --pubkey "YOUR_HEX_PUBIC_KEY"
```

It will then prompt you to input a passphrase. You can also specify the
passphrase with the `--passphrase-file` flag.

### Important

If you tainted a key for security reasons, you should not untaint it.

## List the key pairs

To list your key pairs, use the following command:

```sh
vegawallet key list --name "YOUR_NAME"
```

It will then prompt you to input a passphrase. You can also specify the
passphrase with the `--passphrase-file` flag.

### Important

This will also return the private key. **Never expose this command or its
content to the outside world.**

## Sign and verify messages

To sign and verify any kind of base-64 encoded messages, use the following
commands:

```sh
vegawallet sign --name "YOUR_NAME" --pubkey "YOUR_HEX_PUBIC_KEY" --message "c3BpY2Ugb2YgZHVuZQo="
vegawallet verify --pubkey "YOUR_HEX_PUBIC_KEY" --message "c3BpY2Ugb2YgZHVuZQo=" --signature "76f978asd6fa8s76f"
```

It will then prompt you to input a passphrase. You can also specify the
passphrase with the `--passphrase-file` flag.

## Run the service

Once the service has been initialised, you can run the wallet with the following
command:

```sh
vegawallet service run
```

To run the wallet and open up a local version of Vega Console, the trading UI,
use the following command:

```sh
vegawallet service run --console-proxy
```

To terminate the process, such as if you want to run other commands in Wallet,
use `ctrl+c`.

### Ad-blockers

If you're running an ad/tracker blocker, and you're getting errors, it may be
blocking the node from connecting. Try allowlisting `lb.testnet.vega.xyz` and
refreshing.
