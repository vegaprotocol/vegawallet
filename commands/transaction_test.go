package commands_test

import (
	"testing"

	"code.vegaprotocol.io/go-wallet/commands"
	"code.vegaprotocol.io/go-wallet/wallet/crypto"
	commandspb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/commands/v1"

	"github.com/stretchr/testify/assert"
)

func TestCheckTransaction(t *testing.T) {
	t.Run("Submitting valid transaction succeeds", testSubmittingValidTransactionSucceeds)
	t.Run("Submitting empty transaction fails", testSubmittingEmptyTransactionFails)
	t.Run("Submitting nil transaction fails", testSubmittingNilTransactionFails)
	t.Run("Submitting transaction without input data fails", testSubmittingTransactionWithoutInputDataFails)
	t.Run("Submitting transaction without signature fails", testSubmittingTransactionWithoutSignatureFails)
	t.Run("Submitting transaction without signature value fails", testSubmittingTransactionWithoutSignatureValueFails)
	t.Run("Submitting transaction without signature algo fails", testSubmittingTransactionWithoutSignatureAlgoFails)
	t.Run("Submitting transaction without from fails", testSubmittingTransactionWithoutFromFails)
	t.Run("Submitting transaction without public key fails", testSubmittingTransactionWithoutPubKeyFromFails)
	t.Run("Submitting transaction with unsupported algo fails", testSubmittingTransactionWithUnsupportedAlgoFails)
	t.Run("Submitting transaction with invalid encoding of bytes fails", testSubmittingTransactionWithInvalidEncodingOfValueFails)
	t.Run("Submitting transaction with invalid encoding of bytes fails", testSubmittingTransactionWithInvalidEncodingOfPubKeyFails)
	t.Run("Submitting transaction with invalid signature fails", testSubmittingTransactionWithInvalidSignatureFails)
}

func testSubmittingValidTransactionSucceeds(t *testing.T) {
	err := checkTransaction(newValidTransaction())

	assert.True(t, err.Empty())
}

func testSubmittingEmptyTransactionFails(t *testing.T) {
	err := checkTransaction(&commandspb.Transaction{})

	assert.Error(t, err)
}

func testSubmittingNilTransactionFails(t *testing.T) {
	err := checkTransaction(nil)

	assert.Contains(t, err.Get("tx"), commands.ErrIsRequired)
}

func testSubmittingTransactionWithoutInputDataFails(t *testing.T) {
	tx := newValidTransaction()
	tx.InputData = []byte{}

	err := checkTransaction(tx)

	assert.Contains(t, err.Get("tx.input_data"), commands.ErrIsRequired)
}

func testSubmittingTransactionWithoutSignatureFails(t *testing.T) {
	tx := newValidTransaction()
	tx.Signature = nil

	err := checkTransaction(tx)

	assert.Contains(t, err.Get("tx.signature"), commands.ErrIsRequired)
}

func testSubmittingTransactionWithoutSignatureValueFails(t *testing.T) {
	tx := newValidTransaction()
	tx.Signature.Value = ""

	err := checkTransaction(tx)

	assert.Contains(t, err.Get("tx.signature.value"), commands.ErrIsRequired)
}

func testSubmittingTransactionWithoutSignatureAlgoFails(t *testing.T) {
	tx := newValidTransaction()
	tx.Signature.Algo = ""

	err := checkTransaction(tx)

	assert.Contains(t, err.Get("tx.signature.algo"), commands.ErrIsRequired)
}

func testSubmittingTransactionWithoutFromFails(t *testing.T) {
	tx := newValidTransaction()
	tx.From = nil

	err := checkTransaction(tx)

	assert.Contains(t, err.Get("tx.from"), commands.ErrIsRequired)
}

func testSubmittingTransactionWithoutPubKeyFromFails(t *testing.T) {
	tx := newValidTransaction()
	tx.From = &commandspb.Transaction_PubKey{
		PubKey: "",
	}

	err := checkTransaction(tx)

	assert.Contains(t, err.Get("tx.from.pub_key"), commands.ErrIsRequired)
}

func testSubmittingTransactionWithUnsupportedAlgoFails(t *testing.T) {
	tx := newValidTransaction()
	tx.Signature.Algo = "unsupported-algo"

	err := checkTransaction(tx)

	assert.Contains(t, err.Get("tx.signature.algo"), crypto.ErrUnsupportedSignatureAlgorithm)
}

func testSubmittingTransactionWithInvalidEncodingOfValueFails(t *testing.T) {
	tx := newValidTransaction()
	tx.Signature.Value = "invalid-hex-encoding"

	err := checkTransaction(tx)

	assert.Contains(t, err.Get("tx.signature.value"), commands.ErrShouldBeHexEncoded)
}

func testSubmittingTransactionWithInvalidEncodingOfPubKeyFails(t *testing.T) {
	tx := newValidTransaction()
	tx.From = &commandspb.Transaction_PubKey{
		PubKey: "my-pub-key",
	}

	err := checkTransaction(tx)

	assert.Contains(t, err.Get("tx.from.pub_key"), commands.ErrShouldBeHexEncoded)
}

func testSubmittingTransactionWithInvalidSignatureFails(t *testing.T) {
	tx := newValidTransaction()
	tx.Signature.Value = "8ea1c9baab2919a73b6acd3dae15f515c9d9b191ac2a2cd9e7d7a2f9750da0793a88c8ee96a640e0de64c91d81770299769d4d4d93f81208e17573c836e3a810"

	err := checkTransaction(tx)

	assert.Contains(t, err.Get("tx.signature"), commands.ErrInvalidSignature)
}

func checkTransaction(cmd *commandspb.Transaction) commands.Errors {
	_, err := commands.CheckTransaction(cmd)

	e, ok := err.(commands.Errors)
	if !ok {
		return commands.NewErrors()
	}

	return e
}

func newValidTransaction() *commandspb.Transaction {
	return &commandspb.Transaction{
		InputData: []byte{8, 178, 211, 130, 220, 159, 158, 160, 128, 80, 210, 62, 0},
		Signature: &commandspb.Signature{
			Algo:    "vega/ed25519",
			Value:   "8ea1c9baab2919a73b6acd3dae15f515c9d9b191ac2a2cd9e7d7a2f9750da0793a88c8ee96a640e0de64c91d81770299769d4d4d93f81208e17573c836e3a80d",
			Version: 1,
		},
		From: &commandspb.Transaction_PubKey{
			PubKey: "b82756d3a3c5beff01152d3565e0c5c2235ccbe9c9d29ea4e760d981f53db7c6",
		},
		Version: 1,
	}
}
