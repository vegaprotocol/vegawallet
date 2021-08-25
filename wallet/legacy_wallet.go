package wallet

import (
	"code.vegaprotocol.io/go-wallet/crypto"
	commandspb "code.vegaprotocol.io/protos/vega/commands/v1"
)

type LegacyWallet struct {
	Owner   string        `json:"Owner"`
	KeyRing LegacyKeyRing `json:"Keypairs"`
}

func NewLegacyWallet(name string) *LegacyWallet {
	return &LegacyWallet{
		Owner:   name,
		KeyRing: NewLegacyKeyRing(),
	}
}

func (w *LegacyWallet) Version() uint32 {
	return 0
}

func (w *LegacyWallet) Name() string {
	return w.Owner
}

func (w *LegacyWallet) SetName(newName string) {
	w.Owner = newName
}

func (w *LegacyWallet) DescribePublicKey(pubKey string) (PublicKey, error) {
	keyPair, err := w.KeyRing.FindPair(pubKey)
	if err != nil {
		return nil, err
	}

	return keyPair.ToPublicKey(), nil
}

func (w *LegacyWallet) ListPublicKeys() []PublicKey {
	originalKeys := w.KeyRing.GetKeyPairs()
	keys := make([]PublicKey, len(originalKeys))
	for i, key := range originalKeys {
		keys[i] = key.ToPublicKey()
	}
	return keys
}

// ListKeyPairs lists the key pairs. Be careful, it contains the private key.
func (w *LegacyWallet) ListKeyPairs() []KeyPair {
	originalKeys := w.KeyRing.GetKeyPairs()
	keys := make([]KeyPair, len(originalKeys))
	for i, key := range originalKeys {
		keys[i] = key.DeepCopy()
	}
	return keys
}

func (w *LegacyWallet) GenerateKeyPair(meta []Meta) (KeyPair, error) {
	kp, err := GenKeyPair(crypto.Ed25519, 1)
	if err != nil {
		return nil, err
	}

	kp.MetaList = meta

	w.KeyRing.Upsert(*kp)

	deepCopy := kp.DeepCopy()
	return deepCopy, nil
}

func (w *LegacyWallet) TaintKey(pubKey string) error {
	keyPair, err := w.KeyRing.FindPair(pubKey)
	if err != nil {
		return err
	}

	if err := keyPair.Taint(); err != nil {
		return err
	}

	w.KeyRing.Upsert(keyPair)

	return nil
}

func (w *LegacyWallet) UntaintKey(pubKey string) error {
	keyPair, err := w.KeyRing.FindPair(pubKey)
	if err != nil {
		return err
	}

	keyPair.Tainted = false

	w.KeyRing.Upsert(keyPair)

	return nil
}

func (w *LegacyWallet) UpdateMeta(pubKey string, meta []Meta) error {
	keyPair, err := w.KeyRing.FindPair(pubKey)
	if err != nil {
		return err
	}

	keyPair.MetaList = meta

	w.KeyRing.Upsert(keyPair)
	return nil
}

func (w *LegacyWallet) SignAny(pubKey string, data []byte) ([]byte, error) {
	keyPair, err := w.KeyRing.FindPair(pubKey)
	if err != nil {
		return nil, err
	}

	if keyPair.Tainted {
		return nil, ErrPubKeyIsTainted
	}

	return keyPair.SignAny(data)
}

func (w *LegacyWallet) VerifyAny(pubKey string, data, sig []byte) (bool, error) {
	keyPair, err := w.KeyRing.FindPair(pubKey)
	if err != nil {
		return false, err
	}

	return keyPair.VerifyAny(data, sig)
}

func (w *LegacyWallet) SignTxV2(pubKey string, data []byte) (*commandspb.Signature, error) {
	keyPair, err := w.KeyRing.FindPair(pubKey)
	if err != nil {
		return nil, err
	}

	if keyPair.Tainted {
		return nil, ErrPubKeyIsTainted
	}

	return keyPair.Sign(data)
}
