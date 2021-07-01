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

type LegacyKeyPair struct {
	Pub       string                    `json:"pub"`
	Priv      string                    `json:"priv,omitempty"`
	Algorithm crypto.SignatureAlgorithm `json:"algo"`
	Tainted   bool                      `json:"tainted"`
	MetaList  []Meta                    `json:"meta"`

	// byte version of the public and private keys
	// not being marshalled/sent over the network
	// or saved into the wallet file.
	pubBytes  []byte
	privBytes []byte
}

func GenKeyPair(algorithm string) (*LegacyKeyPair, error) {
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
	return &LegacyKeyPair{
		Priv:      hex.EncodeToString(privBytes),
		Pub:       hex.EncodeToString(pubBytes),
		Algorithm: algo,
		privBytes: privBytes,
		pubBytes:  pubBytes,
	}, err
}

func (k *LegacyKeyPair) PublicKey() string {
	return k.Pub
}

func (k *LegacyKeyPair) PrivateKey() string {
	return k.Priv
}

func (k *LegacyKeyPair) IsTainted() bool {
	return k.Tainted
}

func (k *LegacyKeyPair) Meta() []Meta {
	return k.MetaList
}

func (k *LegacyKeyPair) AlgorithmVersion() uint32 {
	return k.Algorithm.Version()
}

func (k *LegacyKeyPair) AlgorithmName() string {
	return k.Algorithm.Name()
}

func (k *LegacyKeyPair) Taint() error {
	if k.Tainted {
		return ErrPubKeyAlreadyTainted
	}

	k.Tainted = true
	return nil
}

func (k *LegacyKeyPair) SignAny(data []byte) ([]byte, error) {
	return k.Algorithm.Sign(k.privBytes, data)
}

func (k *LegacyKeyPair) VerifyAny(data, sig []byte) (bool, error) {
	return k.Algorithm.Verify(k.pubBytes, data, sig)
}

func (k *LegacyKeyPair) Sign(data []byte) (*commandspb.Signature, error) {
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

func (k *LegacyKeyPair) DeepCopy() *LegacyKeyPair {
	copiedK := *k
	return &copiedK
}

// ToPublicKey ensures the sensitive information doesn't leak outside.
func (k *LegacyKeyPair) ToPublicKey() *LegacyPublicKey {
	return &LegacyPublicKey{
		Pub:       k.Pub,
		Algorithm: k.Algorithm,
		Tainted:   k.Tainted,
		MetaList:  k.MetaList,
	}
}

func (k *LegacyKeyPair) MarshalJSON() ([]byte, error) {
	type alias LegacyKeyPair
	aliasKeypair := (*alias)(k)
	return json.Marshal(aliasKeypair)
}

func (k *LegacyKeyPair) UnmarshalJSON(data []byte) error {
	type alias LegacyKeyPair
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
