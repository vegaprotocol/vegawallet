package wallet

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"math/big"
	"sync"

	types "code.vegaprotocol.io/go-wallet/proto"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
)

var (
	ErrPubKeyDoesNotExist   = errors.New("public key does not exist")
	ErrPubKeyAlreadyTainted = errors.New("public key is already tainted")
	ErrPubKeyIsTainted      = errors.New("public key is tainted")
)

// Auth ...
//go:generate go run github.com/golang/mock/mockgen -destination mocks/auth_mock.go -package mocks code.vegaprotocol.io/go-wallet/wallet Auth
type Auth interface {
	NewSession(walletname string) (string, error)
	VerifyToken(token string) (string, error)
	Revoke(token string) error
}

type Handler struct {
	log      *zap.Logger
	auth     Auth
	rootPath string

	// wallet name -> wallet
	store map[string]Wallet

	// just to make sure we do not access same file conccurently or the map
	mu sync.RWMutex
}

func NewHandler(log *zap.Logger, auth Auth, rootPath string) *Handler {
	return &Handler{
		log:      log,
		auth:     auth,
		rootPath: rootPath,
		store:    map[string]Wallet{},
	}
}

// CreateWallet return the actual token
func (h *Handler) CreateWallet(wallet, passphrase string) (string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.store[wallet]; ok {
		return "", ErrWalletAlreadyExists
	}

	w, err := Create(h.rootPath, wallet, passphrase)
	if err != nil {
		return "", err
	}

	h.store[wallet] = *w
	return h.auth.NewSession(wallet)
}

func (h *Handler) LoginWallet(wallet, passphrase string) (string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	// first check if the user own the wallet
	w, err := Read(h.rootPath, wallet, passphrase)
	if err != nil {
		return "", ErrWalletDoesNotExists
	}

	// then store it in the memory store then
	if _, ok := h.store[wallet]; !ok {
		h.store[wallet] = *w
	}

	return h.auth.NewSession(wallet)
}

func (h *Handler) WalletPath(token string) (string, error) {
	wallet, err := h.auth.VerifyToken(token)
	if err != nil {
		return "", err
	}
	return WalletPath(h.rootPath, wallet)
}

func (h *Handler) RevokeToken(token string) error {
	return h.auth.Revoke(token)
}

func (h *Handler) GenerateKeypair(token, passphrase string) (string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	wname, err := h.auth.VerifyToken(token)
	if err != nil {
		return "", err
	}

	// validate passphrase
	_, err = Read(h.rootPath, wname, passphrase)
	if err != nil {
		return "", err
	}

	w, ok := h.store[wname]
	if !ok {
		// this should never happen as we cannot have a valid session
		// without the actual wallet being loaded in memory but...
		return "", ErrWalletDoesNotExists
	}

	kp, err := GenKeypair(crypto.Ed25519)
	if err != nil {
		return "", err
	}

	w.Keypairs = append(w.Keypairs, *kp)
	_, err = writeWallet(&w, h.rootPath, wname, passphrase)
	if err != nil {
		return "", err
	}

	h.store[wname] = w
	return kp.Pub, nil
}

func (h *Handler) GetPublicKey(token, pubKey string) (*Keypair, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	wname, err := h.auth.VerifyToken(token)
	if err != nil {
		return nil, err
	}

	w, ok := h.store[wname]
	if !ok {
		// this should never happen as we cannot have a valid session
		// without the actual wallet being loaded in memory but...
		return nil, ErrWalletDoesNotExists
	}

	// copy the key so we do not propagate private keys
	for _, v := range w.Keypairs {
		if v.Pub != pubKey {
			continue
		}
		out := v
		out.Priv = ""
		out.privBytes = []byte{}
		return &out, nil
	}

	return nil, ErrPubKeyDoesNotExist
}

func (h *Handler) ListPublicKeys(token string) ([]Keypair, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	wname, err := h.auth.VerifyToken(token)
	if err != nil {
		return nil, err
	}

	w, ok := h.store[wname]
	if !ok {
		// this should never happen as we cannot have a valid session
		// without the actual wallet being loaded in memory but...
		return nil, ErrWalletDoesNotExists
	}

	// copy all keys so we do not propagate private keys
	out := make([]Keypair, 0, len(w.Keypairs))
	for _, v := range w.Keypairs {
		kp := v
		kp.Priv = ""
		kp.privBytes = []byte{}
		out = append(out, kp)
	}

	return out, nil
}

func (h *Handler) SignTx(token, tx, pubkey string) (SignedBundle, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// first the transaction would be in base64, let's decode
	rawtx, err := base64.StdEncoding.DecodeString(tx)
	if err != nil {
		return SignedBundle{}, err
	}

	// then get the wallet name out of the token
	wname, err := h.auth.VerifyToken(token)
	if err != nil {
		return SignedBundle{}, err
	}

	w, ok := h.store[wname]
	if !ok {
		// this should never happen as we cannot have a valid session
		// without the actual wallet being loaded in memory but...
		return SignedBundle{}, ErrWalletDoesNotExists
	}

	// let's retrieve the private key from the public key
	var kp *Keypair
	for i := range w.Keypairs {
		if w.Keypairs[i].Pub == pubkey {
			kp = &w.Keypairs[i]
			break
		}
	}
	// we did not find this pub key
	if kp == nil {
		return SignedBundle{}, ErrPubKeyDoesNotExist
	}

	if kp.Tainted {
		return SignedBundle{}, ErrPubKeyIsTainted
	}

	// we build the transaction payload
	txTy := &types.Transaction{
		InputData: rawtx,
		Nonce:     makeNonce(),
		From: &types.Transaction_PubKey{
			PubKey: kp.pubBytes,
		},
	}

	rawTxTy, err := proto.Marshal(txTy)
	if err != nil {
		return SignedBundle{}, err
	}

	// then lets sign the stuff and return it
	sig, err := kp.Algorithm.Sign(kp.privBytes, rawTxTy)
	if err != nil {
		return SignedBundle{}, err
	}

	return SignedBundle{
		Tx: rawTxTy,
		Sig: Signature{
			Sig:     sig,
			Algo:    kp.Algorithm.Name(),
			Version: kp.Algorithm.Version(),
		},
	}, nil
}

func (h *Handler) TaintKey(token, pubkey, passphrase string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	wname, err := h.auth.VerifyToken(token)
	if err != nil {
		return err
	}

	_, err = Read(h.rootPath, wname, passphrase)
	if err != nil {
		return err
	}

	w, ok := h.store[wname]
	if !ok {
		// this should never happen as we cannot have a valid session
		// without the actual wallet being loaded in memory but...
		return ErrWalletDoesNotExists
	}

	// let's retrieve the private key from the public key
	var kp *Keypair
	for i := range w.Keypairs {
		if w.Keypairs[i].Pub == pubkey {
			kp = &w.Keypairs[i]
			break
		}
	}
	// we did not find this pub key
	if kp == nil {
		return ErrPubKeyDoesNotExist
	}

	if kp.Tainted {
		return ErrPubKeyAlreadyTainted
	}

	kp.Tainted = true

	_, err = writeWallet(&w, h.rootPath, wname, passphrase)
	if err != nil {
		return err
	}

	h.store[wname] = w
	return nil
}

func (h *Handler) UpdateMeta(token, pubkey, passphrase string, meta []Meta) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	wname, err := h.auth.VerifyToken(token)
	if err != nil {
		return err
	}

	_, err = Read(h.rootPath, wname, passphrase)
	if err != nil {
		return err
	}

	w, ok := h.store[wname]
	if !ok {
		// this should never happen as we cannot have a valid session
		// without the actual wallet being loaded in memory but...
		return ErrWalletDoesNotExists
	}

	// let's retrieve the private key from the public key
	var kp *Keypair
	for i := range w.Keypairs {
		if w.Keypairs[i].Pub == pubkey {
			kp = &w.Keypairs[i]
			break
		}
	}
	// we did not find this pub key
	if kp == nil {
		return ErrPubKeyDoesNotExist
	}

	kp.Meta = meta

	_, err = writeWallet(&w, h.rootPath, wname, passphrase)
	if err != nil {
		return err
	}

	h.store[wname] = w
	return nil
}

func makeNonce() uint64 {
	max := &big.Int{}
	// set it to the max value of the uint64
	max.SetUint64(^uint64(0))
	nonce, _ := rand.Int(rand.Reader, max)
	return nonce.Uint64()
}
