package wallet

import (
	"encoding/hex"
	"encoding/json"
	"errors"

	"code.vegaprotocol.io/go-wallet/wallet/crypto"
	commandspb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/commands/v1"
)

var (
	ErrPubKeyDoesNotExist   = errors.New("public key does not exist")
	ErrPubKeyAlreadyTainted = errors.New("public key is already tainted")
)

type KeyRing []KeyPair

func NewKeyRing() KeyRing {
	return []KeyPair{}
}

func (r KeyRing) FindPair(pubKey string) (KeyPair, error) {
	for i := range r {
		if r[i].Pub == pubKey {
			return r[i], nil
		}
	}
	return KeyPair{}, ErrPubKeyDoesNotExist
}

func (r *KeyRing) Upsert(pair KeyPair) {
	for i := range *r {
		if (*r)[i].Pub == pair.Pub && (*r)[i].Priv == pair.Priv {
			(*r)[i] = pair
			return
		}
	}

	*r = append(*r, pair)
}

func (r KeyRing) GetPublicKeys() []PublicKey {
	pubKeys := make([]PublicKey, 0, len(r))
	for _, keyPair := range r {
		pubKeys = append(pubKeys, *keyPair.ToPublicKey())
	}
	return pubKeys
}

type KeyPair struct {
	Pub       string                    `json:"pub"`
	Priv      string                    `json:"priv,omitempty"`
	Algorithm crypto.SignatureAlgorithm `json:"algo"`
	Tainted   bool                      `json:"tainted"`
	Meta      []Meta                    `json:"meta"`

	// byte version of the public and private keys
	// not being marshalled/sent over the network
	// or saved into the wallet file.
	pubBytes  []byte
	privBytes []byte
}

func GenKeyPair(algorithm string) (*KeyPair, error) {
	algo, err := crypto.NewSignatureAlgorithm(algorithm)
	if err != nil {
		return nil, err
	}

	pub, priv, err := algo.GenKey()
	if err != nil {
		return nil, err
	}

	privBytes := priv.([]byte)
	pubBytes := pub.([]byte)
	return &KeyPair{
		Priv:      hex.EncodeToString(privBytes),
		Pub:       hex.EncodeToString(pubBytes),
		Algorithm: algo,
		privBytes: privBytes,
		pubBytes:  pubBytes,
	}, err
}

func (k *KeyPair) MarshalJSON() ([]byte, error) {
	k.Pub = hex.EncodeToString(k.pubBytes)
	k.Priv = hex.EncodeToString(k.privBytes)
	type alias KeyPair
	aliasKeypair := (*alias)(k)
	return json.Marshal(aliasKeypair)
}

func (k *KeyPair) UnmarshalJSON(data []byte) error {
	type alias KeyPair
	aliasKeypair := (*alias)(k)
	if err := json.Unmarshal(data, aliasKeypair); err != nil {
		return err
	}
	var err error
	k.pubBytes, err = hex.DecodeString(k.Pub)
	if err != nil {
		return err
	}
	k.privBytes, err = hex.DecodeString(k.Priv)
	return err
}

func (k *KeyPair) Taint() error {
	if k.Tainted {
		return ErrPubKeyAlreadyTainted
	}

	k.Tainted = true
	return nil
}

func (k *KeyPair) SignAny(data []byte) ([]byte, error) {
	return k.Algorithm.Sign(k.privBytes, data)
}

func (k *KeyPair) VerifyAny(data, sig []byte) (bool, error) {
	return k.Algorithm.Verify(k.pubBytes, data, sig)
}

func (k *KeyPair) Sign(data []byte) (*commandspb.Signature, error) {
	sig, err := k.Algorithm.Sign(k.privBytes, data)
	if err != nil {
		return nil, err
	}

	return &commandspb.Signature{
		Value:   hex.EncodeToString(sig),
		Algo:    k.Algorithm.Name(),
		Version: k.Algorithm.Version(),
	}, nil
}

func (k *KeyPair) DeepCopy() KeyPair {
	copiedK := *k
	return copiedK
}

// ToPublicKey ensures the sensitive information doesn't leak outside.
func (k *KeyPair) ToPublicKey() *PublicKey {
	return &PublicKey{
		Key:       k.Pub,
		Algorithm: k.Algorithm,
		Tainted:   k.Tainted,
		Meta:      k.Meta,
	}
}

type PublicKey struct {
	Key       string                    `json:"pub"`
	Algorithm crypto.SignatureAlgorithm `json:"algo"`
	Tainted   bool                      `json:"tainted"`
	Meta      []Meta                    `json:"meta"`
}
