package wallet

import "code.vegaprotocol.io/go-wallet/wallet/crypto"

func NewKeypair(algo crypto.SignatureAlgorithm, pub, priv []byte) Keypair {
	return Keypair{
		Algorithm: algo,
		pubBytes:  pub,
		privBytes: priv,
	}
}
