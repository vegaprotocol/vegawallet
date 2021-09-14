package crypto

import (
	"crypto"
	"encoding/json"
	"errors"
)

const (
	Ed25519 string = "vega/ed25519"
)

var (
	ErrUnsupportedSignatureAlgorithm = errors.New("unsupported signature algorithm")
)

type SignatureAlgorithm struct {
	impl signatureAlgorithmImpl
}

type signatureAlgorithmImpl interface {
	Sign(priv crypto.PrivateKey, buf []byte) ([]byte, error)
	Verify(pub crypto.PublicKey, message, sig []byte) (bool, error)
	Name() string
	Version() uint32
}

func NewEd25519() SignatureAlgorithm {
	return SignatureAlgorithm{
		impl: newEd25519(),
	}
}

func NewSignatureAlgorithm(name string, version uint32) (SignatureAlgorithm, error) {
	if name == Ed25519 && version == 1 {
		return NewEd25519(), nil
	}
	return SignatureAlgorithm{}, ErrUnsupportedSignatureAlgorithm
}

func (s *SignatureAlgorithm) Sign(priv crypto.PrivateKey, buf []byte) ([]byte, error) {
	return s.impl.Sign(priv, buf)
}

func (s *SignatureAlgorithm) Verify(pub crypto.PublicKey, message, sig []byte) (bool, error) {
	return s.impl.Verify(pub, message, sig)
}

func (s *SignatureAlgorithm) Name() string {
	return s.impl.Name()
}

func (s *SignatureAlgorithm) Version() uint32 {
	return s.impl.Version()
}

// TODO We should rethink this to include the version, based on the
//		jsonAlgorithm implementation.
func (s *SignatureAlgorithm) MarshalJSON() ([]byte, error) {
	if s != nil {
		return json.Marshal(s.Name())
	}
	return nil, errors.New("nil signature")
}

// TODO We should rethink this to include the version, based on the
//		jsonAlgorithm implementation.
func (s *SignatureAlgorithm) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}

	switch name {
	case Ed25519:
		s.impl = newEd25519()
		return nil
	default:
		return ErrUnsupportedSignatureAlgorithm
	}
}
