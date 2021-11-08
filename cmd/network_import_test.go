package cmd_test

import (
	"testing"

	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"code.vegaprotocol.io/vegawallet/cmd"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportNetworkFlags(t *testing.T) {
	t.Run("Valid flags with URL succeeds", testImportNetworkFlagsValidFlagsWithURLSucceeds)
	t.Run("Valid flags with file path succeeds", testImportNetworkFlagsValidFlagsWithFilePathSucceeds)
	t.Run("Missing URL and file path fails", testImportNetworkFlagsMissingURLAndFilePathFails)
	t.Run("Missing public key fails", testImportNetworkFlagsBothURLAndFilePathSpecifiedFails)
}

func testImportNetworkFlagsValidFlagsWithURLSucceeds(t *testing.T) {
	// given
	networkName := vgrand.RandomStr(10)
	url := vgrand.RandomStr(20)

	f := &cmd.ImportNetworkFlags{
		Name:  networkName,
		URL:   url,
		Force: true,
	}

	expectedReq := &network.ImportNetworkFromSourceRequest{
		Name:  networkName,
		URL:   url,
		Force: true,
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, expectedReq, req)
}

func testImportNetworkFlagsValidFlagsWithFilePathSucceeds(t *testing.T) {
	// given
	networkName := vgrand.RandomStr(10)
	filePath := vgrand.RandomStr(20)

	f := &cmd.ImportNetworkFlags{
		Name:     networkName,
		FilePath: filePath,
		Force:    true,
	}

	expectedReq := &network.ImportNetworkFromSourceRequest{
		Name:     networkName,
		FilePath: filePath,
		Force:    true,
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, expectedReq, req)
}

func testImportNetworkFlagsMissingURLAndFilePathFails(t *testing.T) {
	// given
	f := newImportNetworkFlags(t)
	f.URL = ""
	f.FilePath = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.OneOfFlagsMustBeSpecifiedError("from-file", "from-url"))
	assert.Nil(t, req)
}

func testImportNetworkFlagsBothURLAndFilePathSpecifiedFails(t *testing.T) {
	// given
	f := newImportNetworkFlags(t)
	f.URL = vgrand.RandomStr(20)
	f.FilePath = vgrand.RandomStr(20)

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagsMutuallyExclusiveError("from-file", "from-url"))
	assert.Nil(t, req)
}

func newImportNetworkFlags(t *testing.T) *cmd.ImportNetworkFlags {
	t.Helper()

	networkName := vgrand.RandomStr(10)

	return &cmd.ImportNetworkFlags{
		Name:  networkName,
		Force: true,
	}
}
