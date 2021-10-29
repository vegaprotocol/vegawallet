package wallet

const KeyNameMeta = "name"

type Meta struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func GetKeyName(meta []Meta) string {
	for _, m := range meta {
		if m.Key == KeyNameMeta {
			return m.Value
		}
	}

	return "<No name>"
}
