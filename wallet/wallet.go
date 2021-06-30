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

func (w *Wallet) TaintKey(pubKey string) error {
	keyPair, err := w.KeyRing.FindPair(pubKey)
	if err != nil {
		return err
	}

	if err := keyPair.Taint(); err != nil {
		return err
	}

	w.KeyRing.Upsert(keyPair)

	return nil
}

type Meta struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
