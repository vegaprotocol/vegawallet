package cmd_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	vgfs "code.vegaprotocol.io/shared/libs/fs"
	vgrand "code.vegaprotocol.io/shared/libs/rand"
)

func NewTempDir(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Cleanup(func() {
		if err := os.RemoveAll(home); err != nil {
			t.Fatalf(fmt.Sprintf("couldn't remove test folder: %v", err))
		}
	})
	return home
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
