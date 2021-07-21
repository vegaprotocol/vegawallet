package wallet_test

import (
	"testing"

	"code.vegaprotocol.io/go-wallet/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLegacyKeyRing(t *testing.T) {
	t.Run("Adding a new key succeeds", testLegacyKeyRingAddingNewKeySucceeds)
	t.Run("Updating an existing key succeeds", testLegacyKeyRingUpdatingExistingKeySucceeds)
	t.Run("Getting public keys succeeds", testLegacyKeyRingGettingPublicKeysSucceeds)
	t.Run("Getting key pairs succeeds", testLegacyKeyRingGettingKeyPairsSucceeds)
	t.Run("Finding an existing key pair succeeds", testLegacyKeyRingFindingExistingKeyPairSucceeds)
	t.Run("Finding a non-existing key pair fails", testLegacyKeyRingFindingNonExistingKeyPairFails)
}

func testLegacyKeyRingAddingNewKeySucceeds(t *testing.T) {
	// given
	kp := generateLegacyKeyPair()
	ring := wallet.NewLegacyKeyRing()

	// when
	ring.Upsert(*kp)

	// then
	assert.Contains(t, ring, *kp)
}

func testLegacyKeyRingUpdatingExistingKeySucceeds(t *testing.T) {
	// given
	kp := generateLegacyKeyPair()
	ring := wallet.NewLegacyKeyRing()

	// when
	ring.Upsert(*kp)

	// given
	updatedKp := *kp
	updatedKp.Tainted = true

	// when
	ring.Upsert(updatedKp)

	// then
	assert.Contains(t, ring, updatedKp)
	assert.NotContains(t, ring, *kp)
}

func testLegacyKeyRingGettingPublicKeysSucceeds(t *testing.T) {
	// given
	kp1 := wallet.LegacyKeyPair{
		Pub:  "bbbbbb",
		Priv: "111111",
	}
	kp2 := wallet.LegacyKeyPair{
		Pub:  "aaaaaa",
		Priv: "222222",
	}
	ring := wallet.NewLegacyKeyRing()

	// setup
	ring.Upsert(kp1)
	ring.Upsert(kp2)

	// when
	keys := ring.GetPublicKeys()

	// then
	assert.Len(t, keys, 2)
	assert.Equal(t, keys, []wallet.LegacyPublicKey{
		*kp2.ToPublicKey(),
		*kp1.ToPublicKey(),
	})
}

func testLegacyKeyRingGettingKeyPairsSucceeds(t *testing.T) {
	// given
	kp1 := wallet.LegacyKeyPair{
		Pub:  "bbbbbb",
		Priv: "111111",
	}
	kp2 := wallet.LegacyKeyPair{
		Pub:  "aaaaaa",
		Priv: "222222",
	}
	ring := wallet.NewLegacyKeyRing()

	// setup
	ring.Upsert(kp1)
	ring.Upsert(kp2)

	// when
	keys := ring.GetKeyPairs()

	// then
	assert.Len(t, keys, 2)
	assert.Equal(t, keys, []wallet.LegacyKeyPair{kp2, kp1})
}

func testLegacyKeyRingFindingExistingKeyPairSucceeds(t *testing.T) {
	// given
	kp1 := generateLegacyKeyPair()
	kp2 := generateLegacyKeyPair()
	ring := wallet.NewLegacyKeyRing()

	// setup
	ring.Upsert(*kp1)
	ring.Upsert(*kp2)

	// when
	returnedKp, err := ring.FindPair(kp2.Pub)

	// then
	require.NoError(t, err)
	assert.Equal(t, *kp2, returnedKp)
}

func testLegacyKeyRingFindingNonExistingKeyPairFails(t *testing.T) {
	// given
	kp1 := generateLegacyKeyPair()
	kp2 := generateLegacyKeyPair()
	ring := wallet.NewLegacyKeyRing()

	// setup
	ring.Upsert(*kp1)
	ring.Upsert(*kp2)

	// when
	returnedKp, err := ring.FindPair("non-existing public key")

	// then
	assert.Error(t, err)
	assert.Equal(t, wallet.LegacyKeyPair{}, returnedKp)
}
