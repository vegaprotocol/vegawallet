package tests_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	vgfs "code.vegaprotocol.io/shared/libs/fs"
	vgrand "code.vegaprotocol.io/shared/libs/rand"
)

func NewTempDir(t *testing.T) (string, func(t *testing.T)) {
	t.Helper()
	uniqueFolderName := vgrand.RandomStr(10)
	home := filepath.Join("/tmp", "vegawallet", uniqueFolderName)
	if err := vgfs.EnsureDir(home); err != nil {
		t.Fatalf("couldn't create Vega home: %v", err)
	}
	return home, func(t *testing.T) {
		t.Helper()
		if err := os.RemoveAll(home); err != nil {
			t.Fatalf("couldn't remove Vega home: %v", err)
		}
	}
}

func NewPassphraseFile(t *testing.T, path string) (string, string) {
	t.Helper()
	passphrase := vgrand.RandomStr(10)
	passphraseFilePath := NewFile(t, path, "passphrase.txt", passphrase)
	return passphrase, passphraseFilePath
}

func NewFile(t *testing.T, path, fileName, data string) string {
	t.Helper()
	filePath := filepath.Join(path, fileName)
	if err := vgfs.WriteFile(filePath, []byte(data)); err != nil {
		t.Fatalf("couldn't write passphrase file: %v", err)
	}
	return filePath
}

func DefaultMetaName(t *testing.T, wallet string, index int) string {
	t.Helper()
	return fmt.Sprintf("%s key %d", wallet, index)
}
