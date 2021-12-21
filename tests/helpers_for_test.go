package tests_test

import (
	"fmt"
	"path/filepath"
	"testing"

	vgfs "code.vegaprotocol.io/shared/libs/fs"
	vgrand "code.vegaprotocol.io/shared/libs/rand"
)

const testRecoveryPhrase = "swing ceiling chaos green put insane ripple desk match tip melt usual shrug turkey renew icon parade veteran lens govern path rough page render"

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

func FakeNetwork(name string) string {
	return fmt.Sprintf(`
Name = "%s"
Level = "info"
TokenExpiry = "1h0m0s"
Port = 8000
Host = "127.0.0.1"

[API.GRPC]
Retries = 5
Hosts = [
    "example.com:3007",
]

[API.REST]
Hosts = [
    "https://example.com/rest"
]

[API.GraphQL]
Hosts = [
    "https://example.com/gql/query"
]

[Console]
URL = "console.example.com"
LocalPort = 1847
`, name)
}

func RandomPath() string {
	return filepath.Join(vgrand.RandomStr(5), vgrand.RandomStr(5))
}
