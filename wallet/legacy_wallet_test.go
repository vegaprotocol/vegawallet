package wallet_test

import (
	"encoding/json"
	"testing"

	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWallet(t *testing.T) {
	t.Run("Tainting key pair succeeds", testWalletTaintingKeyPairSucceeds)
	t.Run("Tainting key pair that is already tainted fails", testWalletTaintingKeyThatIsAlreadyTaintedFails)
	t.Run("Updating key pair meta succeeds", testWalletUpdatingKeyPairMetaSucceeds)
	t.Run("Updating key pair meta with non-existing public key fails", testWalletUpdatingKeyPairMetaWithNonExistingPublicKeyFails)
	t.Run("Signing transaction request (v2) succeeds", testWalletSigningTxV2Succeeds)
	t.Run("Signing transaction request (v2) with tainted key fails", testWalletSigningTxV2WithTaintedKeyFails)
}

func testWalletTaintingKeyPairSucceeds(t *testing.T) {
	// given
	kp := generateKeyPair()
	name := "jeremy"

	// setup
	w := wallet.NewLegacyWallet(name)
	w.KeyRing.Upsert(*kp)

	// when
	err := w.TaintKey(kp.Pub)

	// then
	require.NoError(t, err)

	// when
	keyPair, err := w.KeyRing.FindPair(kp.Pub)

	// then
	require.NoError(t, err)
	assert.NotNil(t, keyPair)
	assert.True(t, keyPair.Tainted)
}

func testWalletTaintingKeyThatIsAlreadyTaintedFails(t *testing.T) {
	// given
	kp := generateKeyPair()
	kp.Tainted = true
	name := "jeremy"

	// setup
	w := wallet.NewLegacyWallet(name)
	w.KeyRing.Upsert(*kp)

	// when
	err := w.TaintKey(kp.Pub)

	// then
	assert.EqualError(t, err, wallet.ErrPubKeyAlreadyTainted.Error())

	// when
	keyPair, err := w.KeyRing.FindPair(kp.Pub)

	// then
	require.NoError(t, err)
	assert.NotNil(t, keyPair)
	assert.True(t, keyPair.Tainted)
}

func testWalletUpdatingKeyPairMetaSucceeds(t *testing.T) {
	// given
	kp := generateKeyPair()
	name := "jeremy"
	meta := []wallet.Meta{{Key: "primary", Value: "yes"}}

	// setup
	w := wallet.NewLegacyWallet(name)
	w.KeyRing.Upsert(*kp)

	// when
	err := w.UpdateMeta(kp.Pub, meta)

	// then
	require.NoError(t, err)

	// when
	keyPair, err := w.KeyRing.FindPair(kp.Pub)

	// then
	require.NoError(t, err)
	assert.NotNil(t, keyPair)
	assert.Equal(t, meta, keyPair.Meta)
}

func testWalletUpdatingKeyPairMetaWithNonExistingPublicKeyFails(t *testing.T) {
	// given
	kp := generateKeyPair()
	name := "jeremy"
	meta := []wallet.Meta{{Key: "primary", Value: "yes"}}

	// setup
	w := wallet.NewLegacyWallet(name)

	// when
	err := w.UpdateMeta(kp.Pub, meta)

	// then
	require.Error(t, err, wallet.ErrWalletDoesNotExists)
}

func testWalletSigningTxV2Succeeds(t *testing.T) {
	// given
	kp := generateKeyPair()
	name := "jeremy"
	data := []byte("Je ne connaîtrai pas la peur car la peur tue l'esprit.")

	// setup
	w := wallet.NewLegacyWallet(name)
	w.KeyRing.Upsert(*kp)

	// when
	signature, err := w.SignTxV2(kp.Pub, data)

	// then
	require.NoError(t, err)
	assert.Equal(t, kp.Algorithm.Version(), signature.Version)
	assert.Equal(t, kp.Algorithm.Name(), signature.Algo)
	assert.NotEmpty(t, signature.Value)
}

func testWalletSigningTxV2WithTaintedKeyFails(t *testing.T) {
	// given
	kp := generateKeyPair()
	kp.Tainted = true
	name := "jeremy"
	data := []byte("Je ne connaîtrai pas la peur car la peur tue l'esprit.")

	// setup
	w := wallet.NewLegacyWallet(name)
	w.KeyRing.Upsert(*kp)

	// when
	signature, err := w.SignTxV2(kp.Pub, data)

	// then
	require.EqualError(t, err, wallet.ErrPubKeyIsTainted.Error())
	assert.Nil(t, signature)
}

func TestMarshalWallet(t *testing.T) {
	w := wallet.NewLegacyWallet("jeremy")
	w.KeyRing = append(w.KeyRing, wallet.NewKeypair(crypto.NewEd25519(), []byte{1, 2, 3, 4}, []byte{4, 3, 2, 1}))
	expected := `{"Owner":"jeremy","Keypairs":[{"pub":"01020304","priv":"04030201","algo":"vega/ed25519","tainted":false,"meta":null}]}`
	m, err := json.Marshal(&w)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(m))
}

func TestUnMarshalWallet(t *testing.T) {
	w := wallet.LegacyWallet{}
	marshalled := `{"Owner":"jeremy","Keypairs":[{"pub":"01020304","priv":"04030201","algo":"vega/ed25519","tainted":false,"meta":null}]}`
	err := json.Unmarshal([]byte(marshalled), &w)
	assert.NoError(t, err)
	assert.Len(t, w.KeyRing, 1)
	assert.Equal(t, "01020304", w.KeyRing[0].Pub)
	assert.Equal(t, "04030201", w.KeyRing[0].Priv)
	assert.Equal(t, "vega/ed25519", w.KeyRing[0].Algorithm.Name())
}

func TestUnMarshalWalletErrorInvalidAlgorithm(t *testing.T) {
	w := wallet.LegacyWallet{}
	marshalled := `{"Owner":"jeremy","Keypairs":[{"pub":"01020304","priv":"04030201","algo":"notanalgorithm","tainted":false,"meta":null}]}`
	err := json.Unmarshal([]byte(marshalled), &w)
	assert.EqualError(t, err, crypto.ErrUnsupportedSignatureAlgorithm.Error())
}

func generateKeyPair() *wallet.KeyPair {
	kp, err := wallet.GenKeyPair(crypto.Ed25519)
	if err != nil {
		panic(err)
	}
	return kp
}
