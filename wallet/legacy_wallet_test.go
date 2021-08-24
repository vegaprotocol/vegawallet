package wallet_test

import (
	"encoding/json"
	"testing"

	"code.vegaprotocol.io/go-wallet/crypto"
	"code.vegaprotocol.io/go-wallet/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLegacyWallet(t *testing.T) {
	t.Run("Tainting key pair succeeds", testLegacyWalletTaintingKeyPairSucceeds)
	t.Run("Tainting key pair that is already tainted fails", testLegacyWalletTaintingKeyThatIsAlreadyTaintedFails)
	t.Run("Updating key pair metadata succeeds", testLegacyWalletUpdatingKeyPairMetaSucceeds)
	t.Run("Updating key pair metadata with non-existing public key fails", testLegacyWalletUpdatingKeyPairMetaWithNonExistingPublicKeyFails)
	t.Run("Signing transaction request (v2) succeeds", testLegacyWalletSigningTxV2Succeeds)
	t.Run("Signing transaction request (v2) with tainted key fails", testLegacyWalletSigningTxV2WithTaintedKeyFails)
}

func testLegacyWalletTaintingKeyPairSucceeds(t *testing.T) {
	// given
	kp := generateLegacyKeyPair()
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

func testLegacyWalletTaintingKeyThatIsAlreadyTaintedFails(t *testing.T) {
	// given
	kp := generateLegacyKeyPair()
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

func testLegacyWalletUpdatingKeyPairMetaSucceeds(t *testing.T) {
	// given
	kp := generateLegacyKeyPair()
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
	assert.Equal(t, meta, keyPair.MetaList)
}

func testLegacyWalletUpdatingKeyPairMetaWithNonExistingPublicKeyFails(t *testing.T) {
	// given
	kp := generateLegacyKeyPair()
	name := "jeremy"
	meta := []wallet.Meta{{Key: "primary", Value: "yes"}}

	// setup
	w := wallet.NewLegacyWallet(name)

	// when
	err := w.UpdateMeta(kp.Pub, meta)

	// then
	require.Error(t, err, wallet.ErrWalletDoesNotExists)
}

func testLegacyWalletSigningTxV2Succeeds(t *testing.T) {
	// given
	kp := generateLegacyKeyPair()
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

func testLegacyWalletSigningTxV2WithTaintedKeyFails(t *testing.T) {
	// given
	kp := generateLegacyKeyPair()
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
	w.KeyRing = append(w.KeyRing, newKeyPair(crypto.NewEd25519(), "01020304", "04030201"))
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

func generateLegacyKeyPair() *wallet.LegacyKeyPair {
	kp, err := wallet.GenKeyPair(crypto.Ed25519, 1)
	if err != nil {
		panic(err)
	}
	return kp
}

func newKeyPair(algo crypto.SignatureAlgorithm, pub, priv string) wallet.LegacyKeyPair {
	return wallet.LegacyKeyPair{
		Algorithm: algo,
		Pub:       pub,
		Priv:      priv,
	}
}
