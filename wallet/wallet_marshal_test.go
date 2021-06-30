package wallet_test

import (
	"encoding/json"
	"testing"

	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestWallet(t *testing.T) {
	t.Run("Tainting key pair succeeds", testWalletTaintingKeyPairSucceeds)
	t.Run("Tainting key pair that is already tainted fails", testWalletTaintingKeyThatIsAlreadyTaintedFails)
}

func testWalletTaintingKeyPairSucceeds(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	pubKey := "0xDEADBEEF"
	name := "jeremy"

	// setup
	w := wallet.NewWallet(name)
	w.KeyRing.Upsert(wallet.KeyPair{
		Pub:       pubKey,
		Priv:      "0xCOFFEDUDE",
		Algorithm: crypto.SignatureAlgorithm{},
		Tainted:   false,
		Meta:      nil,
	})

	// when
	err := w.TaintKey(pubKey)

	// then
	require.NoError(t, err)

	// when
	keyPair, err := w.KeyRing.FindPair(pubKey)

	// then
	require.NoError(t, err)
	assert.NotNil(t, keyPair)
	assert.True(t, keyPair.Tainted)
}


func testWalletTaintingKeyThatIsAlreadyTaintedFails(t *testing.T) {
	h := getTestHandler(t)
	defer h.ctrl.Finish()

	// given
	pubKey := "0xDEADBEEF"
	name := "jeremy"

	// setup
	w := wallet.NewWallet(name)
	w.KeyRing.Upsert(wallet.KeyPair{
		Pub:       pubKey,
		Priv:      "0xCOFFEDUDE",
		Algorithm: crypto.SignatureAlgorithm{},
		Tainted:   true,
		Meta:      nil,
	})

	// when
	err := w.TaintKey(pubKey)

	// then
	assert.EqualError(t, err, wallet.ErrPubKeyAlreadyTainted.Error())

	// when
	keyPair, err := w.KeyRing.FindPair(pubKey)

	// then
	require.NoError(t, err)
	assert.NotNil(t, keyPair)
	assert.True(t, keyPair.Tainted)
}

func TestMarshalWallet(t *testing.T) {
	w := wallet.NewWallet("jeremy")
	w.KeyRing = append(w.KeyRing, wallet.NewKeypair(crypto.NewEd25519(), []byte{1, 2, 3, 4}, []byte{4, 3, 2, 1}))
	expected := `{"Owner":"jeremy","Keypairs":[{"pub":"01020304","priv":"04030201","algo":"vega/ed25519","tainted":false,"meta":null}]}`
	m, err := json.Marshal(&w)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(m))
}

func TestUnMarshalWallet(t *testing.T) {
	w := wallet.Wallet{}
	marshalled := `{"Owner":"jeremy","Keypairs":[{"pub":"01020304","priv":"04030201","algo":"vega/ed25519","tainted":false,"meta":null}]}`
	err := json.Unmarshal([]byte(marshalled), &w)
	assert.NoError(t, err)
	assert.Len(t, w.KeyRing, 1)
	assert.Equal(t, "01020304", w.KeyRing[0].Pub)
	assert.Equal(t, "04030201", w.KeyRing[0].Priv)
	assert.Equal(t, "vega/ed25519", w.KeyRing[0].Algorithm.Name())
}

func TestUnMarshalWalletErrorInvalidAlgorithm(t *testing.T) {
	w := wallet.Wallet{}
	marshalled := `{"Owner":"jeremy","Keypairs":[{"pub":"01020304","priv":"04030201","algo":"notanalgorithm","tainted":false,"meta":null}]}`
	err := json.Unmarshal([]byte(marshalled), &w)
	assert.EqualError(t, err, crypto.ErrUnsupportedSignatureAlgorithm.Error())
}
