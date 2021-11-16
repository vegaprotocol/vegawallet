package wallet_test

import (
	"encoding/hex"
	"testing"

	"code.vegaprotocol.io/vegawallet/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHDMasterKeypair(t *testing.T) {
	t.Run("New master key pair succeeds", testHDMasterKeyPairNewKeyPairSucceeds)
	t.Run("Signing transaction succeeds", testHDMasterKeyPairSigningAnyMessageSucceeds)
	t.Run("Signing any message succeeds", testHDMasterKeyPairSigningTransactionSucceeds)
}

func testHDMasterKeyPairNewKeyPairSucceeds(t *testing.T) {
	// given
	publicKey := "e20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a"
	rawPublicKey, err := hex.DecodeString(publicKey)
	if err != nil {
		panic(err)
	}
	privateKey := "c55ae95741a0431186bb062501ad9c5a31505fc14d9fd91f08108b11fbec33d9e20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a"
	rawPrivateKey, err := hex.DecodeString(privateKey)
	if err != nil {
		panic(err)
	}

	// when
	kp, err := wallet.NewHDMasterKeyPair(rawPublicKey, rawPrivateKey)

	// then
	require.NoError(t, err)
	assert.NotNil(t, kp)
	assert.Equal(t, publicKey, kp.PublicKey())
	assert.Equal(t, privateKey, kp.PrivateKey())
	assert.Equal(t, "vega/ed25519", kp.AlgorithmName())
	assert.Equal(t, uint32(1), kp.AlgorithmVersion())
}

func testHDMasterKeyPairSigningTransactionSucceeds(t *testing.T) {
	// given
	kp := getHDMasterKeyPair(t)
	data := []byte("Paul Atreides")

	// when
	sig, err := kp.Sign(data)

	// then
	require.NoError(t, err)
	assert.Equal(t, "2ffd9c1a5c28007eb5fe2fbf7be446cf00d6edee2131a658f4a0424b7fc4cd8ef6a3237a0a9d0355e80eabb2dd2716638a5c545a3b9a2ca4a6c5d26898070501", sig.Value)
	assert.Equal(t, kp.AlgorithmName(), sig.Algo)
	assert.Equal(t, kp.AlgorithmVersion(), sig.Version)
}

func testHDMasterKeyPairSigningAnyMessageSucceeds(t *testing.T) {
	// given
	kp := getHDMasterKeyPair(t)
	data := []byte("Paul Atreides")

	// when
	sig, err := kp.SignAny(data)

	// then
	require.NoError(t, err)
	assert.Equal(t, []byte{0x2f, 0xfd, 0x9c, 0x1a, 0x5c, 0x28, 0x0, 0x7e, 0xb5, 0xfe, 0x2f, 0xbf, 0x7b, 0xe4, 0x46, 0xcf, 0x0, 0xd6, 0xed, 0xee, 0x21, 0x31, 0xa6, 0x58, 0xf4, 0xa0, 0x42, 0x4b, 0x7f, 0xc4, 0xcd, 0x8e, 0xf6, 0xa3, 0x23, 0x7a, 0xa, 0x9d, 0x3, 0x55, 0xe8, 0xe, 0xab, 0xb2, 0xdd, 0x27, 0x16, 0x63, 0x8a, 0x5c, 0x54, 0x5a, 0x3b, 0x9a, 0x2c, 0xa4, 0xa6, 0xc5, 0xd2, 0x68, 0x98, 0x7, 0x5, 0x1}, sig)
}

func getHDMasterKeyPair(t *testing.T) *wallet.HDMasterKeyPair {
	t.Helper()

	publicKey := "e20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a"
	rawPublicKey, err := hex.DecodeString(publicKey)
	assert.NoError(t, err)

	privateKey := "c55ae95741a0431186bb062501ad9c5a31505fc14d9fd91f08108b11fbec33d9e20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a"
	rawPrivateKey, err := hex.DecodeString(privateKey)
	assert.NoError(t, err)

	kp, err := wallet.NewHDMasterKeyPair(rawPublicKey, rawPrivateKey)
	assert.NoError(t, err)

	return kp
}
