package crypto

import (
	"encoding/hex"
	"fmt"
)

type VerifyMessageRequest struct {
	Message   []byte
	Signature []byte
	PubKey    string
}

func VerifyMessage(req *VerifyMessageRequest) (bool, error) {
	decodedPubKey, err := hex.DecodeString(req.PubKey)
	if err != nil {
		return false, fmt.Errorf("couldn't decode public key: %w", err)
	}

	signatureAlgorithm := NewEd25519()
	return signatureAlgorithm.Verify(decodedPubKey, req.Message, req.Signature)
}
