package wallet

import (
	"encoding/base64"
	"encoding/json"

	typespb "code.vegaprotocol.io/protos/vega"
)

// here we implement Marhsalling for the SignedBundle
// internally the signed bundle is a bunch of bytes slices
// when being returned to the client they are strings.
// Data and sig become base64 encoded string, as they
// are literaly blobs of binary data.
// the Public key will be encoded into hex format
// as it's a format commonly used by address in the blockchain
// world, and also is easier to read etc.

type Signature struct {
	Sig     []byte `json:"sig"`
	Algo    string `json:"algo"`
	Version uint32 `json:"version"`
}

type signatureMarshalled struct {
	Sig     string `json:"sig"`
	Algo    string `json:"algo"`
	Version uint32 `json:"version"`
}

type SignedBundle struct {
	Tx  []byte    `json:"tx"`
	Sig Signature `json:"sig"`
}

type signedBundleMarshalled struct {
	Tx  string              `json:"tx"`
	Sig signatureMarshalled `json:"sig"`
}

func (s SignedBundle) MarshalJSON() ([]byte, error) {
	stringBundle := signedBundleMarshalled{
		Tx: base64.StdEncoding.EncodeToString(s.Tx),
		Sig: signatureMarshalled{
			Sig:     base64.StdEncoding.EncodeToString(s.Sig.Sig),
			Algo:    s.Sig.Algo,
			Version: s.Sig.Version,
		},
	}
	return json.Marshal(stringBundle)
}

func (s *SignedBundle) IntoProto() *typespb.SignedBundle {
	return &typespb.SignedBundle{
		Tx: s.Tx,
		Sig: &typespb.Signature{
			Sig:     s.Sig.Sig,
			Algo:    s.Sig.Algo,
			Version: s.Sig.Version,
		},
	}

}
