package wallet

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"
)

type Wallet struct {
	Owner    string
	Keypairs KeyRing
}

func NewWallet(name string) *Wallet {
	return &Wallet{
		Owner: name,
	}
}

type Meta struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func New(owner string) Wallet {
	return Wallet{
		Owner:    owner,
		Keypairs: NewKeyRing(),
	}
}

func EnsureBaseFolder(root string) error {
	return fsutil.EnsureDir(filepath.Join(root, walletBaseFolder))
}

func WalletFileExists(root, file string) bool {
	if ok, _ := fsutil.PathExists(filepath.Join(root, walletBaseFolder, file)); ok {
		return true
	}
	return false
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
