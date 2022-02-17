package cmd_test

import (
	"encoding/json"
	"testing"

	commandspb "code.vegaprotocol.io/protos/vega/commands/v1"
	vgrand "code.vegaprotocol.io/shared/libs/rand"
	"code.vegaprotocol.io/vegawallet/cmd"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendTxFlags(t *testing.T) {
	t.Run("Valid flags succeeds", testSendTxFlagsValidFlagsSucceeds)
	t.Run("Missing log level fails", testSendTxFlagsMissingLogLevelFails)
	t.Run("Unsupported log level fails", testSendTxFlagsUnsupportedLogLevelFails)
	t.Run("Missing network and node address fails", testSendTxFlagsMissingNetworkAndNodeAddressFails)
	t.Run("Both network and node address specified fails", testSendTxFlagsBothNetworkAndNodeAddressSpecifiedFails)
	t.Run("Missing tx fails", testSendTxFlagsMissingTxFails)
	t.Run("Malformed tx fails", testSendTxFlagsMalformedTxFails)
}

func testSendTxFlagsValidFlagsSucceeds(t *testing.T) {
	// given
	network := vgrand.RandomStr(10)

	f := &cmd.SendTxFlags{
		Network:     network,
		NodeAddress: "",
		Retries:     10,
		LogLevel:    "debug",
		RawTx:       "ChwIxZXB58qn4K06EMC2BPI+CwoHc29tZS1pZBACEpMBCoABMTM1ZDdmN2Q4MjhkMjg3ZDMyNDQzYjQ2NGEyZDQwNTkyZjQ1OTgwMGQ0MGZmMzY5Y2VhMGFkZDUzZmZjNjYzYzlkZmU2YTI4MGIxZWI4MjdiOTJmYmY2NTY3NzI3MjgwYzMwODBiNjg5NGYyMjYzZmJlYmFkN2I2M2VhN2M4MGYSDHZlZ2EvZWQyNTUxORgBgH0B0j5AZjM4MTc5NjljZDMxNmQ1NmMzN2EzYzE5MjVjMDMyOWM5ZTMxMDQ0ODI5OGZmNzYyMjMwMTVjN2QyY2RiOTFiOQ==",
	}

	expectedReq := &cmd.SendTxRequest{
		Network:     network,
		NodeAddress: "",
		Retries:     10,
		LogLevel:    "debug",
		Tx: &commandspb.Transaction{
			InputData: []uint8{8, 197, 149, 193, 231, 202, 167, 224, 173, 58, 16, 192, 182, 4, 242, 62, 11, 10, 7, 115, 111, 109, 101, 45, 105, 100, 16, 2},
			Signature: &commandspb.Signature{
				Value:   "135d7f7d828d287d32443b464a2d40592f459800d40ff369cea0add53ffc663c9dfe6a280b1eb827b92fbf6567727280c3080b6894f2263fbebad7b63ea7c80f",
				Algo:    "vega/ed25519",
				Version: 1,
			},
			From: &commandspb.Transaction_PubKey{
				PubKey: "f3817969cd316d56c37a3c1925c0329c9e310448298ff76223015c7d2cdb91b9",
			},
			Version: 1,
		},
	}

	// when
	req, err := f.Validate()

	// then
	require.NoError(t, err)
	require.NotNil(t, req)
	expectedJSON, _ := json.Marshal(expectedReq)
	actualJSON, _ := json.Marshal(req)
	assert.Equal(t, expectedJSON, actualJSON)
}

func testSendTxFlagsMissingLogLevelFails(t *testing.T) {
	// given
	f := newSendTxFlags(t)
	f.LogLevel = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagMustBeSpecifiedError("level"))
	assert.Nil(t, req)
}

func testSendTxFlagsUnsupportedLogLevelFails(t *testing.T) {
	// given
	f := newSendTxFlags(t)
	f.LogLevel = vgrand.RandomStr(5)

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, cmd.NewUnsupportedFlagValueError(f.LogLevel))
	assert.Nil(t, req)
}

func testSendTxFlagsMissingNetworkAndNodeAddressFails(t *testing.T) {
	// given
	f := newSendTxFlags(t)
	f.Network = ""
	f.NodeAddress = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.OneOfFlagsMustBeSpecifiedError("network", "node-address"))
	assert.Nil(t, req)
}

func testSendTxFlagsBothNetworkAndNodeAddressSpecifiedFails(t *testing.T) {
	// given
	f := newSendTxFlags(t)
	f.Network = vgrand.RandomStr(10)
	f.NodeAddress = vgrand.RandomStr(10)

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.FlagsMutuallyExclusiveError("network", "node-address"))
	assert.Nil(t, req)
}

func testSendTxFlagsMissingTxFails(t *testing.T) {
	// given
	f := newSendTxFlags(t)
	f.RawTx = ""

	// when
	req, err := f.Validate()

	// then
	assert.ErrorIs(t, err, flags.ArgMustBeSpecifiedError("transaction"))
	assert.Nil(t, req)
}

func testSendTxFlagsMalformedTxFails(t *testing.T) {
	// given
	f := newSendTxFlags(t)
	f.RawTx = vgrand.RandomStr(5)

	// when
	req, err := f.Validate()

	// then
	assert.Error(t, err)
	assert.Nil(t, req)
}

func newSendTxFlags(t *testing.T) *cmd.SendTxFlags {
	t.Helper()

	networkName := vgrand.RandomStr(10)

	return &cmd.SendTxFlags{
		Network:     networkName,
		NodeAddress: "",
		Retries:     10,
		LogLevel:    "debug",
		RawTx:       "ChsItbycz7nhsO4/EPZ38j4LCgdzb21lLWlkEAISkwEKgAE4NjNjY2NhZGU5OTM5NTU5NWFmMmRkYjc4MTRiM2Q0NTE4NTllNDljNGRkZjUwYjRkZTJkOGUwNTBhY2U2YTQzOTM4OGJmMmFiN2E0N2NhZDM3MjQ3YWEwNzU1Yzk5NmMxZDJmMDY4MTI1YzY5NGVlNGNiMmU4ZWEyZmE2YmYwNRIMdmVnYS9lZDI1NTE5GAGAfQHSPkBmMzgxNzk2OWNkMzE2ZDU2YzM3YTNjMTkyNWMwMzI5YzllMzEwNDQ4Mjk4ZmY3NjIyMzAxNWM3ZDJjZGI5MWI5",
	}
}
