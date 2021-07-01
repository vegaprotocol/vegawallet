package wallet

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

func (r LegacyKeyRing) GetPublicKeys() []*LegacyPublicKey {
	pubKeys := make([]*LegacyPublicKey, 0, len(r))
	for _, keyPair := range r {
		pubKeys = append(pubKeys, keyPair.ToPublicKey())
	}
	return pubKeys
}
