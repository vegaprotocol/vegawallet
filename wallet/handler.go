package wallet

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"code.vegaprotocol.io/protos/commands"
	commandspb "code.vegaprotocol.io/protos/vega/commands/v1"
	vgcrypto "code.vegaprotocol.io/shared/libs/crypto"
	"github.com/golang/protobuf/proto"
)

//go:generate go run github.com/golang/mock/mockgen -destination mocks/store_mock.go -package mocks code.vegaprotocol.io/vegawallet/wallet Store
type Store interface {
	WalletExists(name string) bool
	SaveWallet(w Wallet, passphrase string) error
	GetWallet(name, passphrase string) (Wallet, error)
	GetWalletPath(name string) string
	ListWallets() ([]string, error)
}

type GenerateKeyRequest struct {
	Wallet     string
	Metadata   []Meta
	Passphrase string
}

type GenerateKeyResponse struct {
	Wallet struct {
		Name     string `json:"name"`
		FilePath string `json:"filePath"`
		Mnemonic string `json:"mnemonic,omitempty"`
	} `json:"wallet"`
	Key struct {
		KeyPair struct {
			PrivateKey string `json:"privateKey"`
			PublicKey  string `json:"publicKey"`
		} `json:"keyPair"`
		Algorithm struct {
			Name    string `json:"name"`
			Version uint32 `json:"version"`
		} `json:"algorithm"`
		Meta []Meta `json:"meta"`
	} `json:"key"`
}

func GenerateKey(store Store, req *GenerateKeyRequest) (*GenerateKeyResponse, error) {
	resp := &GenerateKeyResponse{}

	walletExists := store.WalletExists(req.Wallet)

	var wal Wallet
	if !walletExists {
		w, mnemonic, err := NewHDWallet(req.Wallet)
		if err != nil {
			return nil, fmt.Errorf("couldn't create HD wallet: %w", err)
		}
		wal = w

		resp.Wallet.Mnemonic = mnemonic
	} else {
		w, err := store.GetWallet(req.Wallet, req.Passphrase)
		if err != nil {
			if errors.Is(err, ErrWrongPassphrase) {
				return nil, err
			}
			return nil, fmt.Errorf("couldn't get wallet %s: %w", req.Wallet, err)
		}
		wal = w
	}

	req.Metadata = addDefaultKeyName(wal, req.Metadata)

	kp, err := wal.GenerateKeyPair(req.Metadata)
	if err != nil {
		return nil, err
	}

	if err := store.SaveWallet(wal, req.Passphrase); err != nil {
		return nil, fmt.Errorf("couldn't save wallet: %w", err)
	}

	resp.Wallet.Name = req.Wallet
	resp.Wallet.FilePath = store.GetWalletPath(req.Wallet)
	resp.Key.KeyPair.PublicKey = kp.PublicKey()
	resp.Key.KeyPair.PrivateKey = kp.PrivateKey()
	resp.Key.Algorithm.Name = kp.AlgorithmName()
	resp.Key.Algorithm.Version = kp.AlgorithmVersion()
	resp.Key.Meta = kp.Meta()

	return resp, nil
}

type AnnotateKeyRequest struct {
	Wallet     string
	PubKey     string
	Metadata   []Meta
	Passphrase string
}

func AnnotateKey(store Store, req *AnnotateKeyRequest) error {
	w, err := getWallet(store, req.Wallet, req.Passphrase)
	if err != nil {
		return err
	}

	if err = w.UpdateMeta(req.PubKey, req.Metadata); err != nil {
		return fmt.Errorf("couldn't update metadata: %w", err)
	}

	if err := store.SaveWallet(w, req.Passphrase); err != nil {
		return fmt.Errorf("couldn't save wallet: %w", err)
	}

	return nil
}

type TaintKeyRequest struct {
	Wallet     string
	PubKey     string
	Passphrase string
}

func TaintKey(store Store, req *TaintKeyRequest) error {
	w, err := getWallet(store, req.Wallet, req.Passphrase)
	if err != nil {
		return err
	}

	if err = w.TaintKey(req.PubKey); err != nil {
		return fmt.Errorf("couldn't taint key: %w", err)
	}

	if err := store.SaveWallet(w, req.Passphrase); err != nil {
		return fmt.Errorf("couldn't save wallet: %w", err)
	}

	return nil
}

type UntaintKeyRequest struct {
	Wallet     string
	PubKey     string
	Passphrase string
}

func UntaintKey(store Store, req *UntaintKeyRequest) error {
	w, err := getWallet(store, req.Wallet, req.Passphrase)
	if err != nil {
		return err
	}

	if err = w.UntaintKey(req.PubKey); err != nil {
		return fmt.Errorf("couldn't untaint key: %w", err)
	}

	if err := store.SaveWallet(w, req.Passphrase); err != nil {
		return fmt.Errorf("couldn't save wallet: %w", err)
	}

	return nil
}

type IsolateKeyRequest struct {
	Wallet     string
	PubKey     string
	Passphrase string
}

type IsolateKeyResponse struct {
	Wallet   string `json:"wallet"`
	FilePath string `json:"filePath"`
}

func IsolateKey(store Store, req *IsolateKeyRequest) (*IsolateKeyResponse, error) {
	w, err := getWallet(store, req.Wallet, req.Passphrase)
	if err != nil {
		return nil, err
	}

	isolatedWallet, err := w.IsolateWithKey(req.PubKey)
	if err != nil {
		return nil, fmt.Errorf("couldn't isolate wallet %s: %w", req.Wallet, err)
	}

	if err := store.SaveWallet(isolatedWallet, req.Passphrase); err != nil {
		return nil, fmt.Errorf("couldn't save isolated wallet %s: %w", isolatedWallet.Name(), err)
	}

	return &IsolateKeyResponse{
		Wallet:   isolatedWallet.Name(),
		FilePath: store.GetWalletPath(isolatedWallet.Name()),
	}, nil
}

type ListKeysRequest struct {
	Wallet     string
	Passphrase string
}

type ListKeysResponse struct {
	Keys []NamedPubKey `json:"keys"`
}

type DescribeKeyRequest struct {
	Wallet     string
	Passphrase string
	PubKey     string
}

type DescribeKeyResponse struct {
	PublicKey string `json:"publicKey"`

	Algorithm struct {
		Name    string `json:"name"`
		Version uint32 `json:"version"`
	} `json:"algorithm"`
	Meta      []Meta `json:"meta"`
	IsTainted bool   `json:"isTainted"`
}

type NamedPubKey struct {
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
}

func ListKeys(store Store, req *ListKeysRequest) (*ListKeysResponse, error) {
	w, err := getWallet(store, req.Wallet, req.Passphrase)
	if err != nil {
		return nil, err
	}

	kps := w.ListKeyPairs()
	keys := make([]NamedPubKey, 0, len(kps))
	for _, kp := range kps {
		keys = append(keys, NamedPubKey{
			Name:      GetKeyName(kp.Meta()),
			PublicKey: kp.PublicKey(),
		})
	}

	return &ListKeysResponse{
		Keys: keys,
	}, nil
}

func DescribeKey(store Store, req *DescribeKeyRequest) (*DescribeKeyResponse, error) {
	w, err := getWallet(store, req.Wallet, req.Passphrase)
	if err != nil {
		return nil, err
	}

	resp := &DescribeKeyResponse{}

	kp, err := w.DescribeKeyPair(req.PubKey)
	if err != nil {
		return nil, err
	}
	resp.PublicKey = kp.PublicKey()
	resp.Algorithm.Name = kp.AlgorithmName()
	resp.Algorithm.Version = kp.AlgorithmVersion()
	resp.Meta = kp.Meta()
	resp.IsTainted = kp.IsTainted()
	return resp, nil
}

type RotateKeyRequest struct {
	Wallet            string
	Passphrase        string
	NewPublicKey      string
	TXBlockHeight     uint32
	TargetBlockHeight uint32
}

type RotateKeyResponse struct {
	MasterPublicKey   string `json:"master_public_key"`
	NewPublicKey      string `json:"new_public_key"`
	Base64Transaction string `json:"base64_transaction"`
}

func RotateKey(store Store, req *RotateKeyRequest) (*RotateKeyResponse, error) {
	w, err := getWallet(store, req.Wallet, req.Passphrase)
	if err != nil {
		return nil, err
	}

	mKeyPair, err := w.GetMasterKeyPair()
	if err != nil {
		return nil, err
	}

	newKeyPair, err := w.GetKeyPair(req.NewPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get new key pair: %w", err)
	}

	inputData := commands.NewInputData(uint64(req.TXBlockHeight))
	inputData.Command = &commandspb.InputData_KeyRotateSubmission{
		KeyRotateSubmission: &commandspb.KeyRotateSubmission{
			KeyNumber:     uint64(newKeyPair.Index()),
			TargetBlock:   uint64(req.TargetBlockHeight),
			Time:          0, // @TODO fill this
			NewPubKeyHash: hex.EncodeToString(vgcrypto.Hash([]byte(req.NewPublicKey))),
		},
	}

	data, err := proto.Marshal(inputData)
	if err != nil {
		return nil, fmt.Errorf("failed marshal key rotate submission input data: %w", err)
	}

	sign, err := mKeyPair.Sign(data)
	if err != nil {
		return nil, fmt.Errorf("failed to sign key rotate submission input data: %w", err)
	}

	protoSignature := &commandspb.Signature{
		Value:   sign.Value,
		Algo:    sign.Algo,
		Version: sign.Version,
	}

	transaction := commands.NewTransaction(mKeyPair.PublicKey(), data, protoSignature)
	transactionRaw, err := proto.Marshal(transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %w", err)
	}

	return &RotateKeyResponse{
		MasterPublicKey:   mKeyPair.PublicKey(),
		NewPublicKey:      req.NewPublicKey,
		Base64Transaction: base64.StdEncoding.EncodeToString(transactionRaw),
	}, nil
}

type GetWalletInfoRequest struct {
	Wallet     string
	Passphrase string
}

type GetWalletInfoResponse struct {
	Type    string `json:"type"`
	Version uint32 `json:"version"`
	ID      string `json:"id"`
}

func GetWalletInfo(store Store, req *GetWalletInfoRequest) (*GetWalletInfoResponse, error) {
	w, err := getWallet(store, req.Wallet, req.Passphrase)
	if err != nil {
		return nil, err
	}

	return &GetWalletInfoResponse{
		Type:    w.Type(),
		Version: w.Version(),
		ID:      w.ID(),
	}, nil
}

type ImportWalletRequest struct {
	Wallet     string
	Mnemonic   string
	Version    uint32
	Passphrase string
}

type ImportWalletResponse struct {
	Name     string `json:"name"`
	FilePath string `json:"filePath"`
}

func ImportWallet(store Store, req *ImportWalletRequest) (*ImportWalletResponse, error) {
	if store.WalletExists(req.Wallet) {
		return nil, ErrWalletAlreadyExists
	}

	w, err := ImportHDWallet(req.Wallet, req.Mnemonic, req.Version)
	if err != nil {
		return nil, fmt.Errorf("couldn't import the wallet: %w", err)
	}

	if err := store.SaveWallet(w, req.Passphrase); err != nil {
		return nil, fmt.Errorf("couldn't save wallet %s: %w", w.Name(), err)
	}

	return &ImportWalletResponse{
		Name:     w.Name(),
		FilePath: store.GetWalletPath(w.Name()),
	}, nil
}

type ListWalletsResponse struct {
	Wallets []string `json:"wallets"`
}

func ListWallets(store Store) (*ListWalletsResponse, error) {
	ws, err := store.ListWallets()
	if err != nil {
		return nil, err
	}

	resp := &ListWalletsResponse{}
	resp.Wallets = make([]string, 0, len(ws))
	resp.Wallets = append(resp.Wallets, ws...)

	return resp, nil
}

type SignMessageRequest struct {
	Wallet     string
	PubKey     string
	Message    []byte
	Passphrase string
}

type SignMessageResponse struct {
	Base64 string `json:"hexSignature"`
	Bytes  []byte `json:"bytesSignature"`
}

func SignMessage(store Store, req *SignMessageRequest) (*SignMessageResponse, error) {
	w, err := getWallet(store, req.Wallet, req.Passphrase)
	if err != nil {
		return nil, err
	}

	sig, err := w.SignAny(req.PubKey, req.Message)
	if err != nil {
		return nil, fmt.Errorf("couldn't sign message: %w", err)
	}

	return &SignMessageResponse{
		Base64: base64.StdEncoding.EncodeToString(sig),
		Bytes:  sig,
	}, nil
}

func getWallet(store Store, wallet, passphrase string) (Wallet, error) {
	if !store.WalletExists(wallet) {
		return nil, ErrWalletDoesNotExists
	}

	w, err := store.GetWallet(wallet, passphrase)
	if err != nil {
		if errors.Is(err, ErrWrongPassphrase) {
			return nil, err
		}
		return nil, fmt.Errorf("couldn't get wallet %s: %w", wallet, err)
	}

	return w, nil
}

func addDefaultKeyName(w Wallet, meta []Meta) []Meta {
	for _, m := range meta {
		if m.Key == KeyNameMeta {
			return meta
		}
	}

	if len(meta) == 0 {
		meta = []Meta{}
	}

	nextID := len(w.ListKeyPairs()) + 1

	meta = append(meta, Meta{
		Key:   KeyNameMeta,
		Value: fmt.Sprintf("%s key %d", w.Name(), nextID),
	})
	return meta
}
