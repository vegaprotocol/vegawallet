package wallet

import (
	"encoding/json"
	"fmt"

	"code.vegaprotocol.io/go-wallet/crypto"
	typespb "code.vegaprotocol.io/go-wallet/internal/proto"
	commandspb "code.vegaprotocol.io/go-wallet/internal/proto/commands/v1"

	"github.com/golang/protobuf/proto"
	"github.com/tyler-smith/go-bip39"
	"github.com/vegaprotocol/go-slip10"
)

const (
	MagicIndex = 1789
	// OriginIndex is a constant index used to derived a node from the master
	// node. The resulting node will be used to generate the cryptographic keys.
	OriginIndex = slip10.FirstHardenedIndex + MagicIndex
)

type HDWallet struct {
	version uint32
	name    string
	keyRing *HDKeyRing

	// node is the node from which the cryptographic keys are generated. This is
	// not the master node. This is a node derived from the master. Its
	// derivation index is constant (see OriginIndex). This node is referred as
	// "wallet node".
	node *slip10.Node
}

// NewHDWallet creates a wallet with auto-generated mnemonic. This is useful to
// create a brand new wallet, without having to take care of the mnemonic
// generation.
// The generated mnemonic is returned alongside the created wallet.
func NewHDWallet(name string) (*HDWallet, string, error) {
	mnemonic, err := NewMnemonic()
	if err != nil {
		return nil, "", err
	}

	w, err := ImportHDWallet(name, mnemonic)
	if err != nil {
		return nil, "", err
	}

	return w, mnemonic, err
}

// ImportHDWallet creates a wallet based on the mnemonic in input. This is
// useful import or retrieve a wallet.
func ImportHDWallet(name, mnemonic string) (*HDWallet, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, ErrInvalidMnemonic
	}

	walletNode, err := deriveWalletNodeFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	return &HDWallet{
		version: 1,
		name:    name,
		keyRing: NewHDKeyRing(),
		node:    walletNode,
	}, nil
}

func (w *HDWallet) Version() uint32 {
	return w.version
}

func (w *HDWallet) Name() string {
	return w.name
}

// DescribePublicKey returns all the information associated to a public key,
// except the private key.
func (w *HDWallet) DescribePublicKey(pubKey string) (PublicKey, error) {
	keyPair, ok := w.keyRing.FindPair(pubKey)
	if !ok {
		return nil, ErrPubKeyDoesNotExist
	}

	publicKey := keyPair.ToPublicKey()
	return &publicKey, nil
}

// ListPublicKeys lists the public keys with their information. The private keys
// are not returned.
func (w *HDWallet) ListPublicKeys() []PublicKey {
	originalKeys := w.keyRing.ListKeyPairs()
	keys := make([]PublicKey, len(originalKeys))
	for i, key := range originalKeys {
		publicKey := key.ToPublicKey()
		keys[i] = &publicKey
	}
	return keys
}

// ListKeyPairs lists the key pairs. Be careful, it contains the private key.
func (w *HDWallet) ListKeyPairs() []KeyPair {
	originalKeys := w.keyRing.ListKeyPairs()
	keys := make([]KeyPair, len(originalKeys))
	for i, key := range originalKeys {
		keys[i] = key.DeepCopy()
	}
	return keys
}

// GenerateKeyPair generates a new key pair from a node, that is derived from
// the wallet node.
func (w *HDWallet) GenerateKeyPair(meta []Meta) (KeyPair, error) {
	nextIndex := w.keyRing.NextIndex()
	childNode, err := w.node.Derive(OriginIndex + nextIndex)
	if err != nil {
		return nil, err
	}

	publicKey, privateKey := childNode.Keypair()
	keyPair, err := NewHDKeyPair(nextIndex, publicKey, privateKey)
	if err != nil {
		return nil, err
	}

	keyPair.meta = meta

	w.keyRing.Upsert(*keyPair)

	return keyPair.DeepCopy(), nil
}

// TaintKey marks a key as tainted.
func (w *HDWallet) TaintKey(pubKey string) error {
	keyPair, ok := w.keyRing.FindPair(pubKey)
	if !ok {
		return ErrPubKeyDoesNotExist
	}

	if err := keyPair.Taint(); err != nil {
		return err
	}

	w.keyRing.Upsert(keyPair)

	return nil
}

// UpdateMeta replaces the key's meta by the new ones.
func (w *HDWallet) UpdateMeta(pubKey string, meta []Meta) error {
	keyPair, ok := w.keyRing.FindPair(pubKey)
	if !ok {
		return ErrPubKeyDoesNotExist
	}

	keyPair.meta = meta

	w.keyRing.Upsert(keyPair)
	return nil
}

func (w *HDWallet) SignAny(pubKey string, data []byte) ([]byte, error) {
	keyPair, ok := w.keyRing.FindPair(pubKey)
	if !ok {
		return nil, ErrPubKeyDoesNotExist
	}

	if keyPair.IsTainted() {
		return nil, ErrPubKeyIsTainted
	}

	return keyPair.SignAny(data)
}

func (w *HDWallet) VerifyAny(pubKey string, data, sig []byte) (bool, error) {
	keyPair, ok := w.keyRing.FindPair(pubKey)
	if !ok {
		return false, ErrPubKeyDoesNotExist
	}

	return keyPair.VerifyAny(data, sig)
}

func (w *HDWallet) SignTxV1(pubKey string, data []byte, blockHeight uint64) (SignedBundle, error) {
	keyPair, ok := w.keyRing.FindPair(pubKey)
	if !ok {
		return SignedBundle{}, ErrPubKeyDoesNotExist
	}

	if keyPair.IsTainted() {
		return SignedBundle{}, ErrPubKeyIsTainted
	}

	txTy := &typespb.Transaction{
		InputData:   data,
		Nonce:       crypto.NewNonce(),
		BlockHeight: blockHeight,
		From: &typespb.Transaction_PubKey{
			PubKey: keyPair.publicKey.bytes,
		},
	}

	rawTxTy, err := proto.Marshal(txTy)
	if err != nil {
		return SignedBundle{}, err
	}

	sig, err := keyPair.SignAny(rawTxTy)
	if err != nil {
		return SignedBundle{}, err
	}

	return SignedBundle{
		Tx: rawTxTy,
		Sig: Signature{
			Sig:     sig,
			Algo:    keyPair.AlgorithmName(),
			Version: keyPair.AlgorithmVersion(),
		},
	}, nil
}

func (w *HDWallet) SignTxV2(pubKey string, data []byte) (*commandspb.Signature, error) {
	keyPair, ok := w.keyRing.FindPair(pubKey)
	if !ok {
		return nil, ErrPubKeyDoesNotExist
	}

	if keyPair.IsTainted() {
		return nil, ErrPubKeyIsTainted
	}

	return keyPair.Sign(data)
}

type jsonHDWallet struct {
	Version uint32       `json:"version"`
	Name    string       `json:"name"`
	Node    *slip10.Node `json:"node"`
	Keys    []HDKeyPair  `json:"keys"`
}

func (w *HDWallet) MarshalJSON() ([]byte, error) {
	jsonW := jsonHDWallet{
		Version: w.Version(),
		Name:    w.Name(),
		Keys:    w.keyRing.ListKeyPairs(),
		Node:    w.node,
	}
	return json.Marshal(jsonW)
}

func (w *HDWallet) UnmarshalJSON(data []byte) error {
	jsonW := &jsonHDWallet{}
	if err := json.Unmarshal(data, jsonW); err != nil {
		return err
	}

	*w = HDWallet{
		version: jsonW.Version,
		name:    jsonW.Name,
		keyRing: LoadHDKeyRing(jsonW.Keys),
		node:    jsonW.Node,
	}

	return nil
}

// NewMnemonic generates a mnemonic with an entropy of 256 bits.
func NewMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", fmt.Errorf("couldn't create new wallet: %v", err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("couldn't create new wallet: %v", err)
	}
	return mnemonic, nil
}

func deriveWalletNodeFromMnemonic(mnemonic string) (*slip10.Node, error) {
	seed := bip39.NewSeed(mnemonic, "")
	masterNode, err := slip10.NewMasterNode(seed)
	if err != nil {
		return nil, fmt.Errorf("couldn't create master node: %v", err)
	}
	walletNode, err := masterNode.Derive(OriginIndex)
	if err != nil {
		return nil, fmt.Errorf("couldn't derive wallet node: %v", err)
	}
	return walletNode, nil
}
