package cmd_test

import (
	"encoding/json"
	"testing"

	"code.vegaprotocol.io/protos/vega"
	v1 "code.vegaprotocol.io/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/protos/vega/wallet/v1"
	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"code.vegaprotocol.io/vegawallet/cmd"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendCommandFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testSendCommandFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testSendCommandFlagsMissingWalletFails)
	t.Run("Missing log level fails", testSendCommandFlagsMissingLogLevelFails)
	t.Run("Unsupported log level fails", testSendCommandFlagsUnsupportedLogLevelFails)
	t.Run("Missing network and node address fails", testSendCommandFlagsMissingNetworkAndNodeAddressFails)
	t.Run("Both network and node address specified fails", testSendCommandFlagsBothNetworkAndNodeAddressSpecifiedFails)
	t.Run("Missing public key fails", testSendCommandFlagsMissingPubKeyFails)
	t.Run("Missing request fails", testSendCommandFlagsMissingRequestFails)
	t.Run("Malformed request fails", testSendCommandFlagsMalformedRequestFails)
	t.Run("Invalid request fails", testSendCommandFlagsInvalidRequestFails)
	t.Run("Request with public key set in it fails", testSendCommandFlagsRequestWithPubKeyFails)
}

func testSendCommandFlagsValidFlagsSucceeds(t *testing.T) {
	testDir := t.TempDir()

	// given
	passphrase, passphraseFilePath := NewPassphraseFile(t, testDir)
	network := vgrand.RandomStr(10)
	walletName := vgrand.RandomStr(10)
	pubKey := vgrand.RandomStr(20)

	f := &cmd.SendCommandFlags{
		Network:        network,
		NodeAddress:    "",
		Wallet:         walletName,
		PubKey:         pubKey,
		Retries:        10,
		LogLevel:       "debug",
		PassphraseFile: passphraseFilePath,
		RawCommand:     `{"voteSubmission": {"proposalId": "some-id", "value": "VALUE_YES"}}`,
	}

	expectedReq := &cmd.SendCommandRequest{
		Network:     network,
		NodeAddress: "",
		Wallet:      walletName,
		Retries:     10,
		LogLevel:    "debug",
		Passphrase:  passphrase,
		Request: &walletpb.SubmitTransactionRequest{
			PubKey:    pubKey,
			Propagate: true,
			Command: &walletpb.SubmitTransactionRequest_VoteSubmission{
				VoteSubmission: &v1.VoteSubmission{
					ProposalId: "some-id",
					Value:      vega.Vote_VALUE_YES,
				},
			},
		},
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	require.NotNil(t, req)

	expectedJson, _ := json.Marshal(expectedReq)
	actualJson, _ := json.Marshal(req)
	assert.Equal(t, expectedJson, actualJson)
}

func testSendCommandFlagsMissingWalletFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newSendCommandFlags(t, testDir)
	f.Wallet = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("wallet"))
	assert.Nil(t, req)
}

func testSendCommandFlagsMissingLogLevelFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newSendCommandFlags(t, testDir)
	f.LogLevel = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("level"))
	assert.Nil(t, req)
}

func testSendCommandFlagsUnsupportedLogLevelFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newSendCommandFlags(t, testDir)
	f.LogLevel = vgrand.RandomStr(5)

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, cmd.NewUnsupportedFlagValueError(f.LogLevel))
	assert.Nil(t, req)
}

func testSendCommandFlagsMissingNetworkAndNodeAddressFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newSendCommandFlags(t, testDir)
	f.Network = ""
	f.NodeAddress = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.OneOfFlagsMustBeSpecifiedError("network", "node-address"))
	assert.Nil(t, req)
}

func testSendCommandFlagsBothNetworkAndNodeAddressSpecifiedFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newSendCommandFlags(t, testDir)
	f.Network = vgrand.RandomStr(10)
	f.NodeAddress = vgrand.RandomStr(10)

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagsMutuallyExclusiveError("network", "node-address"))
	assert.Nil(t, req)
}

func testSendCommandFlagsMissingPubKeyFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newSendCommandFlags(t, testDir)
	f.PubKey = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("pubkey"))
	assert.Nil(t, req)
}

func testSendCommandFlagsMissingRequestFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newSendCommandFlags(t, testDir)
	f.RawCommand = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.ArgMustBeSpecifiedError("command"))
	assert.Nil(t, req)
}

func testSendCommandFlagsMalformedRequestFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newSendCommandFlags(t, testDir)
	f.RawCommand = vgrand.RandomStr(5)

	// when
	req, err := f.Validate()

	// then
	assert.Error(t, err)
	assert.Nil(t, req)
}

func testSendCommandFlagsInvalidRequestFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newSendCommandFlags(t, testDir)
	f.RawCommand = `{"voteSubmission": {}}`

	// when
	req, err := f.Validate()

	// then
	assert.Error(t, err)
	assert.Nil(t, req)
}

func testSendCommandFlagsRequestWithPubKeyFails(t *testing.T) {
	testDir := t.TempDir()

	// given
	f := newSendCommandFlags(t, testDir)
	f.RawCommand = `{"pubKey": "qwerty123456", "voteSubmission": {"proposalId": "some-id", "value": "VALUE_YES"}}`

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, cmd.ErrDoNotSetPubKeyInCommand)
	assert.Nil(t, req)
}

func newSendCommandFlags(t *testing.T, testDir string) *cmd.SendCommandFlags {
	t.Helper()

	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	networkName := vgrand.RandomStr(10)
	walletName := vgrand.RandomStr(10)
	pubKey := vgrand.RandomStr(20)

	return &cmd.SendCommandFlags{
		Network:        networkName,
		NodeAddress:    "",
		Retries:        10,
		LogLevel:       "debug",
		RawCommand:     `{"voteSubmission": {"proposalId": "some-id", "value": "VALUE_YES"}}`,
		Wallet:         walletName,
		PubKey:         pubKey,
		PassphraseFile: passphraseFilePath,
	}
}
