package wallet

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"
)

var (
	ErrWalletAlreadyExists = errors.New("a wallet with the same name already exists")
	ErrWalletDoesNotExists = errors.New("wallet does not exist")
)

const (
	walletBaseFolder = "wallets"
)

type Wallet struct {
	Owner    string
	Keypairs []Keypair
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

type Meta struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func New(owner string) Wallet {
	return Wallet{
		Owner:    owner,
		Keypairs: []Keypair{},
	}
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

func NewKeypair(algo crypto.SignatureAlgorithm, pub, priv []byte) Keypair {
	return Keypair{
		Algorithm: algo,
		pubBytes:  pub,
		privBytes: priv,
	}
}

func EnsureBaseFolder(root string) error {
	return fsutil.EnsureDir(filepath.Join(root, walletBaseFolder))
}

func CreateWalletFile(walletpath, owner, passphrase string) (*Wallet, error) {
	w := Wallet{
		Owner: owner,
	}
	// make sure this do not exists already
	if ok, _ := fsutil.PathExists(walletpath); ok {
		return nil, ErrWalletAlreadyExists
	}

	return WriteWalletFile(&w, walletpath, passphrase)
}

func Create(root, owner, passphrase string) (*Wallet, error) {
	// build walletpath
	walletpath := filepath.Join(root, walletBaseFolder, owner)
	return CreateWalletFile(walletpath, owner, passphrase)
}

// WalletPath get the path to the wallet file, check if actually is a file
func WalletPath(root, owner string) (string, error) {
	path := filepath.Join(root, walletBaseFolder, owner)
	if ok, err := fsutil.FileExists(path); !ok || err != nil {
		return "", ErrWalletDoesNotExists
	}
	return path, nil
}

func AddKeypair(kp *Keypair, root, owner, passphrase string) (*Wallet, error) {
	w, err := Read(root, owner, passphrase)
	if err != nil {
		return nil, err
	}

	w.Keypairs = append(w.Keypairs, *kp)

	return writeWallet(w, root, owner, passphrase)
}

func ReadWalletFile(walletpath, passphrase string) (*Wallet, error) {
	// make sure this do not exists already
	if ok, _ := fsutil.PathExists(walletpath); !ok {
		return nil, ErrWalletDoesNotExists
	}

	// read file
	buf, err := ioutil.ReadFile(walletpath)
	if err != nil {
		return nil, err
	}

	// decrypt the buffer
	decBuf, err := crypto.Decrypt(buf, passphrase)
	if err != nil {
		return nil, err
	}

	// unmarshal the wallet now an return
	w := &Wallet{}
	return w, json.Unmarshal(decBuf, w)
}

func Read(root, owner, passphrase string) (*Wallet, error) {
	// build walletpath
	walletpath := filepath.Join(root, walletBaseFolder, owner)

	return ReadWalletFile(walletpath, passphrase)
}

func Write(w *Wallet, root, owner, passphrase string) (*Wallet, error) {
	return writeWallet(w, root, owner, passphrase)
}

func WriteWalletFile(w *Wallet, walletpath, passphrase string) (*Wallet, error) {
	// marshal our wallet
	buf, err := json.Marshal(w)
	if err != nil {
		return nil, err
	}

	// encrypt our data
	encBuf, err := crypto.Encrypt(buf, passphrase)
	if err != nil {
		return nil, err
	}

	// create and write file
	f, err := os.Create(walletpath)
	if err != nil {
		return nil, err
	}
	f.Write(encBuf)
	f.Close()

	return w, nil
}

func writeWallet(w *Wallet, root, owner, passphrase string) (*Wallet, error) {
	// build walletpath
	walletpath := filepath.Join(root, walletBaseFolder, owner)
	return WriteWalletFile(w, walletpath, passphrase)
}

func Sign(alg crypto.SignatureAlgorithm, kp *Keypair, message []byte) ([]byte, error) {
	return alg.Sign(kp.privBytes, message)
}

func Verify(alg crypto.SignatureAlgorithm, kp *Keypair, message []byte, sig []byte) (bool, error) {
	return alg.Verify(kp.pubBytes, message, sig)
}
