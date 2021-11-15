package cmd_test

import (
	"encoding/base64"
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"code.vegaprotocol.io/vegawallet/cmd"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyMessageFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testVerifyMessageFlagsValidFlagsSucceeds)
	t.Run("Missing public key fails", testVerifyMessageFlagsMissingPubKeyFails)
	t.Run("Missing message fails", testVerifyMessageFlagsMissingMessageFails)
	t.Run("Malformed message fails", testVerifyMessageFlagsMalformedMessageFails)
	t.Run("Missing signature fails", testVerifyMessageFlagsMissingSignatureFails)
	t.Run("Malformed signature fails", testVerifyMessageFlagsMalformedSignatureFails)
}

func testVerifyMessageFlagsValidFlagsSucceeds(t *testing.T) {
	// given
	pubKey := vgrand.RandomStr(20)
	decodedMessage := []byte(vgrand.RandomStr(20))
	decodedSignature := []byte(vgrand.RandomStr(20))

	f := &cmd.VerifyMessageFlags{
		PubKey:    pubKey,
		Message:   base64.StdEncoding.EncodeToString(decodedMessage),
		Signature: base64.StdEncoding.EncodeToString(decodedSignature),
	}

	expectedReq := &crypto.VerifyMessageRequest{
		PubKey:    pubKey,
		Signature: decodedSignature,
		Message:   decodedMessage,
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, expectedReq, req)
}

func testVerifyMessageFlagsMissingPubKeyFails(t *testing.T) {
	// given
	f := newVerifyMessageFlags(t)
	f.PubKey = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("pubkey"))
	assert.Nil(t, req)
}

func testVerifyMessageFlagsMissingMessageFails(t *testing.T) {
	// given
	f := newVerifyMessageFlags(t)
	f.Message = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("message"))
	assert.Nil(t, req)
}

func testVerifyMessageFlagsMalformedMessageFails(t *testing.T) {
	// given
	f := newVerifyMessageFlags(t)
	f.Message = "not-base-64"

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.MustBase64EncodedError("message"))
	assert.Nil(t, req)
}

func testVerifyMessageFlagsMissingSignatureFails(t *testing.T) {
	// given
	f := newVerifyMessageFlags(t)
	f.Signature = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("signature"))
	assert.Nil(t, req)
}

func testVerifyMessageFlagsMalformedSignatureFails(t *testing.T) {
	// given
	f := newVerifyMessageFlags(t)
	f.Signature = "not-base-64"

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.MustBase64EncodedError("signature"))
	assert.Nil(t, req)
}

func newVerifyMessageFlags(t *testing.T) *cmd.VerifyMessageFlags {
	t.Helper()

	pubKey := vgrand.RandomStr(20)
	decodedMessage := []byte(vgrand.RandomStr(20))
	decodedSignature := []byte(vgrand.RandomStr(20))

	return &cmd.VerifyMessageFlags{
		PubKey:    pubKey,
		Message:   base64.StdEncoding.EncodeToString(decodedMessage),
		Signature: base64.StdEncoding.EncodeToString(decodedSignature),
	}
}
