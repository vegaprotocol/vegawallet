package wallet

import (
	"crypto/ed25519"
	"encoding/hex"

	"code.vegaprotocol.io/vegawallet/crypto"
)

type HDMasterKeyPair struct {
	publicKey  *key
	privateKey *key
	algo       crypto.SignatureAlgorithm
}

func NewHDMasterKeyPair(
	publicKey ed25519.PublicKey,
	privateKey ed25519.PrivateKey,
) (*HDMasterKeyPair, error) {
	algo, err := crypto.NewSignatureAlgorithm(crypto.Ed25519, 1)
	if err != nil {
		return nil, err
	}

	return &HDMasterKeyPair{
		publicKey: &key{
			bytes:   publicKey,
			encoded: hex.EncodeToString(publicKey),
		},
		privateKey: &key{
			bytes:   privateKey,
			encoded: hex.EncodeToString(privateKey),
		},
		algo: algo,
	}, nil
}

func (k *HDMasterKeyPair) PublicKey() string {
	return k.publicKey.encoded
}

func (k *HDMasterKeyPair) PrivateKey() string {
	return k.privateKey.encoded
}

func (k *HDMasterKeyPair) AlgorithmVersion() uint32 {
	return k.algo.Version()
}

func (k *HDMasterKeyPair) AlgorithmName() string {
	return k.algo.Name()
}

func (k *HDMasterKeyPair) SignAny(data []byte) ([]byte, error) {
	return k.algo.Sign(k.privateKey.bytes, data)
}

func (k *HDMasterKeyPair) Sign(data []byte) (*Signature, error) {
	sig, err := k.algo.Sign(k.privateKey.bytes, data)
	if err != nil {
		return nil, err
	}

	return &Signature{
		Value:   hex.EncodeToString(sig),
		Algo:    k.algo.Name(),
		Version: k.algo.Version(),
	}, nil
}
