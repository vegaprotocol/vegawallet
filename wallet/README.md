Wallet
======

This package provides the basic cryptography to sign vega transactions, and a basic key management system: `wallet service`. It can be run alongside the core, but is not required for the operation of a Vega node, and API clients are free to implement their own transaction signing.

A wallet takes the form of a file saved on the file system and is encrypted using the passphrase chosen by the user.
A wallet is composed of a list of key pairs (Ed25519) used to sign transactions for the user of a wallet.

To read more about authentication & signing in Vega, see [Authentication & Signing](../design/authentication-and-signing.md)

## Generate configuration

The package provides a way to generate the configuration of the service before starting it, it can be used through the vega command line like so:

```shell
vega wallet service init --genrsakey -f
```

Where `--genrsakey` generates an RSA key that will be used to sign your JWT (JSON web token), and `-f` overwrites any existing configuration files (if found).
In short: this command will generate the RSA key, and the configuration files required by the wallet service.

Start the vega wallet service with:

```shell
vega wallet service run
```

## Create a wallet

Creating a wallet is done using a name and passphrase. If a wallet already exists, the action is aborted. New wallets are marshalled, encrypted (using the passphrase) and saved to a file on the file system.
A session and accompanying JWT is created, and the JWT is returned to the user.

* Request:

  ```json
  {
    "wallet": "walletname",
    "passphrase": "supersecret"
  }
  ```
* Command:

  ```shell
  curl -s -XPOST -d 'requestjson' http://127.0.0.1:1789/api/v1/wallets
  ```
* Response:

  ```json
  {
    "token": "verylongJWT"
  }
  ```

## Logging in to a wallet

Logging in to a wallet is done using the wallet name and passphrase.
The operation fails should the wallet not exist, or if the passphrase used is incorrect (i.e. the passphrase cannot be used to decrypt the wallet).
On success, the wallet is loaded, a session is created and a JWT is returned to the user.

* Request:

  ```json
  {
    "wallet": "walletname",
    "passphrase": "supersecret"
  }
  ```
* Command:

  ```shell
  curl -s -XPOST -d 'requestjson' http://127.0.0.1:1789/api/v1/auth/token
  ```
* Response:

  ```json
  {
    "token": "verylongJWT"
  }
  ```

## Logging out from a wallet

Using the JWT returned when logging in, the session is recovered and removed from the service. The wallet can no longer be accessed using the token from this point on.

* Request: n/a
* Command:

  ```shell
  curl -s -XDELETE -H 'Authorization: Bearer verylongJWT' http://127.0.0.1:1789/api/v1/auth/token
  ```
* Response:

  ```json
  {
    "success": true
  }
  ```

## List keys

Users can list all their public keys (with taint status, and metadata), if they provide the correct JWT. The service extracts the session from this token, and uses it to fetch the relevant wallet information to send back to the user.

* Request: n/a
* Command:

  ```shell
  curl -s -XGET -H "Authorization: Bearer verylongJWT" http://127.0.0.1:1789/api/v1/keys
  ```
* Response:

  ```json
  {
    "keys": [
      {
        "pub": "1122aabb...",
        "algo": "ed25519",
        "tainted": false,
        "meta": [
          {
            "key": "somekey",
            "value": "somevalue"
          }
        ]
      }
    ]
  }
  ```

## Generate a new key pair

The user submits a valid JWT, and a passphrase. We recover the session of the user, and attempt to open the wallet using the passphrase. If the JWT is invalid, the session could not be recovered, or the wallet could not be opened, an error is returned.
If all went well, a new key pair is generated, saved in the wallet, and the public key is returned.

* Request:

  ```json
  {
    "passphrase": "supersecret",
    "meta": [
      {
        "key": "somekey",
        "value": "somevalue"
      }
    ]
  }
  ```
* Command:

  ```shell
  curl -s -XPOST -H 'Authorization: Bearer verylongJWT' -d 'requestjson' http://127.0.0.1:1789/api/v1/keys
  ```
* Response:

  ```json
  {
    "key": {
      "pub": "1122aabb...",
      "algo": "ed25519",
      "tainted": false,
      "meta": [
        {
          "key": "somekey",
          "value": "somevalue"
        }
      ]
    }
  }
  ```

## Sign a transaction

Sign a transaction using the specified keypair.

* Request:

  ```json
  {
    "tx": "dGVzdGRhdGEK",
    "pubKey": "1122aabb...",
    "propagate": false
  }
  ```
* Command:

  ```shell
  curl -s -XPOST -H "Authorization: Bearer verylongJWT" -d 'requestjson' http://127.0.0.1:1789/api/v1/messages

  ```
* Response:

  ```json
  {
    "signedTx": {
      "data": "dGVzdGRhdGEK",
      "sig": "...",
      "pubKey": "1122aabb..."
    }
  }
  ```
## Taint a key

* Request:

  ```json
  {
    "passphrase": "supersecret"
  }
  ```
* Command:

  ```shell
  curl -s -XPUT -H "Authorization: Bearer verylongJWT" -d 'requestjson' http://127.0.0.1:1789/api/v1/keys/1122aabb/taint

  ```
* Response:

  ```json
  {
    "success": true
  }
  ```

## Update key metadata

Overwrite all existing metadata with the new metadata.

* Request:

  ```json
  {
    "passphrase": "supersecret",
    "meta": [
      {
        "key": "newkey",
        "value": "newvalue"
      }
    ]
  }
  ```
* Command:

  ```shell
  curl -s -XPUT -H "Authorization: Bearer verylongJWT" -d 'requestjson' http://127.0.0.1:1789/api/v1/keys/1122aabb/metadata

  ```
* Response:

  ```json
  {
    "success": true
  }
  ```
