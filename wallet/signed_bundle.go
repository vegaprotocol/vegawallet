package wallet

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"

	types "code.vegaprotocol.io/go-wallet/proto"
)

// here we implement Marhsalling for the SignedBundle
// internally the signed bundle is a bunch of bytes slices
// when being returned to the client they are strings.
// Data and sig become base64 encoded string, as they
// are literaly blobs of binary data.
// the Public key will be encoded into hex format
// as it's a format commonly used by address in the blockchain
// world, and also is easier to read etc.

type SignedBundle struct {
	Data   []byte `json:"data"`
	Sig    []byte `json:"sig"`
	PubKey []byte `json:"pubKey"`
}

type signedBundleStrings struct {
	Data   string `json:"data"`
	Sig    string `json:"sig"`
	PubKey string `json:"pubKey"`
}

func (s SignedBundle) MarshalJSON() ([]byte, error) {
	stringBundle := signedBundleStrings{
		Data:   base64.StdEncoding.EncodeToString(s.Data),
		Sig:    base64.StdEncoding.EncodeToString(s.Sig),
		PubKey: hex.EncodeToString(s.PubKey),
	}
	return json.Marshal(stringBundle)
}

func (s *SignedBundle) IntoProto() *types.SignedBundle {
	return &types.SignedBundle{
		Data: s.Data,
		Sig:  s.Sig,
		Auth: &types.SignedBundle_PubKey{
			PubKey: s.PubKey,
		},
	}

}
