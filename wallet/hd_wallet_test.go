package wallet_test

import (
	"encoding/json"
	"testing"

	"code.vegaprotocol.io/go-wallet/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	TestMnemonic1 = "swing ceiling chaos green put insane ripple desk match tip melt usual shrug turkey renew icon parade veteran lens govern path rough page render"
	TestMnemonic2 = "green put insane ripple desk match tip melt usual shrug turkey renew icon parade veteran lens govern path rough page render swing ceiling chaos"
)

func TestHDWallet(t *testing.T) {
	t.Run("Creating wallet succeeds", testHDWalletCreateWalletSucceeds)
	t.Run("Importing wallet succeeds", testHDWalletImportingWalletSucceeds)
	t.Run("Importing wallet with invalid mnemonic fails", testHDWalletImportingWalletWithInvalidMnemonicFails)
	t.Run("Tainting key pair succeeds", testHDWalletTaintingKeyPairSucceeds)
	t.Run("Tainting key pair that is already tainted fails", testHDWalletTaintingKeyThatIsAlreadyTaintedFails)
	t.Run("Tainting unknown key pair fails", testHDWalletTaintingUnknownKeyFails)
	t.Run("Updating key pair meta succeeds", testHDWalletUpdatingKeyPairMetaSucceeds)
	t.Run("Updating key pair meta with unknown public key fails", testHDWalletUpdatingKeyPairMetaWithUnknownPublicKeyFails)
	t.Run("Describing public key succeeds", testHDWalletDescribingPublicKeysSucceeds)
	t.Run("Describing unknown public key fails", testHDWalletDescribingUnknownPublicKeysFails)
	t.Run("Listing public keys succeeds", testHDWalletListingPublicKeysSucceeds)
	t.Run("Listing key pairs succeeds", testHDWalletListingKeyPairsSucceeds)
	t.Run("Signing transaction request (v1) succeeds", testHDWalletSigningTxV1Succeeds)
	t.Run("Signing transaction request (v1) with tainted key fails", testHDWalletSigningTxV1WithTaintedKeyFails)
	t.Run("Signing transaction request (v1) with unknown key fails", testHDWalletSigningTxV1WithUnknownKeyFails)
	t.Run("Signing transaction request (v2) succeeds", testHDWalletSigningTxV2Succeeds)
	t.Run("Signing transaction request (v2) with tainted key fails", testHDWalletSigningTxV2WithTaintedKeyFails)
	t.Run("Signing transaction request (v2) with unknown key fails", testHDWalletSigningTxV2WithUnknownKeyFails)
	t.Run("Signing any message succeeds", testHDWalletSigningAnyMessageSucceeds)
	t.Run("Signing any message with tainted key fails", testHDWalletSigningAnyMessageWithTaintedKeyFails)
	t.Run("Signing any message with unknown key fails", testHDWalletSigningAnyMessageWithUnknownKeyFails)
	t.Run("Verifying any message succeeds", testHDWalletVerifyingAnyMessageSucceeds)
	t.Run("Verifying any message with unknown key fails", testHDWalletVerifyingAnyMessageWithUnknownKeyFails)
	t.Run("Marshaling wallet succeeds", testHDWalletMarshalingSucceeds)
	t.Run("Unmarshaling wallet succeeds", testHDWalletUnmarshalingWalletSucceeds)
}

func testHDWalletCreateWalletSucceeds(t *testing.T) {
	// given
	name := "jeremy"

	// when
	w, mnemonic, err := wallet.NewHDWallet(name)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, mnemonic)
	assert.NotNil(t, w)
}

func testHDWalletImportingWalletSucceeds(t *testing.T) {
	// given
	name := "jeremy"

	// when
	w, err := wallet.ImportHDWallet(name, TestMnemonic1)

	// then
	require.NoError(t, err)
	assert.NotNil(t, w)
}

func testHDWalletImportingWalletWithInvalidMnemonicFails(t *testing.T) {
	// given
	name := "jeremy"

	// when
	w, err := wallet.ImportHDWallet(name, "vladimir harkonnen doesn't like trees")

	// then
	require.EqualError(t, err, wallet.ErrInvalidMnemonic.Error())
	assert.Nil(t, w)
}

func testHDWalletTaintingKeyPairSucceeds(t *testing.T) {
	// given
	name := "jeremy"

	// when
	w, mnemonic, err := wallet.NewHDWallet(name)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, mnemonic)
	assert.NotNil(t, w)

	// when
	kp, err := w.GenerateKeyPair()

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp)

	// when
	err = w.TaintKey(kp.PublicKey())

	// then
	require.NoError(t, err)

	// when
	pubKey, err := w.DescribePublicKey(kp.PublicKey())

	// then
	require.NoError(t, err)
	assert.NotNil(t, pubKey)
	assert.True(t, pubKey.IsTainted())
}

func testHDWalletTaintingKeyThatIsAlreadyTaintedFails(t *testing.T) {
	// given
	name := "jeremy"

	// when
	w, mnemonic, err := wallet.NewHDWallet(name)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, mnemonic)
	assert.NotNil(t, w)

	// when
	kp, err := w.GenerateKeyPair()

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp)

	// when
	err = w.TaintKey(kp.PublicKey())

	// then
	require.NoError(t, err)

	// when
	err = w.TaintKey(kp.PublicKey())

	// then
	assert.EqualError(t, err, wallet.ErrPubKeyAlreadyTainted.Error())

	// when
	pubKey, err := w.DescribePublicKey(kp.PublicKey())

	// then
	require.NoError(t, err)
	assert.NotNil(t, pubKey)
	assert.True(t, pubKey.IsTainted())
}

func testHDWalletTaintingUnknownKeyFails(t *testing.T) {
	// given
	name := "jeremy"

	// when
	w, mnemonic, err := wallet.NewHDWallet(name)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, mnemonic)
	assert.NotNil(t, w)

	// when
	err = w.TaintKey("vladimirharkonnen")

	// then
	assert.EqualError(t, err, wallet.ErrPubKeyDoesNotExist.Error())
}

func testHDWalletUpdatingKeyPairMetaSucceeds(t *testing.T) {
	// given
	name := "jeremy"
	meta := []wallet.Meta{{Key: "primary", Value: "yes"}}

	// when
	w, mnemonic, err := wallet.NewHDWallet(name)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, mnemonic)
	assert.NotNil(t, w)

	// when
	kp, err := w.GenerateKeyPair()

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp)

	// when
	err = w.UpdateMeta(kp.PublicKey(), meta)

	// then
	require.NoError(t, err)

	// when
	pubKey, err := w.DescribePublicKey(kp.PublicKey())

	// then
	require.NoError(t, err)
	assert.NotNil(t, pubKey)
	assert.Equal(t, meta, pubKey.Meta())
}

func testHDWalletUpdatingKeyPairMetaWithUnknownPublicKeyFails(t *testing.T) {
	// given
	name := "jeremy"
	meta := []wallet.Meta{{Key: "primary", Value: "yes"}}

	// when
	w, mnemonic, err := wallet.NewHDWallet(name)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, mnemonic)
	assert.NotNil(t, w)

	// when
	err = w.UpdateMeta("somekey", meta)

	// then
	require.Error(t, err, wallet.ErrWalletDoesNotExists)
}

func testHDWalletDescribingPublicKeysSucceeds(t *testing.T) {
	// given
	name := "jeremy"

	// when
	w, err := wallet.ImportHDWallet(name, TestMnemonic1)

	// then
	require.NoError(t, err)
	assert.NotNil(t, w)

	// when
	kp1, err := w.GenerateKeyPair()

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp1)

	// when
	pubKey, err := w.DescribePublicKey(kp1.PublicKey())

	// then
	require.NoError(t, err)
	assert.Equal(t, kp1.PublicKey(), pubKey.Key())
	assert.Equal(t, kp1.Meta(), pubKey.Meta())
	assert.Equal(t, kp1.IsTainted(), pubKey.IsTainted())
	assert.Equal(t, kp1.AlgorithmName(), pubKey.AlgorithmName())
	assert.Equal(t, kp1.AlgorithmVersion(), pubKey.AlgorithmVersion())
}

func testHDWalletDescribingUnknownPublicKeysFails(t *testing.T) {
	// given
	name := "jeremy"

	// when
	w, err := wallet.ImportHDWallet(name, TestMnemonic1)

	// then
	require.NoError(t, err)
	assert.NotNil(t, w)

	// when
	pubKey, err := w.DescribePublicKey("vladimirharkonnen")

	// then
	require.EqualError(t, err, wallet.ErrPubKeyDoesNotExist.Error())
	assert.Empty(t, pubKey)
}

func testHDWalletListingPublicKeysSucceeds(t *testing.T) {
	// given
	name := "jeremy"

	// when
	w, err := wallet.ImportHDWallet(name, TestMnemonic1)

	// then
	require.NoError(t, err)
	assert.NotNil(t, w)

	// when
	kp1, err := w.GenerateKeyPair()

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp1)

	// when
	kp2, err := w.GenerateKeyPair()

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp2)

	// when
	keys := w.ListPublicKeys()

	// then
	assert.Len(t, keys, 2)
	assert.Equal(t, keys[0].Key(), kp1.PublicKey())
	assert.Equal(t, keys[1].Key(), kp2.PublicKey())
}

func testHDWalletListingKeyPairsSucceeds(t *testing.T) {
	// given
	name := "jeremy"

	// when
	w, err := wallet.ImportHDWallet(name, TestMnemonic1)

	// then
	require.NoError(t, err)
	assert.NotNil(t, w)

	// when
	kp1, err := w.GenerateKeyPair()

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp1)

	// when
	kp2, err := w.GenerateKeyPair()

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp2)

	// when
	keys := w.ListKeyPairs()

	// then
	assert.Equal(t, keys, []wallet.KeyPair{kp1, kp2})
}

func testHDWalletSigningTxV1Succeeds(t *testing.T) {
	// given
	name := "jeremy"
	data := []byte("Je ne connaîtrai pas la peur car la peur tue l'esprit.")

	// when
	w, err := wallet.ImportHDWallet(name, TestMnemonic1)

	// then
	require.NoError(t, err)
	assert.NotNil(t, w)

	// when
	kp, err := w.GenerateKeyPair()

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp)

	// when
	signature, err := w.SignTxV1(kp.PublicKey(), data, 1)

	// then
	require.NoError(t, err)
	assert.NotEmpty(t, signature.Tx)
	assert.Equal(t, kp.AlgorithmVersion(), signature.Sig.Version)
	assert.Equal(t, kp.AlgorithmName(), signature.Sig.Algo)
	assert.NotEmpty(t, signature.Sig.Sig)
}

func testHDWalletSigningTxV1WithTaintedKeyFails(t *testing.T) {
	// given
	name := "jeremy"
	data := []byte("Je ne connaîtrai pas la peur car la peur tue l'esprit.")

	// when
	w, err := wallet.ImportHDWallet(name, TestMnemonic1)

	// then
	require.NoError(t, err)
	assert.NotNil(t, w)

	// when
	kp, err := w.GenerateKeyPair()

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp)

	// when
	err = w.TaintKey(kp.PublicKey())

	// then
	require.NoError(t, err)

	// when
	signature, err := w.SignTxV1(kp.PublicKey(), data, 1)

	// then
	require.EqualError(t, err, wallet.ErrPubKeyIsTainted.Error())
	assert.Empty(t, signature)
}

func testHDWalletSigningTxV1WithUnknownKeyFails(t *testing.T) {
	// given
	name := "jeremy"
	data := []byte("Je ne connaîtrai pas la peur car la peur tue l'esprit.")

	// when
	w, err := wallet.ImportHDWallet(name, TestMnemonic1)

	// then
	require.NoError(t, err)
	assert.NotNil(t, w)

	// when
	kp, err := w.GenerateKeyPair()

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp)

	// when
	signature, err := w.SignTxV1("vladimirharkonnen", data, 1)

	// then
	require.EqualError(t, err, wallet.ErrPubKeyDoesNotExist.Error())
	assert.Empty(t, signature)
}

func testHDWalletSigningTxV2Succeeds(t *testing.T) {
	// given
	name := "jeremy"
	data := []byte("Je ne connaîtrai pas la peur car la peur tue l'esprit.")

	// when
	w, err := wallet.ImportHDWallet(name, TestMnemonic1)

	// then
	require.NoError(t, err)
	assert.NotNil(t, w)

	// when
	kp, err := w.GenerateKeyPair()

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp)

	// when
	signature, err := w.SignTxV2(kp.PublicKey(), data)

	// then
	require.NoError(t, err)
	assert.Equal(t, kp.AlgorithmVersion(), signature.Version)
	assert.Equal(t, kp.AlgorithmName(), signature.Algo)
	assert.Equal(t, "3849965c2f327f0b148e3b122cdc89a17fa07611e2a4178b1605dea5442ab7cfadb35d0b0ef527522f6477a5633b8f22d3b2d1e619d306111b7851a9d6100d02", signature.Value)
}

func testHDWalletSigningTxV2WithTaintedKeyFails(t *testing.T) {
	// given
	name := "jeremy"
	data := []byte("Je ne connaîtrai pas la peur car la peur tue l'esprit.")

	// when
	w, err := wallet.ImportHDWallet(name, TestMnemonic1)

	// then
	require.NoError(t, err)
	assert.NotNil(t, w)

	// when
	kp, err := w.GenerateKeyPair()

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp)

	// when
	err = w.TaintKey(kp.PublicKey())

	// then
	require.NoError(t, err)

	// when
	signature, err := w.SignTxV2(kp.PublicKey(), data)

	// then
	require.EqualError(t, err, wallet.ErrPubKeyIsTainted.Error())
	assert.Nil(t, signature)
}

func testHDWalletSigningTxV2WithUnknownKeyFails(t *testing.T) {
	// given
	name := "jeremy"
	data := []byte("Je ne connaîtrai pas la peur car la peur tue l'esprit.")

	// when
	w, err := wallet.ImportHDWallet(name, TestMnemonic1)

	// then
	require.NoError(t, err)
	assert.NotNil(t, w)

	// when
	kp, err := w.GenerateKeyPair()

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp)

	// when
	signature, err := w.SignTxV2("vladimirharkonnen", data)

	// then
	require.EqualError(t, err, wallet.ErrPubKeyDoesNotExist.Error())
	assert.Empty(t, signature)
}

func testHDWalletSigningAnyMessageSucceeds(t *testing.T) {
	// given
	name := "jeremy"
	data := []byte("Je ne connaîtrai pas la peur car la peur tue l'esprit.")

	// when
	w, err := wallet.ImportHDWallet(name, TestMnemonic1)

	// then
	require.NoError(t, err)
	assert.NotNil(t, w)

	// when
	kp, err := w.GenerateKeyPair()

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp)

	// when
	signature, err := w.SignAny(kp.PublicKey(), data)

	// then
	require.NoError(t, err)
	assert.Equal(t, []byte{0x38, 0x49, 0x96, 0x5c, 0x2f, 0x32, 0x7f, 0xb, 0x14, 0x8e, 0x3b, 0x12, 0x2c, 0xdc, 0x89, 0xa1, 0x7f, 0xa0, 0x76, 0x11, 0xe2, 0xa4, 0x17, 0x8b, 0x16, 0x5, 0xde, 0xa5, 0x44, 0x2a, 0xb7, 0xcf, 0xad, 0xb3, 0x5d, 0xb, 0xe, 0xf5, 0x27, 0x52, 0x2f, 0x64, 0x77, 0xa5, 0x63, 0x3b, 0x8f, 0x22, 0xd3, 0xb2, 0xd1, 0xe6, 0x19, 0xd3, 0x6, 0x11, 0x1b, 0x78, 0x51, 0xa9, 0xd6, 0x10, 0xd, 0x2}, signature)
}

func testHDWalletSigningAnyMessageWithTaintedKeyFails(t *testing.T) {
	// given
	name := "jeremy"
	data := []byte("Je ne connaîtrai pas la peur car la peur tue l'esprit.")

	// when
	w, err := wallet.ImportHDWallet(name, TestMnemonic1)

	// then
	require.NoError(t, err)
	assert.NotNil(t, w)

	// when
	kp, err := w.GenerateKeyPair()

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp)

	// when
	err = w.TaintKey(kp.PublicKey())

	// then
	require.NoError(t, err)

	// when
	signature, err := w.SignAny(kp.PublicKey(), data)

	// then
	require.EqualError(t, err, wallet.ErrPubKeyIsTainted.Error())
	assert.Empty(t, signature)
}

func testHDWalletSigningAnyMessageWithUnknownKeyFails(t *testing.T) {
	// given
	name := "jeremy"
	data := []byte("Je ne connaîtrai pas la peur car la peur tue l'esprit.")

	// when
	w, err := wallet.ImportHDWallet(name, TestMnemonic1)

	// then
	require.NoError(t, err)
	assert.NotNil(t, w)

	// when
	signature, err := w.SignAny("vladimirharkonnen", data)

	// then
	require.EqualError(t, err, wallet.ErrPubKeyDoesNotExist.Error())
	assert.Empty(t, signature)
}

func testHDWalletVerifyingAnyMessageSucceeds(t *testing.T) {
	// given
	name := "jeremy"
	data := []byte("Je ne connaîtrai pas la peur car la peur tue l'esprit.")
	sig := []byte{0x38, 0x49, 0x96, 0x5c, 0x2f, 0x32, 0x7f, 0xb, 0x14, 0x8e, 0x3b, 0x12, 0x2c, 0xdc, 0x89, 0xa1, 0x7f, 0xa0, 0x76, 0x11, 0xe2, 0xa4, 0x17, 0x8b, 0x16, 0x5, 0xde, 0xa5, 0x44, 0x2a, 0xb7, 0xcf, 0xad, 0xb3, 0x5d, 0xb, 0xe, 0xf5, 0x27, 0x52, 0x2f, 0x64, 0x77, 0xa5, 0x63, 0x3b, 0x8f, 0x22, 0xd3, 0xb2, 0xd1, 0xe6, 0x19, 0xd3, 0x6, 0x11, 0x1b, 0x78, 0x51, 0xa9, 0xd6, 0x10, 0xd, 0x2}

	// when
	w, err := wallet.ImportHDWallet(name, TestMnemonic1)

	// then
	require.NoError(t, err)
	assert.NotNil(t, w)

	// when
	kp, err := w.GenerateKeyPair()

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp)

	// when
	verified, err := w.VerifyAny(kp.PublicKey(), data, sig)

	// then
	require.NoError(t, err)
	assert.True(t, verified)
}

func testHDWalletVerifyingAnyMessageWithUnknownKeyFails(t *testing.T) {
	// given
	name := "jeremy"
	data := []byte("Je ne connaîtrai pas la peur car la peur tue l'esprit.")
	sig := []byte{0xd5, 0xc4, 0x9e, 0xfd, 0x13, 0x73, 0x9b, 0xdd, 0x36, 0x81, 0x75, 0xcc, 0x59, 0xc8, 0xbe, 0xe1, 0x20, 0x25, 0xe4, 0xb9, 0x14, 0x7a, 0x22, 0xbb, 0xa4, 0x84, 0xef, 0x7e, 0xe7, 0x2f, 0x55, 0x13, 0x5f, 0x52, 0x55, 0xad, 0x90, 0x35, 0x67, 0x6c, 0x91, 0x9d, 0xbb, 0x91, 0x21, 0x1f, 0x98, 0x53, 0xcc, 0x68, 0xe, 0x58, 0x5b, 0x4c, 0x26, 0xd7, 0xea, 0x20, 0x1, 0x50, 0x6c, 0x41, 0xcb, 0x3}

	// when
	w, err := wallet.ImportHDWallet(name, TestMnemonic1)

	// then
	require.NoError(t, err)
	assert.NotNil(t, w)

	// when
	signature, err := w.VerifyAny("vladimirharkonnen", data, sig)

	// then
	require.EqualError(t, err, wallet.ErrPubKeyDoesNotExist.Error())
	assert.Empty(t, signature)
}

func testHDWalletMarshalingSucceeds(t *testing.T) {
	// given
	name := "jeremy"

	// when
	w, err := wallet.ImportHDWallet(name, TestMnemonic1)

	// then
	require.NoError(t, err)
	assert.NotNil(t, w)

	// when
	kp, err := w.GenerateKeyPair()

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp)

	// when
	m, err := json.Marshal(&w)

	// then
	assert.NoError(t, err)
	expected := `{"version":1,"name":"jeremy","node":"PjI6zxEu4dtcTu92dYlB/2Da+rvSpg7KzvmLMQ9wv6i6n75/ftik1rPYiZ/nTfBzqVttvNnoswyldTjPCjV5kw==","keys":[{"index":1,"public_key":"30ebce58d94ad37c4ff6a9014c955c20e12468da956163228cc7ec9b98d3a371","private_key":"1bbd4efb460d0bf457251e866697d5d2e9b58c5dcb96a964cd9cfff1a712a2b930ebce58d94ad37c4ff6a9014c955c20e12468da956163228cc7ec9b98d3a371","meta":null,"tainted":false,"algorithm":{"name":"vega/ed25519","version":1}}]}`
	assert.Equal(t, expected, string(m))
}

func testHDWalletUnmarshalingWalletSucceeds(t *testing.T) {
	// given
	w := wallet.HDWallet{}
	marshalled := `{"version":1,"name":"jeremy","node":"CZ13XhuFZ8K7TxNTAdKmMXh+OIVX6TFxTToXgnAqGlcO5eTY/5AVqZkWRIU3zfr8hvE7i2yIYAB6HT28ibi1fg==","keys":[{"index":1,"public_key":"e4997f2886f3f0fae5c4353f45c50560e93971e00b3b9350ede8abd491b5fbde","private_key":"07757d2c86c98e36c041d2a8a0fdd7c70bb3e88794328d27a9a29c159a38b23fe4997f2886f3f0fae5c4353f45c50560e93971e00b3b9350ede8abd491b5fbde","meta":null,"tainted":false,"algorithm":{"name":"vega/ed25519","version":1}}]}`

	// when
	err := json.Unmarshal([]byte(marshalled), &w)

	// then
	assert.NoError(t, err)
	assert.Equal(t, w.Version(), uint32(1))
	assert.Equal(t, w.Name(), "jeremy")
	keyPairs := w.ListKeyPairs()
	assert.Len(t, keyPairs, 1)
	assert.Equal(t, "e4997f2886f3f0fae5c4353f45c50560e93971e00b3b9350ede8abd491b5fbde", keyPairs[0].PublicKey())
	assert.Equal(t, "07757d2c86c98e36c041d2a8a0fdd7c70bb3e88794328d27a9a29c159a38b23fe4997f2886f3f0fae5c4353f45c50560e93971e00b3b9350ede8abd491b5fbde", keyPairs[0].PrivateKey())
	assert.Equal(t, uint32(1), keyPairs[0].AlgorithmVersion())
	assert.Equal(t, "vega/ed25519", keyPairs[0].AlgorithmName())
	assert.False(t, keyPairs[0].IsTainted())
	assert.Nil(t, keyPairs[0].Meta())
}
