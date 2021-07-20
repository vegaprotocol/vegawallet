package wallet

import (
	"encoding/json"

	crypto2 "code.vegaprotocol.io/go-wallet/crypto"
)

type LegacyPublicKey struct {
	Pub       string                     `json:"pub"`
	Algorithm crypto2.SignatureAlgorithm `json:"algo"`
	Tainted   bool                       `json:"tainted"`
	MetaList  []Meta                     `json:"meta"`
}

func (k *LegacyPublicKey) Key() string {
	return k.Pub
}

func (k *LegacyPublicKey) IsTainted() bool {
	return k.Tainted
}

func (k *LegacyPublicKey) Meta() []Meta {
	return k.MetaList
}

func (k *LegacyPublicKey) AlgorithmVersion() uint32 {
	return k.Algorithm.Version()
}

func (k *LegacyPublicKey) AlgorithmName() string {
	return k.Algorithm.Name()
}

func (k *LegacyPublicKey) MarshalJSON() ([]byte, error) {
	type alias LegacyPublicKey
	aliasPublicKey := (*alias)(k)
	return json.Marshal(aliasPublicKey)
}

func (k *LegacyPublicKey) UnmarshalJSON(data []byte) error {
	type alias LegacyPublicKey
	aliasPublicKey := (*alias)(k)
	return json.Unmarshal(data, aliasPublicKey)
}
