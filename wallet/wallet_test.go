package wallet_test

import (
	"os"
	"testing"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/wallet"
	wcrypto "code.vegaprotocol.io/go-wallet/wallet/crypto"

	"github.com/stretchr/testify/assert"
)

var (
	rootDirPath = "/tmp/vegatests/wallet/"
)

func TestWallet(t *testing.T) {
	t.Run("create a wallet success", testCreateWallet)
	t.Run("create a wallet failure", testCreateWalletFailure)
	t.Run("read a wallet success", testReadWallet)
	t.Run("read a wallet failure invalid passphrase", testReadWalletFailureInvalidPassphrase)
	t.Run("read a wallet failure does not exist", testReadWalletFailureDoesNotExist)
	t.Run("add a keypair to a wallet", testAddKeyPairToWallet)
}

func testCreateWallet(t *testing.T) {
	rootDir := rootDir()
	fsutil.EnsureDir(rootDir)
	wallet.EnsureBaseFolder(rootDir)

	w, err := wallet.Create(rootDir, "jeremy", "thisisasecurepassphraseinnit")
	assert.NoError(t, err)
	assert.NotNil(t, w)

	assert.NoError(t, os.RemoveAll(rootDir))
}

func testCreateWalletFailure(t *testing.T) {
	rootDir := rootDir()
	fsutil.EnsureDir(rootDir)
	wallet.EnsureBaseFolder(rootDir)

	// fist time is fine
	w, err := wallet.Create(rootDir, "jeremy", "thisisasecurepassphraseinnit")
	assert.NoError(t, err)
	assert.NotNil(t, w)

	w, err = wallet.Create(rootDir, "jeremy", "whateverthepassphrase")
	assert.EqualError(t, err, wallet.ErrWalletAlreadyExists.Error())
	assert.Nil(t, w)

	assert.NoError(t, os.RemoveAll(rootDir))
}

func testReadWallet(t *testing.T) {
	rootDir := rootDir()
	fsutil.EnsureDir(rootDir)
	wallet.EnsureBaseFolder(rootDir)

	w, err := wallet.Create(rootDir, "jeremy", "thisisasecurepassphraseinnit")
	assert.NoError(t, err)
	assert.NotNil(t, w)

	// no try to read the same wallet
	w1, err1 := wallet.Read(rootDir, "jeremy", "thisisasecurepassphraseinnit")
	assert.NoError(t, err1)
	assert.NotNil(t, w1)

	assert.NoError(t, os.RemoveAll(rootDir))
}

func testReadWalletFailureDoesNotExist(t *testing.T) {
	rootDir := rootDir()
	fsutil.EnsureDir(rootDir)
	wallet.EnsureBaseFolder(rootDir)

	// try to read a wallet which do not exists
	w, err := wallet.Read(rootDir, "jeremy", "thisisasecurepassphraseinnit")
	assert.EqualError(t, err, wallet.ErrWalletDoesNotExists.Error())
	assert.Nil(t, w)

	assert.NoError(t, os.RemoveAll(rootDir))
}

func testReadWalletFailureInvalidPassphrase(t *testing.T) {
	rootDir := rootDir()
	fsutil.EnsureDir(rootDir)
	wallet.EnsureBaseFolder(rootDir)

	w, err := wallet.Create(rootDir, "jeremy", "thisisasecurepassphraseinnit")
	assert.NoError(t, err)
	assert.NotNil(t, w)

	// no try to read the same wallet
	w1, err1 := wallet.Read(rootDir, "jeremy", "thisisasecurepassph")
	assert.EqualError(t, err1, "cipher: message authentication failed")
	assert.Nil(t, w1)

	assert.NoError(t, os.RemoveAll(rootDir))
}

func testAddKeyPairToWallet(t *testing.T) {
	rootDir := rootDir()
	fsutil.EnsureDir(rootDir)
	wallet.EnsureBaseFolder(rootDir)

	w, err := wallet.Create(rootDir, "jeremy", "thisisasecurepassphraseinnit")
	assert.NoError(t, err)
	assert.NotNil(t, w)

	// now try to read the same wallet
	w1, err1 := wallet.Read(rootDir, "jeremy", "thisisasecurepassphraseinnit")
	assert.NoError(t, err1)
	assert.NotNil(t, w1)
	assert.Len(t, w1.Keypairs, 0)

	// create the keypair
	kp := wallet.NewKeypair(wcrypto.NewEd25519(), []byte{1, 2, 3, 255}, []byte{253, 3, 2, 1})

	// now try to add the keypair to the wallet
	w2, err2 := wallet.AddKeypair(&kp, rootDir, "jeremy", "thisisasecurepassphraseinnit")
	assert.NoError(t, err2)
	assert.NotNil(t, w2)

	// now try to read the same wallet a last time, and make sur the values are right
	w3, err3 := wallet.Read(rootDir, "jeremy", "thisisasecurepassphraseinnit")
	assert.NoError(t, err3)
	assert.NotNil(t, w3)
	assert.Len(t, w3.Keypairs, 1)
	// check the hex value of the keypair
	assert.Equal(t, "010203ff", w3.Keypairs[0].Pub)
	assert.Equal(t, "fd030201", w3.Keypairs[0].Priv)

	assert.NoError(t, os.RemoveAll(rootDir))
}
