package wallet

import (
	"encoding/json"
)

type HDPublicKey struct {
	Idx       uint32        `json:"index"`
	PublicKey string        `json:"pub"`
	Algorithm jsonAlgorithm `json:"algorithm"`
	Tainted   bool          `json:"tainted"`
	MetaList  []Meta        `json:"meta"`
}

func (k *HDPublicKey) Index() uint32 {
	return k.Idx
}

func (k *HDPublicKey) Key() string {
	return k.PublicKey
}

func (k *HDPublicKey) IsTainted() bool {
	return k.Tainted
}

func (k *HDPublicKey) Meta() []Meta {
	return k.MetaList
}

func (k *HDPublicKey) AlgorithmVersion() uint32 {
	return k.Algorithm.Version
}

func (k *HDPublicKey) AlgorithmName() string {
	return k.Algorithm.Name
}

func (k *HDPublicKey) MarshalJSON() ([]byte, error) {
	type alias HDPublicKey
	aliasPublicKey := (*alias)(k)
	return json.Marshal(aliasPublicKey)
}

func (k *HDPublicKey) UnmarshalJSON(data []byte) error {
	type alias HDPublicKey
	aliasPublicKey := (*alias)(k)
	return json.Unmarshal(data, aliasPublicKey)
}
