package crypto_test

import (
	"testing"

	crypto2 "code.vegaprotocol.io/go-wallet/crypto"
	"github.com/stretchr/testify/assert"
)

func TestEncryption(t *testing.T) {
	t.Run("create encrypt and decrypt success", testEncryptDecryptOK)
	t.Run("decrypt fail wrong passphrase", testDecryptFailWrongPassphrase)
}

func testEncryptDecryptOK(t *testing.T) {
	data := []byte("hello world")
	passphrase := "oh yea?"

	buf, err := crypto2.Encrypt(data, passphrase)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)

	buf1, err := crypto2.Decrypt(buf, passphrase)
	assert.NoError(t, err)
	assert.Equal(t, data, buf1)
}

func testDecryptFailWrongPassphrase(t *testing.T) {
	data := []byte("hello world")
	passphrase := "oh yea?"
	wrongpassphrase := "oh really!"

	buf, err := crypto2.Encrypt(data, passphrase)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)

	buf1, err := crypto2.Decrypt(buf, wrongpassphrase)
	assert.Error(t, err)
	assert.NotEqual(t, data, buf1)
}
