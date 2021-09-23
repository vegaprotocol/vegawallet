package wallet

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"

	"code.vegaprotocol.io/go-wallet/crypto"
)

type HDKeyPair struct {
	index      uint32
	publicKey  *key
	privateKey *key
	meta       []Meta
	tainted    bool
	algo       crypto.SignatureAlgorithm
}

type key struct {
	bytes   []byte
	encoded string
}

func NewHDKeyPair(
	index uint32,
	publicKey ed25519.PublicKey,
	privateKey ed25519.PrivateKey,
) (*HDKeyPair, error) {
	algo, err := crypto.NewSignatureAlgorithm(crypto.Ed25519, 1)
	if err != nil {
		return nil, err
	}

	return &HDKeyPair{
		index: index,
		publicKey: &key{
			bytes:   publicKey,
			encoded: hex.EncodeToString(publicKey),
		},
		privateKey: &key{
			bytes:   privateKey,
			encoded: hex.EncodeToString(privateKey),
		},
		algo:    algo,
		meta:    nil,
		tainted: false,
	}, nil
}

func (k *HDKeyPair) Index() uint32 {
	return k.index
}

func (k *HDKeyPair) PublicKey() string {
	return k.publicKey.encoded
}

func (k *HDKeyPair) PrivateKey() string {
	return k.privateKey.encoded
}

func (k *HDKeyPair) IsTainted() bool {
	return k.tainted
}

func (k *HDKeyPair) Meta() []Meta {
	return k.meta
}

func (k *HDKeyPair) AlgorithmVersion() uint32 {
	return k.algo.Version()
}

func (k *HDKeyPair) AlgorithmName() string {
	return k.algo.Name()
}

func (k *HDKeyPair) Taint() error {
	if k.tainted {
		return ErrPubKeyAlreadyTainted
	}

	k.tainted = true
	return nil
}

func (k *HDKeyPair) Untaint() error {
	if !k.tainted {
		return ErrPubKeyNotTainted
	}

	k.tainted = false
	return nil
}

func (k *HDKeyPair) SignAny(data []byte) ([]byte, error) {
	if k.tainted {
		return nil, ErrPubKeyIsTainted
	}

	return k.algo.Sign(k.privateKey.bytes, data)
}

func (k *HDKeyPair) VerifyAny(data, sig []byte) (bool, error) {
	return k.algo.Verify(k.publicKey.bytes, data, sig)
}

func (k *HDKeyPair) Sign(data []byte) (*Signature, error) {
	if k.tainted {
		return nil, ErrPubKeyIsTainted
	}

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

func (k *HDKeyPair) DeepCopy() *HDKeyPair {
	copiedK := *k
	return &copiedK
}

// ToPublicKey ensures the sensitive information doesn't leak outside.
func (k *HDKeyPair) ToPublicKey() HDPublicKey {
	return HDPublicKey{
		Idx:       k.Index(),
		PublicKey: k.PublicKey(),
		Algorithm: Algorithm{
			Name:    k.algo.Name(),
			Version: k.algo.Version(),
		},
		Tainted:  k.tainted,
		MetaList: k.meta,
	}
}

type jsonHDKeyPair struct {
	Index      uint32    `json:"index"`
	PublicKey  string    `json:"public_key"`
	PrivateKey string    `json:"private_key"`
	Meta       []Meta    `json:"meta"`
	Tainted    bool      `json:"tainted"`
	Algorithm  Algorithm `json:"algorithm"`
}

func (k *HDKeyPair) MarshalJSON() ([]byte, error) {
	jsonKp := jsonHDKeyPair{
		Index:      k.index,
		PublicKey:  k.publicKey.encoded,
		PrivateKey: k.privateKey.encoded,
		Meta:       k.meta,
		Tainted:    k.tainted,
		Algorithm: Algorithm{
			Name:    k.algo.Name(),
			Version: k.algo.Version(),
		},
	}
	return json.Marshal(jsonKp)
}

func (k *HDKeyPair) UnmarshalJSON(data []byte) error {
	jsonKp := &jsonHDKeyPair{}
	if err := json.Unmarshal(data, jsonKp); err != nil {
		return err
	}

	algo, err := crypto.NewSignatureAlgorithm(jsonKp.Algorithm.Name, jsonKp.Algorithm.Version)
	if err != nil {
		return err
	}

	pubKeyBytes, err := hex.DecodeString(jsonKp.PublicKey)
	if err != nil {
		return err
	}

	privKeyBytes, err := hex.DecodeString(jsonKp.PrivateKey)
	if err != nil {
		return err
	}

	*k = HDKeyPair{
		index: jsonKp.Index,
		publicKey: &key{
			bytes:   pubKeyBytes,
			encoded: jsonKp.PublicKey,
		},
		privateKey: &key{
			bytes:   privKeyBytes,
			encoded: jsonKp.PrivateKey,
		},
		meta:    jsonKp.Meta,
		tainted: jsonKp.Tainted,
		algo:    algo,
	}

	return nil
}
