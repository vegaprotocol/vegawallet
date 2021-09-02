package v1_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
)

type walletsDir struct {
	path string
}

func newWalletsDir() walletsDir {
	rootPath := filepath.Join("/tmp", "vegatests", "wallet", vgrand.RandomStr(10))

	return walletsDir{
		path: rootPath,
	}
}

func (d walletsDir) WalletsPath() string {
	return d.path
}

func (d walletsDir) WalletPath(name string) string {
	return filepath.Join(d.path, name)
}

func (d walletsDir) WalletContent(name string) string {
	buf, err := ioutil.ReadFile(d.WalletPath(name))
	if err != nil {
		panic(err)
	}
	return string(buf)
}

func (d walletsDir) Remove() {
	err := os.RemoveAll(d.path)
	if err != nil {
		panic(err)
	}
}
