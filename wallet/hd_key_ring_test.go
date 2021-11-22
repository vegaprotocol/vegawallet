package wallet_test

import (
	"encoding/hex"
	"testing"

	"code.vegaprotocol.io/vegawallet/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHDKeyRing(t *testing.T) {
	t.Run("New key ring succeeds", testHDKeyRingNewKeyRingSucceeds)
	t.Run("Loading key ring succeeds", testHDKeyRingLoadingKeyRingSucceeds)
	t.Run("Adding a new key succeeds", testHDKeyRingAddingNewKeySucceeds)
	t.Run("Updating an existing key succeeds", testHDKeyRingUpdatingExistingKeySucceeds)
	t.Run("Getting public keys succeeds", testHDKeyRingGettingPublicKeysSucceeds)
	t.Run("Getting key pairs succeeds", testHDKeyRingGettingKeyPairsSucceeds)
	t.Run("Finding an existing key pair succeeds", testHDKeyRingFindingExistingKeyPairSucceeds)
	t.Run("Finding a non-existing key pair fails", testHDKeyRingFindingNonExistingKeyPairFails)
	t.Run("Incrementing next index succeeds", testIncrementingNextIndexSucceeds)
}

func testHDKeyRingNewKeyRingSucceeds(t *testing.T) {
	// when
	ring := wallet.NewHDKeyRing()

	// then
	assert.NotNil(t, ring)
	assert.Empty(t, ring.ListKeyPairs())
	assert.Equal(t, uint32(1), ring.NextIndex())
}

func testHDKeyRingLoadingKeyRingSucceeds(t *testing.T) {
	// given
	kp1 := bundleHDKeyPair(1,
		"e20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a",
		"c55ae95741a0431186bb062501ad9c5a31505fc14d9fd91f08108b11fbec33d9e20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a",
	)
	kp2 := bundleHDKeyPair(2,
		"a20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a",
		"a55ae95741a0431186bb062501ad9c5a31505fc14d9fd91f08108b11fbec33d9e20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a",
	)

	// when
	ring := wallet.LoadHDKeyRing([]wallet.HDKeyPair{*kp1, *kp2})

	// then
	assert.NotNil(t, ring)
	keyPairs := ring.ListKeyPairs()
	assert.Equal(t, keyPairs, []wallet.HDKeyPair{*kp1, *kp2})
	assert.Equal(t, uint32(3), ring.NextIndex())
}

func testHDKeyRingAddingNewKeySucceeds(t *testing.T) {
	// given
	kp := generateHDKeyPair(t)
	ring := wallet.NewHDKeyRing()

	// when
	ring.Upsert(*kp)

	// then
	keyPairs := ring.ListKeyPairs()
	assert.Contains(t, keyPairs, *kp)
	assert.Equal(t, kp.Index()+1, ring.NextIndex())
}

func testHDKeyRingUpdatingExistingKeySucceeds(t *testing.T) {
	// given
	kp := generateHDKeyPair(t)
	ring := wallet.NewHDKeyRing()

	// when
	ring.Upsert(*kp)

	// given
	updatedKp := *kp

	// when
	err := updatedKp.Taint()

	// then
	require.NoError(t, err)

	// when
	ring.Upsert(updatedKp)

	// then
	keyPairs := ring.ListKeyPairs()
	assert.Contains(t, keyPairs, updatedKp)
	assert.NotContains(t, keyPairs, *kp)
}

func testHDKeyRingGettingPublicKeysSucceeds(t *testing.T) {
	// given
	kp1 := bundleHDKeyPair(1,
		"e20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a",
		"c55ae95741a0431186bb062501ad9c5a31505fc14d9fd91f08108b11fbec33d9e20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a",
	)
	kp2 := bundleHDKeyPair(2,
		"a20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a",
		"a55ae95741a0431186bb062501ad9c5a31505fc14d9fd91f08108b11fbec33d9e20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a",
	)
	ring := wallet.NewHDKeyRing()

	// setup
	ring.Upsert(*kp2)
	ring.Upsert(*kp1)

	// when
	keys := ring.ListPublicKeys()

	// then
	assert.Len(t, keys, 2)
	assert.Equal(t, keys, []wallet.HDPublicKey{
		kp1.ToPublicKey(),
		kp2.ToPublicKey(),
	})
}

func testHDKeyRingGettingKeyPairsSucceeds(t *testing.T) {
	// given
	kp1 := bundleHDKeyPair(1,
		"e20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a",
		"c55ae95741a0431186bb062501ad9c5a31505fc14d9fd91f08108b11fbec33d9e20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a",
	)
	kp2 := bundleHDKeyPair(2,
		"a20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a",
		"a55ae95741a0431186bb062501ad9c5a31505fc14d9fd91f08108b11fbec33d9e20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a",
	)
	ring := wallet.NewHDKeyRing()

	// setup
	ring.Upsert(*kp2)
	ring.Upsert(*kp1)

	// when
	keys := ring.ListKeyPairs()

	// then
	assert.Len(t, keys, 2)
	assert.Equal(t, []wallet.HDKeyPair{*kp1, *kp2}, keys)
}

func testHDKeyRingFindingExistingKeyPairSucceeds(t *testing.T) {
	// given
	kp1 := generateHDKeyPair(t)
	kp2 := generateHDKeyPair(t)
	ring := wallet.NewHDKeyRing()

	// setup
	ring.Upsert(*kp1)
	ring.Upsert(*kp2)

	// when
	returnedKp, found := ring.FindPair(kp2.PublicKey())

	// then
	require.True(t, found)
	assert.Equal(t, *kp2, returnedKp)
}

func testHDKeyRingFindingNonExistingKeyPairFails(t *testing.T) {
	// given
	kp1 := generateHDKeyPair(t)
	kp2 := generateHDKeyPair(t)
	ring := wallet.NewHDKeyRing()

	// setup
	ring.Upsert(*kp1)
	ring.Upsert(*kp2)

	// when
	returnedKp, found := ring.FindPair("non-existing public key")

	// then
	assert.False(t, found)
	assert.Equal(t, wallet.HDKeyPair{}, returnedKp)
}

func testIncrementingNextIndexSucceeds(t *testing.T) {
	// given
	kp1 := bundleHDKeyPair(1,
		"e20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a",
		"c55ae95741a0431186bb062501ad9c5a31505fc14d9fd91f08108b11fbec33d9e20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a",
	)
	kp2 := bundleHDKeyPair(2,
		"a20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a",
		"a55ae95741a0431186bb062501ad9c5a31505fc14d9fd91f08108b11fbec33d9e20fcd0aa4cc2ea7c18d55dc083d14006377295418f657ba52593eae14e3a98a",
	)
	ring := wallet.NewHDKeyRing()

	// setup
	ring.Upsert(*kp2)
	ring.Upsert(*kp1)

	// when
	nextIndex := ring.NextIndex()

	// then
	assert.Equal(t, uint32(3), nextIndex)
}

func bundleHDKeyPair(index uint32, publicKey, privateKey string) *wallet.HDKeyPair {
	rawPublicKey, err := hex.DecodeString(publicKey)
	if err != nil {
		panic(err)
	}
	rawPrivateKey, err := hex.DecodeString(privateKey)
	if err != nil {
		panic(err)
	}

	kp, err := wallet.NewHDKeyPair(index, rawPublicKey, rawPrivateKey)
	if err != nil {
		panic(err)
	}

	return kp
}
