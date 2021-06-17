package wallet_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"code.vegaprotocol.io/go-wallet/crypto"
)

func rootDir() string {
	path := filepath.Join(rootDirPath, crypto.RandomStr(10))
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return path
}

type configDir struct {
	path string
}

func newConfigDir() configDir {
	path := filepath.Join("/tmp/vegatests/wallet/", crypto.RandomStr(10))

	return configDir{
		path: path,
	}
}

func (d configDir) RootPath() string {
	return d.path
}

func (d configDir) WalletsPath() string {
	return filepath.Join(d.path, "wallets")
}

func (d configDir) WalletPath(name string) string {
	return filepath.Join(d.path, "wallets", name)
}

func (d configDir) WalletContent(name string) string {
	buf, err := ioutil.ReadFile(d.WalletPath(name))
	if err != nil {
		panic(err)
	}
	return string(buf)
}

func (d configDir) Create() {
	err := os.MkdirAll(d.path, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (d configDir) Remove() {
	err := os.RemoveAll(d.path)
	if err != nil {
		panic(err)
	}
}

