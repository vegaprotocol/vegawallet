package wallet

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/tyler-smith/go-bip39"
	"github.com/vegaprotocol/go-slip10"
)

const (
	// MaxEntropyByteSize is the entropy bytes size used for recovery phrase
	// generation.
	MaxEntropyByteSize = 256
	// MagicIndex is the registered HD wallet index for Vega's wallets.
	MagicIndex = 1789
	// OriginIndex is a constant index used to derive a node from the master
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
	id   string
}

// NewHDWallet creates a wallet with auto-generated recovery phrase. This is
// useful to create a brand-new wallet, without having to take care of the
// recovery phrase generation.
// The generated recovery phrase is returned alongside the created wallet.
func NewHDWallet(name string) (*HDWallet, string, error) {
	recoveryPhrase, err := NewRecoveryPhrase()
	if err != nil {
		return nil, "", err
	}

	w, err := ImportHDWallet(name, recoveryPhrase, LatestVersion)
	if err != nil {
		return nil, "", err
	}

	return w, recoveryPhrase, err
}

// ImportHDWallet creates a wallet based on the recovery phrase in input. This
// is useful import or retrieve a wallet.
func ImportHDWallet(name, recoveryPhrase string, version uint32) (*HDWallet, error) {
	if !bip39.IsMnemonicValid(recoveryPhrase) {
		return nil, ErrInvalidRecoveryPhrase
	}

	if !IsVersionSupported(version) {
		return nil, NewUnsupportedWalletVersionError(version)
	}

	walletNode, err := deriveWalletNodeFromRecoveryPhrase(recoveryPhrase)
	if err != nil {
		return nil, err
	}

	return &HDWallet{
		version: version,
		name:    name,
		keyRing: NewHDKeyRing(),
		node:    walletNode,
		id:      walletID(walletNode),
	}, nil
}

func (w *HDWallet) Version() uint32 {
	return w.version
}

func (w *HDWallet) Name() string {
	return w.name
}

func (w *HDWallet) ID() string {
	return w.id
}

func (w *HDWallet) Type() string {
	if w.IsIsolated() {
		return "HD wallet (isolated)"
	}
	return "HD wallet"
}

func (w *HDWallet) SetName(newName string) {
	w.name = newName
}

// DescribeKeyPair returns all the information associated with a public key.
func (w *HDWallet) DescribeKeyPair(pubKey string) (KeyPair, error) {
	keyPair, ok := w.keyRing.FindPair(pubKey)
	if !ok {
		return nil, ErrPubKeyDoesNotExist
	}
	return &keyPair, nil
}

// GetMasterKeyPair returns all the information associated to a master key pair.
func (w *HDWallet) GetMasterKeyPair() (MasterKeyPair, error) {
	if w.IsIsolated() {
		return nil, ErrIsolatedWalletDoesNotHaveMasterKey
	}

	pubKey, priKey := w.node.Keypair()
	keyPair, err := NewHDMasterKeyPair(pubKey, priKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get master key pair: %w", err)
	}

	return keyPair, nil
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
	if w.IsIsolated() {
		return nil, ErrIsolatedWalletCantGenerateKeyPairs
	}
	nextIndex := w.keyRing.NextIndex()

	keyNode, err := w.deriveKeyNode(nextIndex)
	if err != nil {
		return nil, err
	}

	publicKey, privateKey := keyNode.Keypair()
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

// UntaintKey remove the taint on a key.
func (w *HDWallet) UntaintKey(pubKey string) error {
	keyPair, ok := w.keyRing.FindPair(pubKey)
	if !ok {
		return ErrPubKeyDoesNotExist
	}

	if err := keyPair.Untaint(); err != nil {
		return err
	}

	w.keyRing.Upsert(keyPair)

	return nil
}

// UpdateMeta replaces the key's metadata by the new ones.
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

	return keyPair.SignAny(data)
}

func (w *HDWallet) VerifyAny(pubKey string, data, sig []byte) (bool, error) {
	keyPair, ok := w.keyRing.FindPair(pubKey)
	if !ok {
		return false, ErrPubKeyDoesNotExist
	}

	return keyPair.VerifyAny(data, sig)
}

func (w *HDWallet) SignTx(pubKey string, data []byte) (*Signature, error) {
	keyPair, ok := w.keyRing.FindPair(pubKey)
	if !ok {
		return nil, ErrPubKeyDoesNotExist
	}

	return keyPair.Sign(data)
}

func (w *HDWallet) IsolateWithKey(pubKey string) (Wallet, error) {
	keyPair, ok := w.keyRing.FindPair(pubKey)
	if !ok {
		return nil, ErrPubKeyDoesNotExist
	}

	if keyPair.IsTainted() {
		return nil, ErrPubKeyIsTainted
	}

	return &HDWallet{
		version: w.version,
		name:    fmt.Sprintf("%s.%s.isolated", w.name, keyPair.PublicKey()[0:8]),
		keyRing: LoadHDKeyRing([]HDKeyPair{keyPair}),
		id:      w.id,
	}, nil
}

func (w *HDWallet) IsIsolated() bool {
	return w.node == nil
}

type jsonHDWallet struct {
	// The wallet name is retrieved from the file name it is stored in, so no
	// need to serialize it.

	Version uint32       `json:"version"`
	Node    *slip10.Node `json:"node,omitempty"`
	ID      string       `json:"id,omitempty"`
	Keys    []HDKeyPair  `json:"keys"`
}

func (w *HDWallet) MarshalJSON() ([]byte, error) {
	jsonW := jsonHDWallet{
		Version: w.Version(),
		Keys:    w.keyRing.ListKeyPairs(),
		Node:    w.node,
		ID:      w.id,
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
		keyRing: LoadHDKeyRing(jsonW.Keys),
		node:    jsonW.Node,
		id:      jsonW.ID,
	}

	if len(w.id) == 0 {
		w.id = walletID(jsonW.Node)
	}

	return nil
}

func (w *HDWallet) deriveKeyNode(nextIndex uint32) (*slip10.Node, error) {
	var derivationFn func(uint32) (*slip10.Node, error)
	switch w.version {
	case Version1:
		derivationFn = w.deriveKeyNodeV1
	case Version2:
		derivationFn = w.deriveKeyNodeV2
	default:
		return nil, NewUnsupportedWalletVersionError(w.version)
	}

	return derivationFn(nextIndex)
}

func (w *HDWallet) deriveKeyNodeV1(nextIndex uint32) (*slip10.Node, error) {
	keyNode, err := w.node.Derive(OriginIndex + nextIndex)
	if err != nil {
		return nil, fmt.Errorf("couldn't derive key node for index %d: %w", OriginIndex+nextIndex, err)
	}
	return keyNode, nil
}

func (w *HDWallet) deriveKeyNodeV2(nextIndex uint32) (*slip10.Node, error) {
	defaultSubNode, err := w.node.Derive(slip10.FirstHardenedIndex)
	if err != nil {
		return nil, fmt.Errorf("couldn't derive default sub-node: %w", err)
	}
	keyNode, err := defaultSubNode.Derive(slip10.FirstHardenedIndex + nextIndex)
	if err != nil {
		return nil, fmt.Errorf("couldn't derive key node for index %d: %w", OriginIndex+nextIndex, err)
	}
	return keyNode, nil
}

// NewRecoveryPhrase generates a recovery phrase with an entropy of 256 bits.
func NewRecoveryPhrase() (string, error) {
	entropy, err := bip39.NewEntropy(MaxEntropyByteSize)
	if err != nil {
		return "", fmt.Errorf("couldn't create new wallet: %w", err)
	}
	recoveryPhrase, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("couldn't create recovery phrase: %w", err)
	}
	return recoveryPhrase, nil
}

func deriveWalletNodeFromRecoveryPhrase(recoveryPhrase string) (*slip10.Node, error) {
	seed := bip39.NewSeed(recoveryPhrase, "")
	masterNode, err := slip10.NewMasterNode(seed)
	if err != nil {
		return nil, fmt.Errorf("couldn't create master node: %w", err)
	}
	walletNode, err := masterNode.Derive(OriginIndex)
	if err != nil {
		return nil, fmt.Errorf("couldn't derive wallet node: %w", err)
	}
	return walletNode, nil
}

func walletID(walletNode *slip10.Node) string {
	pubKey, _ := walletNode.Keypair()
	return hex.EncodeToString(pubKey)
}
