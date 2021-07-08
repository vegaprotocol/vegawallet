# Wallet CLI

## Root flags

By default, the wallet will be stored at a specific location. If you want to
specify a different location for test or isolation purposes, use
the ``--root-path`` flag to do so.

## List of available commands

There are 3 ways to list the available commands

```console
vegawallet
vegawallet -h
vegawallet help
```

## Create a wallet

To create a new wallet, generate your first key pair using
the following command:

```console
vegawallet key generate --name "YOUR_USERNAME"
```

The `--name` flag sets the name of your wallet.

It will then prompt you to input a passphrase, and then confirm that passphrase.
You'll use this username and passphrase to login to Vega Console. You can also
specify the passphrase with the ``--passphrase`` flag.

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

```console
vegawallet import --name "YOUR_USERNAME" --mnemonic "YOUR_MNEMONIC"
```

It will then prompt you to input a passphrase, and then confirm that passphrase.
You'll use this username and passphrase to login to Vega Console. You can also
specify the passphrase with the ``--passphrase`` flag.

This command is only able to import the wallet from which you can re-generate
your key pairs.

## Generate a key pair

To generate a key pair on the given wallet, use the following command:

```console
vegawallet key generate --name "YOUR_USERNAME"
```

It will then prompt you to input a passphrase. You can also specify the
passphrase with the ``--passphrase`` flag.

If the wallet does not exist, it will automatically create one. See
["Create a wallet"](#create-a-wallet) for more information.

You have the opportunity to attach metadata to your key with the ``--meta``
flag (more on this below).

## Add metadata to your keys

For better key management, you may want to add metadata to your key pairs. This
is done with the following command:

```console
vegawallet key meta --name "YOUR_USERNAME" --meta "key1:value1;key2:value2" --pubkey "YOUR_HEX_PUBLIC_KEY"
```

An item of metadata is represented as a key-value property.

### Give an alias to a key

You can give to each key pair a nickname/alias with a meta `name`. For example:

```console
vegawallet key meta --name "YOUR_USERNAME" --meta "name:my-alias" --pubkey "YOUR_HEX_PUBLIC_KEY"
```

### Important

This command does not insert the new metadata into the existing ones, **it
replaces them**. If you want to keep the previous metadata, ensure to add them
to your update.

## Tainting a key pair

You may want to prevent the use of a key by "tainting" it with
the following command:

```console
vegawallet key taint --name "YOUR_NAME" --pubkey "YOUR_HEX_PUBIC_KEY"
```

It will then prompt you to input a passphrase. You can also specify the
passphrase with the ``--passphrase`` flag.

## List the key pairs

To list your key pairs, use the following command:

```console
vegawallet key list --name "YOUR_NAME"
```

It will then prompt you to input a passphrase. You can also specify the
passphrase with the ``--passphrase`` flag.

### Important

This will also return the private key. **Never expose this command or its
content to the outside world.**

## Sign and verify messages

To sign and verify any kind of base-64 encoded messages, use the following
commands:

```console
vegawallet sign --name "YOUR_NAME" --pubkey "YOUR_HEX_PUBIC_KEY" --message "c3BpY2Ugb2YgZHVuZQo="
vegawallet verify --name "YOUR_NAME" --pubkey "YOUR_HEX_PUBIC_KEY" --message "c3BpY2Ugb2YgZHVuZQo=" --signature "76f978asd6fa8s76f"
```

It will then prompt you to input a passphrase. You can also specify the
passphrase with the ``--passphrase`` flag.

## Initialise the service

Before using the service, you need to initialise it with the following command:

```console
vegawallet service init
```

This creates the configuration file and RSA keys needed by the service to
operate.

## Run the service

Once the service has been initialised, you can run the wallet with the following
command:

```console
vegawallet service run --console-proxy
```

To terminate the process, such as if you want to run other commands in Wallet,
use `ctrl+c`.

### Important

If you're running an ad/tracker blocker, and you're getting errors, it may be
blocking the node from connecting. Try `allowlisting lb.testnet.vega.xyz` and
refreshing.
