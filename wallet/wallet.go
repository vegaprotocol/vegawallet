package wallet

type Wallet interface {
	Version() uint32
	Name() string
	SetName(newName string)
	ID() string
	Type() string
	DescribePublicKey(pubKey string) (PublicKey, error)
	DescribeKeyPair(pubKey string) (KeyPair, error)
	ListPublicKeys() []PublicKey
	ListKeyPairs() []KeyPair
	GetMasterKeyPair() (MasterKeyPair, error)
	GenerateKeyPair(meta []Meta) (KeyPair, error)
	TaintKey(pubKey string) error
	UntaintKey(pubKey string) error
	UpdateMeta(pubKey string, meta []Meta) error
	SignAny(pubKey string, data []byte) ([]byte, error)
	VerifyAny(pubKey string, data, sig []byte) (bool, error)
	SignTx(pubKey string, data []byte) (*Signature, error)
	IsolateWithKey(pubKey string) (Wallet, error)
}

type KeyPair interface {
	PublicKey() string
	PrivateKey() string
	IsTainted() bool
	Meta() []Meta
	Index() uint32
	AlgorithmVersion() uint32
	AlgorithmName() string
	SignAny(data []byte) ([]byte, error)
	VerifyAny(data, sig []byte) (bool, error)
	Sign(data []byte) (*Signature, error)
}

type PublicKey interface {
	Key() string
	IsTainted() bool
	Meta() []Meta
	AlgorithmVersion() uint32
	AlgorithmName() string

	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

type MasterKeyPair interface {
	PublicKey() string
	PrivateKey() string
	AlgorithmVersion() uint32
	AlgorithmName() string

	SignAny(data []byte) ([]byte, error)
	Sign(data []byte) (*Signature, error)
}
