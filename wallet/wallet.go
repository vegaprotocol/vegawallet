package wallet

type Wallet struct {
	Name    string  `json:"Owner"`
	KeyRing KeyRing `json:"Keypairs"`
}

func NewWallet(name string) *Wallet {
	return &Wallet{
		Name:    name,
		KeyRing: NewKeyRing(),
	}
}

type Meta struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
