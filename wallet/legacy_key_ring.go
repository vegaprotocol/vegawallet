package wallet

import "sort"

type LegacyKeyRing []LegacyKeyPair

func NewLegacyKeyRing() LegacyKeyRing {
	return []LegacyKeyPair{}
}

func (r LegacyKeyRing) FindPair(pubKey string) (LegacyKeyPair, error) {
	for i := range r {
		if r[i].Pub == pubKey {
			return r[i], nil
		}
	}
	return LegacyKeyPair{}, ErrPubKeyDoesNotExist
}

func (r *LegacyKeyRing) Upsert(pair LegacyKeyPair) {
	for i := range *r {
		if (*r)[i].Pub == pair.Pub && (*r)[i].Priv == pair.Priv {
			(*r)[i] = pair
			return
		}
	}

	*r = append(*r, pair)
}

func (r LegacyKeyRing) GetPublicKeys() []LegacyPublicKey {
	keyPairs := r.GetKeyPairs()
	pubKeys := make([]LegacyPublicKey, len(keyPairs))
	for i, keyPair := range keyPairs {
		pubKeys[i] = *keyPair.ToPublicKey()
	}
	return pubKeys
}

func (r LegacyKeyRing) GetKeyPairs() []LegacyKeyPair {
	keysList := make([]LegacyKeyPair, len(r))
	for i, key := range r {
		keysList[i] = key
	}
	sort.Sort(byPubKey(keysList))
	return keysList
}

type byPubKey []LegacyKeyPair

func (a byPubKey) Len() int {
	return len(a)
}

func (a byPubKey) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a byPubKey) Less(i, j int) bool {
	return a[i].Pub < a[j].Pub
}
