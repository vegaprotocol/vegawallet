package v1_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"code.vegaprotocol.io/go-wallet/crypto"
)

type configDir struct {
	rootPath string
}

func newConfigDir() configDir {
	rootPath := filepath.Join("tmp", "vegatests", "wallet", crypto.RandomStr(10))

	return configDir{
		rootPath: rootPath,
	}
}

func (d configDir) RootPath() string {
	return d.rootPath
}

func (d configDir) RSAKeysPath() string {
	return filepath.Join(d.rootPath, "wallet_rsa")
}

func (d configDir) WalletsPath() string {
	return filepath.Join(d.rootPath, "wallets")
}

func (d configDir) WalletPath(name string) string {
	return filepath.Join(d.rootPath, "wallets", name)
}

func (d configDir) WalletContent(name string) string {
	buf, err := ioutil.ReadFile(d.WalletPath(name))
	if err != nil {
		panic(err)
	}
	return string(buf)
}

func (d configDir) Remove() {
	err := os.RemoveAll(d.rootPath)
	if err != nil {
		panic(err)
	}
}

