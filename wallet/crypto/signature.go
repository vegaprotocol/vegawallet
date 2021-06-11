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
	GenKey() (crypto.PublicKey, crypto.PrivateKey, error)
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

func NewSignatureAlgorithm(algo string) (SignatureAlgorithm, error) {
	switch algo {
	case Ed25519:
		return NewEd25519(), nil
	default:
		return SignatureAlgorithm{}, ErrUnsupportedSignatureAlgorithm
	}

}

func (s *SignatureAlgorithm) GenKey() (crypto.PublicKey, crypto.PrivateKey, error) {
	return s.impl.GenKey()
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

func (s *SignatureAlgorithm) MarshalJSON() ([]byte, error) {
	if s != nil {
		return json.Marshal(s.Name())
	}
	return nil, errors.New("nil signature")
}

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
