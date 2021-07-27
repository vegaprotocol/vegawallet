package commands

import (
	"encoding/hex"
	"errors"

	"code.vegaprotocol.io/go-wallet/crypto"
	commandspb "code.vegaprotocol.io/go-wallet/internal/proto/commands/v1"
	"github.com/golang/protobuf/proto"
)

var (
	ErrInvalidSignature   = errors.New("invalid signature")
	ErrShouldBeHexEncoded = errors.New("should be hex encoded")
)

func NewTransaction(pubKey string, data []byte, signature *commandspb.Signature) *commandspb.Transaction {
	return &commandspb.Transaction{
		InputData: data,
		Signature: signature,
		From: &commandspb.Transaction_PubKey{
			PubKey: pubKey,
		},
		Version: 1,
	}
}

func NewInputData(height uint64) *commandspb.InputData {
	return &commandspb.InputData{
		Nonce:       crypto.NewNonce(),
		BlockHeight: height,
	}
}

func CheckTransaction(tx *commandspb.Transaction) (*commandspb.InputData, error) {
	errs := NewErrors()

	if tx == nil {
		return nil, errs.FinalAddForProperty("tx", ErrIsRequired)
	}

	if len(tx.InputData) == 0 {
		errs.AddForProperty("tx.input_data", ErrIsRequired)
	}
	if tx.Signature == nil {
		errs.AddForProperty("tx.signature", ErrIsRequired)
	} else {
		errs.Merge(checkSignature(tx.Signature))
	}
	if tx.From == nil {
		errs.AddForProperty("tx.from", ErrIsRequired)
	} else if len(tx.GetPubKey()) == 0 {
		errs.AddForProperty("tx.from.pub_key", ErrIsRequired)
	}

	if !errs.Empty() {
		return nil, errs.ErrorOrNil()
	}

	errs.Merge(validateSignature(tx.InputData, tx.Signature, tx.GetPubKey()))
	if !errs.Empty() {
		return nil, errs.ErrorOrNil()
	}

	inputData, errs := checkInputData(tx.InputData)
	if !errs.Empty() {
		return nil, errs.ErrorOrNil()
	}
	return inputData, nil
}

func validateSignature(inputData []byte, signature *commandspb.Signature, pubKey string) Errors {
	errs := NewErrors()

	validator, err := crypto.NewSignatureAlgorithm(signature.Algo, signature.Version)
	if err != nil {
		return errs.FinalAddForProperty("tx.signature.algo", err)
	}

	decodedSig, err := hex.DecodeString(signature.Value)
	if err != nil {
		return errs.FinalAddForProperty("tx.signature.value", ErrShouldBeHexEncoded)
	}

	decodedPubKey, err := hex.DecodeString(pubKey)
	if err != nil {
		return errs.FinalAddForProperty("tx.from.pub_key", ErrShouldBeHexEncoded)
	}

	ok, err := validator.Verify(decodedPubKey, inputData, decodedSig)
	if err != nil {
		return errs.FinalAdd(err)
	}
	if !ok {
		return errs.FinalAddForProperty("tx.signature", ErrInvalidSignature)
	}
	return errs
}

func checkInputData(inputData []byte) (*commandspb.InputData, Errors) {
	errs := NewErrors()

	input := commandspb.InputData{}
	err := proto.Unmarshal(inputData, &input)
	if err != nil {
		return nil, errs.FinalAdd(err)
	}

	if input.Nonce == 0 {
		errs.AddForProperty("tx.input_data.nonce", ErrMustBePositive)
	}

	if input.Command == nil {
		errs.AddForProperty("tx.input_data.command", ErrIsRequired)
	} else {
		switch cmd := input.Command.(type) {
		case *commandspb.InputData_OrderSubmission:
			errs.Merge(checkOrderSubmission(cmd.OrderSubmission))
		case *commandspb.InputData_OrderCancellation:
			break // No verification to be made
		case *commandspb.InputData_OrderAmendment:
			errs.Merge(checkOrderAmendment(cmd.OrderAmendment))
		case *commandspb.InputData_VoteSubmission:
			errs.Merge(checkVoteSubmission(cmd.VoteSubmission))
		case *commandspb.InputData_WithdrawSubmission:
			errs.Merge(checkWithdrawSubmission(cmd.WithdrawSubmission))
		case *commandspb.InputData_LiquidityProvisionSubmission:
			errs.Merge(checkLiquidityProvisionSubmission(cmd.LiquidityProvisionSubmission))
		case *commandspb.InputData_ProposalSubmission:
			errs.Merge(checkProposalSubmission(cmd.ProposalSubmission))
		case *commandspb.InputData_NodeRegistration:
			errs.Merge(checkNodeRegistration(cmd.NodeRegistration))
		case *commandspb.InputData_NodeVote:
			errs.Merge(checkNodeVote(cmd.NodeVote))
		case *commandspb.InputData_NodeSignature:
			errs.Merge(checkNodeSignature(cmd.NodeSignature))
		case *commandspb.InputData_ChainEvent:
			errs.Merge(checkChainEvent(cmd.ChainEvent))
		case *commandspb.InputData_OracleDataSubmission:
			errs.Merge(checkOracleDataSubmission(cmd.OracleDataSubmission))
		default:
			errs.AddForProperty("tx.input_data.command", ErrIsNotSupported)
		}
	}

	return &input, errs
}

func checkSignature(signature *commandspb.Signature) Errors {
	errs := NewErrors()
	if len(signature.Value) == 0 {
		errs.AddForProperty("tx.signature.value", ErrIsRequired)
	}
	if len(signature.Algo) == 0 {
		errs.AddForProperty("tx.signature.algo", ErrIsRequired)
	}
	return errs
}
