package wallet

type Wallet struct {
	Owner   string
	KeyRing KeyRing `json:"Keypairs"`
}

func NewWallet(name string) *Wallet {
	return &Wallet{
		Owner: name,
	}
}

type Meta struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func New(owner string) Wallet {
	return Wallet{
		Owner:   owner,
		KeyRing: NewKeyRing(),
	}
}


