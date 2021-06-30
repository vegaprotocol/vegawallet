package wallet

type Wallet struct {
	Owner   string
	KeyRing KeyRing `json:"Keypairs"`
}

func NewWallet(name string) *Wallet {
	return &Wallet{
		Owner: name,
		KeyRing: NewKeyRing(),
	}
}

type Meta struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}


