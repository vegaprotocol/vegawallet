package wallet

import (
	"sort"
)

type HDKeyRing struct {
	keys      map[string]HDKeyPair
	nextIndex uint32
}

func NewHDKeyRing() *HDKeyRing {
	return &HDKeyRing{
		keys:      map[string]HDKeyPair{},
		nextIndex: 1,
	}
}

func LoadHDKeyRing(keyPairs []HDKeyPair) *HDKeyRing {
	keyRing := NewHDKeyRing()
	for _, keyPair := range keyPairs {
		keyRing.Upsert(keyPair)
	}
	return keyRing
}

func (r *HDKeyRing) FindPair(pubKey string) (HDKeyPair, bool) {
	keyPair, ok := r.keys[pubKey]
	return keyPair, ok
}

func (r *HDKeyRing) Upsert(keyPair HDKeyPair) {
	r.keys[keyPair.PublicKey()] = keyPair
	if r.nextIndex <= keyPair.Index() {
		r.nextIndex = keyPair.Index() + 1
	}
}

// ListPublicKeys returns the list of public keys sorted by key index.
func (r *HDKeyRing) ListPublicKeys() []HDPublicKey {
	sortedKeyPairs := r.ListKeyPairs()
	pubKeys := make([]HDPublicKey, len(r.keys))
	for i, keyPair := range sortedKeyPairs {
		pubKeys[i] = keyPair.ToPublicKey()
	}
	return pubKeys
}

func (r *HDKeyRing) NextIndex() uint32 {
	return r.nextIndex
}

// ListKeyPairs returns the list of key pairs sorted by key index.
func (r *HDKeyRing) ListKeyPairs() []HDKeyPair {
	keysList := make([]HDKeyPair, len(r.keys))
	i := 0
	for _, key := range r.keys {
		keysList[i] = key
		i += 1
	}
	sort.SliceStable(keysList, func(i, j int) bool {
		return keysList[i].Index() < keysList[j].Index()
	})
	return keysList
}
