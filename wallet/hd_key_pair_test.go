package wallet_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"code.vegaprotocol.io/vegawallet/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testPublicKey  = "e20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a"
	testPrivateKey = "c55ae95741a0431186bb062501ad9c5a31505fc14d9fd91f08108b11fbec33d9e20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a"
)

func TestHDKeypair(t *testing.T) {
	t.Run("New key pair succeeds", testHDKeyPairNewKeyPairSucceeds)
	t.Run("Deep copying key pair succeeds", testHDKeyPairDeepCopyingKeyPairSucceeds)
	t.Run("Tainting a key pair succeeds", testHDKeyPairTaintingKeyPairSucceeds)
	t.Run("Tainting an already tainted key pair fails", testHDKeyPairTaintingAlreadyTaintedKeyPairFails)
	t.Run("Untainting a key pair succeeds", testHDKeyPairUntaintingKeyPairSucceeds)
	t.Run("Untainting a not-tainted key pair fails", testHDKeyPairUntaintingNotTaintedKeyPairFails)
	t.Run("Secure copy of key pair removes sensitive information", testHDKeyPairToPublicKeyRemovesSensitiveInformation)
	t.Run("Signing a transaction succeeds", testHDKeyPairSigningTransactionSucceeds)
	t.Run("Signing a transaction with tainted key fails", testHDKeyPairSigningTransactionWithTaintedKeyFails)
	t.Run("Signing any message succeeds", testHDKeyPairSigningAnyMessageSucceeds)
	t.Run("Signing any message with tainted key fails", testHDKeyPairSigningAnyMessageWithTaintedKeyFails)
	t.Run("Verifying any message succeeds", testHDKeyPairVerifyingAnyMessageSucceeds)
	t.Run("Verifying any message with invalid signature fails", testHDKeyPairVerifyingAnyMessageWithInvalidSignatureFails)
	t.Run("Marshaling key pair succeeds", testHDKeyPairMarshalingKeyPairSucceeds)
	t.Run("Unmarshaling key pair succeeds", testHDKeyPairUnmarshalingKeyPairSucceeds)
}

func testHDKeyPairNewKeyPairSucceeds(t *testing.T) {
	// given
	publicKey := testPublicKey
	rawPublicKey, err := hex.DecodeString(publicKey)
	if err != nil {
		t.Fatalf("couldn't decode public key: %v", err)
	}
	privateKey := testPrivateKey
	rawPrivateKey, err := hex.DecodeString(privateKey)
	if err != nil {
		t.Fatalf("couldn't decode private key: %v", err)
	}

	// when
	kp, err := wallet.NewHDKeyPair(1, rawPublicKey, rawPrivateKey)

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp)
	assert.Equal(t, uint32(1), kp.Index())
	assert.Equal(t, publicKey, kp.PublicKey())
	assert.Equal(t, privateKey, kp.PrivateKey())
	assert.False(t, kp.IsTainted())
	assert.Equal(t, "vega/ed25519", kp.AlgorithmName())
	assert.Equal(t, uint32(1), kp.AlgorithmVersion())
	assert.Empty(t, kp.Meta())
}

func testHDKeyPairDeepCopyingKeyPairSucceeds(t *testing.T) {
	// given
	kp := generateHDKeyPair(t)

	// when
	copiedKp := kp.DeepCopy()

	// then
	require.NotNil(t, copiedKp)
	assert.NotSame(t, kp, copiedKp)
}

func testHDKeyPairTaintingKeyPairSucceeds(t *testing.T) {
	// given
	kp := generateHDKeyPair(t)

	// when
	err := kp.Taint()

	// then
	require.NoError(t, err)
	assert.True(t, kp.IsTainted())
}

func testHDKeyPairTaintingAlreadyTaintedKeyPairFails(t *testing.T) {
	// given
	kp := generateHDKeyPair(t)

	// when
	err := kp.Taint()

	// then
	require.NoError(t, err)
	assert.True(t, kp.IsTainted())

	// when
	err = kp.Taint()

	// then
	assert.Error(t, err)
	assert.True(t, kp.IsTainted())
}

func testHDKeyPairUntaintingKeyPairSucceeds(t *testing.T) {
	// given
	kp := generateHDKeyPair(t)

	// when
	err := kp.Taint()

	// then
	require.NoError(t, err)
	assert.True(t, kp.IsTainted())

	// when
	err = kp.Untaint()

	// then
	require.NoError(t, err)
	assert.False(t, kp.IsTainted())
}

func testHDKeyPairUntaintingNotTaintedKeyPairFails(t *testing.T) {
	// given
	kp := generateHDKeyPair(t)

	// when
	err := kp.Untaint()

	// then
	assert.Error(t, err)
	assert.False(t, kp.IsTainted())
}

func testHDKeyPairToPublicKeyRemovesSensitiveInformation(t *testing.T) {
	// given
	kp := generateHDKeyPair(t)

	// when
	secureKp := kp.ToPublicKey()

	// then
	assert.Equal(t, uint32(1), secureKp.Index())
	assert.Equal(t, kp.PublicKey(), secureKp.Key())
	assert.Equal(t, kp.IsTainted(), secureKp.IsTainted())
	assert.Equal(t, kp.AlgorithmName(), secureKp.AlgorithmName())
	assert.Equal(t, kp.AlgorithmVersion(), secureKp.AlgorithmVersion())
	assert.Equal(t, kp.Meta(), secureKp.Meta())
}

func testHDKeyPairSigningTransactionSucceeds(t *testing.T) {
	// given
	kp := generateHDKeyPair(t)
	data := []byte("Paul Atreides")

	// when
	sig, err := kp.Sign(data)

	// then
	require.NoError(t, err)
	assert.Equal(t, "2ffd9c1a5c28007eb5fe2fbf7be446cf00d6edee2131a658f4a0424b7fc4cd8ef6a3237a0a9d0355e80eabb2dd2716638a5c545a3b9a2ca4a6c5d26898070501", sig.Value)
	assert.Equal(t, kp.AlgorithmName(), sig.Algo)
	assert.Equal(t, kp.AlgorithmVersion(), sig.Version)
}

func testHDKeyPairSigningTransactionWithTaintedKeyFails(t *testing.T) {
	// given
	kp := generateHDKeyPair(t)
	data := []byte("Paul Atreides")

	// setup
	err := kp.Taint()
	require.NoError(t, err)

	// when
	sig, err := kp.Sign(data)

	// then
	require.Error(t, err)
	assert.Nil(t, sig)
}

func testHDKeyPairSigningAnyMessageSucceeds(t *testing.T) {
	// given
	kp := generateHDKeyPair(t)
	data := []byte("Paul Atreides")

	// when
	sig, err := kp.SignAny(data)

	// then
	require.NoError(t, err)
	assert.Equal(t, []byte{0x2f, 0xfd, 0x9c, 0x1a, 0x5c, 0x28, 0x0, 0x7e, 0xb5, 0xfe, 0x2f, 0xbf, 0x7b, 0xe4, 0x46, 0xcf, 0x0, 0xd6, 0xed, 0xee, 0x21, 0x31, 0xa6, 0x58, 0xf4, 0xa0, 0x42, 0x4b, 0x7f, 0xc4, 0xcd, 0x8e, 0xf6, 0xa3, 0x23, 0x7a, 0xa, 0x9d, 0x3, 0x55, 0xe8, 0xe, 0xab, 0xb2, 0xdd, 0x27, 0x16, 0x63, 0x8a, 0x5c, 0x54, 0x5a, 0x3b, 0x9a, 0x2c, 0xa4, 0xa6, 0xc5, 0xd2, 0x68, 0x98, 0x7, 0x5, 0x1}, sig)
}

func testHDKeyPairSigningAnyMessageWithTaintedKeyFails(t *testing.T) {
	// given
	kp := generateHDKeyPair(t)
	data := []byte("Paul Atreides")

	// setup
	err := kp.Taint()
	require.NoError(t, err)

	// when
	sig, err := kp.SignAny(data)

	// then
	require.Error(t, err)
	assert.Nil(t, sig)
}

func testHDKeyPairVerifyingAnyMessageSucceeds(t *testing.T) {
	// given
	kp := generateHDKeyPair(t)
	data := []byte("Paul Atreides")
	sig := []byte{0x2f, 0xfd, 0x9c, 0x1a, 0x5c, 0x28, 0x0, 0x7e, 0xb5, 0xfe, 0x2f, 0xbf, 0x7b, 0xe4, 0x46, 0xcf, 0x0, 0xd6, 0xed, 0xee, 0x21, 0x31, 0xa6, 0x58, 0xf4, 0xa0, 0x42, 0x4b, 0x7f, 0xc4, 0xcd, 0x8e, 0xf6, 0xa3, 0x23, 0x7a, 0xa, 0x9d, 0x3, 0x55, 0xe8, 0xe, 0xab, 0xb2, 0xdd, 0x27, 0x16, 0x63, 0x8a, 0x5c, 0x54, 0x5a, 0x3b, 0x9a, 0x2c, 0xa4, 0xa6, 0xc5, 0xd2, 0x68, 0x98, 0x7, 0x5, 0x1}

	// when
	verified, err := kp.VerifyAny(data, sig)

	// then
	require.NoError(t, err)
	assert.True(t, verified)
}

func testHDKeyPairVerifyingAnyMessageWithInvalidSignatureFails(t *testing.T) {
	// given
	kp := generateHDKeyPair(t)
	data := []byte("Paul Atreides")
	sig := []byte("Vladimir Harkonnen")

	// when
	verified, err := kp.VerifyAny(data, sig)

	// then
	require.NoError(t, err)
	assert.False(t, verified)
}

func testHDKeyPairMarshalingKeyPairSucceeds(t *testing.T) {
	// given
	kp := generateHDKeyPair(t)

	// when
	m, err := json.Marshal(kp)

	// then
	assert.NoError(t, err)
	expected := fmt.Sprintf(`{"index":1,"public_key":"%s","private_key":"%s","meta":null,"tainted":false,"algorithm":{"name":"vega/ed25519","version":1}}`, testPublicKey, testPrivateKey)
	assert.Equal(t, expected, string(m))
}

func testHDKeyPairUnmarshalingKeyPairSucceeds(t *testing.T) {
	// given
	kp := wallet.HDKeyPair{}
	marshalled := fmt.Sprintf(`{"index":1,"public_key":"%s","private_key":"%s","meta":null,"tainted":false,"algorithm":{"name":"vega/ed25519","version":1}}`, testPublicKey, testPrivateKey)

	// when
	err := json.Unmarshal([]byte(marshalled), &kp)

	// then
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), kp.Index())
	assert.Equal(t, testPublicKey, kp.PublicKey())
	assert.Equal(t, testPrivateKey, kp.PrivateKey())
	assert.Equal(t, uint32(1), kp.AlgorithmVersion())
	assert.Equal(t, "vega/ed25519", kp.AlgorithmName())
	assert.False(t, kp.IsTainted())
	assert.Nil(t, kp.Meta())
}

func generateHDKeyPair(t *testing.T) *wallet.HDKeyPair {
	t.Helper()

	rawPublicKey, err := hex.DecodeString(testPublicKey)
	if err != nil {
		t.Fatalf("couldn't decode public key: %v", err)
	}
	rawPrivateKey, err := hex.DecodeString(testPrivateKey)
	if err != nil {
		t.Fatalf("couldn't decode private key: %v", err)
	}

	kp, err := wallet.NewHDKeyPair(1, rawPublicKey, rawPrivateKey)
	if err != nil {
		t.Fatalf("couldn't create HD key pair: %v", err)
	}

	return kp
}
