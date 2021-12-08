# Wallet API

## Authentication

### Logging in to a wallet

`POST api/v1/auth/token`

Logging in to a wallet is done using the wallet name and passphrase. The
operation fails if the wallet not exist, or if the passphrase used is incorrect.
On success, the wallet is loaded, a session is created and a JWT is returned to
the user.

#### Example

##### Request

```json
{
  "wallet": "your_wallet_name",
  "passphrase": "super-secret"
}
```

##### Command

```sh
curl -s -XPOST -d 'YOUR_REQUEST' http://127.0.0.1:1789/api/v1/auth/token
```

##### Response

```json
{
  "token": "abcd.efgh.ijkl"
}
```

### Logging out from a wallet

`DELETE api/v1/auth/token`

Using the JWT returned when logging in, the session is recovered and removed
from the service. The wallet can no longer be accessed using the token from this
point on.

#### Example

##### Command

```sh
curl -s -XDELETE -H 'Authorization: Bearer abcd.efgh.ijkl' http://127.0.0.1:1789/api/v1/auth/token
```

##### Response

```json
{
  "success": true
}
```

## Network management

### Get current network configuration

`GET api/v1/network`

### Example

#### Command

```sh
curl -s -XPOST -d 'YOUR_REQUEST' http://127.0.0.1:1789/api/v1/network
```

#### Response

```json
{
  "network": {
    "name": "mainnet"
  }
}
```

## Wallet management

### Create a wallet

`POST api/v1/wallets`

Creating a wallet is done using a name and passphrase. If a wallet with the same
name already exists, the action is aborted. The new wallets is encrypted (using
the passphrase) and saved to a file on the file system. A session and
accompanying JWT is created, and the JWT is returned to the user.

#### Example

##### Request

```json
{
  "wallet": "your_wallet_name",
  "passphrase": "super-secret"
}
```

##### Command

```sh
curl -s -XPOST -d 'YOUR_REQUEST' http://127.0.0.1:1789/api/v1/wallets
```

##### Response

```json
{
  "token": "abcd.efgh.ijkl"
}
```

### Import a wallet

`POST api/v1/wallets/import`

Import a wallet is done using a name, a passphrase, and a recoveryPhrase. If a wallet
with the same name already exists, the action is aborted. The imported wallet is
encrypted (using the passphrase) and saved to a file on the file system. A
session and accompanying JWT is created, and the JWT is returned to the user.

#### Example

##### Request

```json
{
  "wallet": "your_wallet_name",
  "passphrase": "super-secret",
  "recoveryPhrase": "my twenty four words recovery phrase"
}
```

##### Command

```sh
curl -s -XPOST -d 'YOUR_REQUEST' http://127.0.0.1:1789/api/v1/wallets
```

##### Response

```json
{
  "token": "abcd.efgh.ijkl"
}
```

## Key management

### Generate a key pair

`POST api/v1/keys`

**Authentication required.**

It generates a new key pair into the logged wallet, and returns the generated
public key.

#### Example

##### Request

```json
{
  "passphrase": "super-secret",
  "meta": [
    {
      "key": "somekey",
      "value": "somevalue"
    }
  ]
}
```

##### Command

```sh
curl -s -XPOST -H 'Authorization: Bearer abcd.efgh.ijkl' -d 'YOUR_REQUEST' http://127.0.0.1:1789/api/v1/keys
```

##### Response

```json
{
  "key": {
    "pub": "1122aabb",
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

### List keys

`GET api/v1/keys`

**Authentication required.**

Users can list all the public keys (with taint status, and metadata) of the
logged wallet.

#### Example

##### Command

```sh
curl -s -XGET -H "Authorization: Bearer abcd.efgh.ijkl" http://127.0.0.1:1789/api/v1/keys
```

##### Response

```json
{
  "keys": [
    {
      "pub": "1122aabb",
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

### Describe a key pair

`GET api/v1/keys/:keyid`

**Authentication required.**

Return the information associated the public key `:keyid`, from the logged
wallet. The private key is not returned.

#### Example

##### Command

```sh
  curl -s -XPUT -H "Authorization: Bearer abcd.efgh.ijkl" -d 'YOUR_REQUEST' http://127.0.0.1:1789api/v1/keys/1122aabb
```

##### Response

```json
{
  "key": {
    "index": 1,
    "pub": "1122aabb"
  }
}
```

### Taint a key pair

`PUT api/v1/keys/:keyid/taint`

**Authentication required.**

Taint the key pair matching the public key `:keyid`, from the logged wallet. The
key pair must belong to the logged wallet.

#### Example

##### Request

```json
{
  "passphrase": "super-secret"
}
```

##### Command

```sh
  curl -s -XPUT -H "Authorization: Bearer abcd.efgh.ijkl" -d 'YOUR_REQUEST' http://127.0.0.1:1789/api/v1/keys/1122aabb/taint
```

##### Response

```json
{
  "success": true
}
```

### Annotate a key pair

`PUT api/v1/keys/:keyid/metadata`

**Authentication required.**

Annotating a key pair replace the metadata matching the public key `:keyid`,
from the logged wallet. The key pair must belong to the logged wallet.

#### Example

##### Request

```json
{
  "passphrase": "super-secret",
  "meta": [
    {
      "key": "newkey",
      "value": "newvalue"
    }
  ]
}
```

##### Command

```sh
  curl -s -XPUT -H "Authorization: Bearer abcd.efgh.ijkl" -d 'YOUR_REQUEST' http://127.0.0.1:1789/api/v1/keys/1122aabb/metadata
```

##### Response

```json
{
  "success": true
}
```

## Commands

### Sign a command

`POST api/v1/command`

**Authentication required.**

Sign a Vega command using the specified key pair, and returns the signed
transaction. The key pair must belong to the logged wallet.

#### Example

##### Request

```json
{
  "pubKey": "1122aabb",
  "propagate": true,
  "orderCancellation": {
    "marketId": "YESYESYES"
  }
}
```

##### Command

```sh
  curl -s -XPOST -H "Authorization: Bearer abcd.efgh.ijkl" -d 'YOUR_REQUEST' http://127.0.0.1:1789/api/v1/command
```

##### Response

```json
{
  "transaction": {
    "inputData": "dGVzdGRhdG9837420b4b3yb23ybc4o1ui23yEK",
    "signature": {
      "value": "7f6g9sf8f8s76dfa867fda",
      "algo": "vega/ed25519",
      "version": 1
    },
    "from": {
      "pubKey": "1122aabb"
    },
    "version": 1
  }
}
```

#### Propagate

In the request payload, when the `propagate` field can be set to true, the
wallet service send the transaction on your behalf to the registered nodes after
signing it successfully.

### Sign data

`POST api/v1/sign`

**Authentication required.**

Sign any base64-encoded data using the specified key pair, and returns the
signed transaction. The key pair must belong to the logged wallet.

#### Example

##### Request

```json
{
  "inputData": "dGVzdGRhdGEK==",
  "pubKey": "1122aabb"
}
```

##### Command

```sh
  curl -s -XPOST -H "Authorization: Bearer abcd.efgh.ijkl" -d 'YOUR_REQUEST' http://127.0.0.1:1789/api/v1/sign
```

##### Response

```json
{
  "hexSignature": "0xhafdsf86df876af",
  "base64Signature": "fad7h34k1jh3g413g=="
}
```

### Verify data

`POST api/v1/verify`

Verify any base64-encoded data using the specified public key, and returns the
confirmation.

#### Example

##### Request

```json
{
  "inputData": "dGVzdGRhdGEK==",
  "pubKey": "1122aabb"
}
```

##### Command

```sh
  curl -s -XPOST -H "Authorization: Bearer abcd.efgh.ijkl" -d 'YOUR_REQUEST' http://127.0.0.1:1789/api/v1/sign
```

##### Response

```json
{
  "hexSignature": "0xhafdsf86df876af",
  "base64Signature": "fad7h34k1jh3g413g=="
}
```
