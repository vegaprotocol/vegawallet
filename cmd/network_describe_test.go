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

func TestDescribeNetworkFlags(t *testing.T) {
	t.Run("Valid flags with Network succeeds", testDescribeNetworkValidFlagsWithNetworkSucceeds)
	t.Run("Invalid flags without Network fails", testDescribeNetworkInvalidFlagsWithoutNetworkFails)
}

func testDescribeNetworkValidFlagsWithNetworkSucceeds(t *testing.T) {
	// given
	networkName := vgrand.RandomStr(10)

	f := &cmd.DescribeNetworkFlags{
		Network: networkName,
	}

	expectedReq := &network.DescribeNetworkRequest{
		Name: networkName,
	}
	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	require.NotNil(t, req)
	assert.Equal(t, expectedReq, req)
}

func testDescribeNetworkInvalidFlagsWithoutNetworkFails(t *testing.T) {
	// given
	f := &cmd.DescribeNetworkFlags{}

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("network"))
	require.Nil(t, req)
}
