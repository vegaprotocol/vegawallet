package cmd_test

import (
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"code.vegaprotocol.io/vegawallet/cmd"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/wallet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateKeyFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testGenerateKeyFlagsValidFlagsSucceeds)
	t.Run("Missing wallet fails", testGenerateKeyFlagsMissingWalletFails)
	t.Run("Invalid metadata fails", testGenerateKeyFlagsInvalidMetadataFails)
}

func testGenerateKeyFlagsValidFlagsSucceeds(t *testing.T) {
	// given
	testDir := t.TempDir()
	walletName := vgrand.RandomStr(10)
	passphrase, passphraseFilePath := NewPassphraseFile(t, testDir)

	f := &cmd.GenerateKeyFlags{
		Wallet:         walletName,
		PassphraseFile: passphraseFilePath,
		RawMetadata:    []string{"name:my-wallet", "role:validation"},
	}

	expectedReq := &wallet.GenerateKeyRequest{
		Wallet: walletName,
		Metadata: []wallet.Meta{
			{Key: "name", Value: "my-wallet"},
			{Key: "role", Value: "validation"},
		},
		Passphrase: passphrase,
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, expectedReq, req)
}

func testGenerateKeyFlagsMissingWalletFails(t *testing.T) {
	// given
	f := newGenerateKeyFlags(t)
	f.Wallet = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("wallet"))
	assert.Nil(t, req)
}

func testGenerateKeyFlagsInvalidMetadataFails(t *testing.T) {
	// given
	f := newGenerateKeyFlags(t)
	f.RawMetadata = []string{"is=invalid"}

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.InvalidFlagFormatError("meta"))
	assert.Nil(t, req)
}

func newGenerateKeyFlags(t *testing.T) *cmd.GenerateKeyFlags {
	t.Helper()
	testDir := t.TempDir()
	_, passphraseFilePath := NewPassphraseFile(t, testDir)
	return &cmd.GenerateKeyFlags{
		Wallet:         vgrand.RandomStr(5),
		PassphraseFile: passphraseFilePath,
		RawMetadata:    []string{"name:my-wallet", "role:validation"},
	}
}
