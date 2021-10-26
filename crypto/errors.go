package crypto

import (
	"errors"
)

var (
	// ErrBadED25519PrivateKeyLength is returned if a private key with incorrect length is supplied.
	ErrBadED25519PrivateKeyLength = errors.New("bad ed25519 private key length")

	// ErrBadED25519PublicKeyLength is returned if a public key with incorrect length is supplied.
	ErrBadED25519PublicKeyLength = errors.New("bad ed25519 public key length")

	ErrCouldNotCastPrivateKeyToBytes = errors.New("couldn't cast private key to bytes")
	ErrCouldNotCastPublicKeyToBytes  = errors.New("couldn't cast public key to bytes")

	ErrSignatureIsNil = errors.New("signature is nil")
)
