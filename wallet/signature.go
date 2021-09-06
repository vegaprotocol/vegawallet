package wallet

type Signature struct {
	// Value is hex-encoded
	Value   string `json:"value"`
	Algo    string `json:"algo"`
	Version uint32 `json:"version"`
}
