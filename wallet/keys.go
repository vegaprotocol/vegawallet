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

type KeyRing []Keypair

func NewKeyRing() KeyRing {
	return []Keypair{}
}

func (r KeyRing) FindPair(pubKey string) (Keypair, error) {
	for i := range r {
		if r[i].Pub == pubKey {
			return r[i], nil
		}
	}
	return Keypair{}, ErrPubKeyDoesNotExist
}

func (r *KeyRing) Upsert(pair Keypair) {
	for i := range *r {
		if (*r)[i].Pub == pair.Pub && (*r)[i].Priv == pair.Priv {
			(*r)[i] = pair
			return
		}
	}

	*r = append(*r, pair)
}

// GetPubKeys copy all keys so we do not propagate private keys
func (r KeyRing) GetPubKeys() []Keypair {
	pairs := make([]Keypair, 0, len(r))
	for _, keyPair := range r {
		pairs = append(pairs, keyPair.SecureCopy())
	}
	return pairs
}

type Keypair struct {
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

func GenKeypair(algorithm string) (*Keypair, error) {
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
	return &Keypair{
		Priv:      hex.EncodeToString(privBytes),
		Pub:       hex.EncodeToString(pubBytes),
		Algorithm: algo,
		privBytes: privBytes,
		pubBytes:  pubBytes,
	}, err
}

func (k *Keypair) MarshalJSON() ([]byte, error) {
	k.Pub = hex.EncodeToString(k.pubBytes)
	k.Priv = hex.EncodeToString(k.privBytes)
	type alias Keypair
	aliasKeypair := (*alias)(k)
	return json.Marshal(aliasKeypair)
}

func (k *Keypair) UnmarshalJSON(data []byte) error {
	type alias Keypair
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

func (k *Keypair) Taint() error {
	if k.Tainted {
		return ErrPubKeyAlreadyTainted
	}

	k.Tainted = true
	return nil
}

func (k *Keypair) Sign(marshalledData []byte) (*commandspb.Signature, error) {
	sig, err := k.Algorithm.Sign(k.privBytes, marshalledData)
	if err != nil {
		return nil, err
	}

	return &commandspb.Signature{
		Value:   hex.EncodeToString(sig),
		Algo:    k.Algorithm.Name(),
		Version: k.Algorithm.Version(),
	}, nil
}

func (k *Keypair) DeepCopy() Keypair {
	copiedK := *k
	return copiedK
}

// SecureCopy ensures the sensitive information doesn't leak outside.
func (k *Keypair) SecureCopy() Keypair {
	copiedK := *k
	copiedK.Priv = ""
	copiedK.privBytes = []byte{}
	return copiedK
}
