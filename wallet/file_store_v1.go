package wallet

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"
)

const (
	walletBaseFolder = "wallets"
)

type FileV1 struct {
	rootPath string
}

func NewFileStoreV1(rootPath string) (*FileV1, error) {
	err := fsutil.EnsureDir(filepath.Join(rootPath, walletBaseFolder))
	if err != nil {
		return nil, err
	}

	return &FileV1{
		rootPath: rootPath,
	}, nil
}

func (c *FileV1) WalletExists(name string) bool {
	walletPath := c.walletPath(name)

	ok, _ := fsutil.PathExists(walletPath)
	return ok
}

func (c *FileV1) GetWallet(name, passphrase string) (Wallet, error) {
	walletPath := c.walletPath(name)

	if ok, _ := fsutil.PathExists(walletPath); !ok {
		return Wallet{}, ErrWalletDoesNotExists
	}

	buf, err := ioutil.ReadFile(walletPath)
	if err != nil {
		return Wallet{}, err
	}

	decBuf, err := crypto.Decrypt(buf, passphrase)
	if err != nil {
		return Wallet{}, err
	}

	w := &Wallet{}
	err = json.Unmarshal(decBuf, w)
	return *w, err
}

func (c *FileV1) SaveWallet(w Wallet, passphrase string) error {
	buf, err := json.Marshal(w)
	if err != nil {
		return err
	}

	encBuf, err := crypto.Encrypt(buf, passphrase)
	if err != nil {
		return err
	}

	f, err := os.Create(c.walletPath(w.Owner))
	if err != nil {
		return err
	}

	_, err = f.Write(encBuf)
	if err != nil {
		return err
	}

	return f.Close()
}

func (c *FileV1) GetWalletPath(name string) string {
	return c.walletPath(name)
}

func (c *FileV1) walletPath(name string) string {
	return filepath.Join(c.rootPath, walletBaseFolder, name)
}
