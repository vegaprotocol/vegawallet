package wallet

import (
	"code.vegaprotocol.io/go-wallet/commands"
	walletpb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/wallet/v1"
)

func CheckSubmitTransactionRequest(req *walletpb.SubmitTransactionRequest) commands.Errors {
	errs := commands.NewErrors()

	if len(req.PubKey) == 0 {
		errs.AddForProperty("submit_transaction_request.pub_key", commands.ErrIsRequired)
	}

	if req.Command == nil {
		return errs.FinalAddForProperty("submit_transaction_request.command", commands.ErrIsRequired)
	}

	var cmdErr error
	switch cmd := req.Command.(type) {
	case *walletpb.SubmitTransactionRequest_OrderSubmission:
		cmdErr = commands.CheckOrderSubmission(cmd.OrderSubmission)
	case *walletpb.SubmitTransactionRequest_OrderCancellation:
		cmdErr = commands.NewErrors()
	case *walletpb.SubmitTransactionRequest_OrderAmendment:
		cmdErr = commands.CheckOrderAmendment(cmd.OrderAmendment)
	case *walletpb.SubmitTransactionRequest_VoteSubmission:
		cmdErr = commands.CheckVoteSubmission(cmd.VoteSubmission)
	case *walletpb.SubmitTransactionRequest_WithdrawSubmission:
		cmdErr = commands.CheckWithdrawSubmission(cmd.WithdrawSubmission)
	case *walletpb.SubmitTransactionRequest_LiquidityProvisionSubmission:
		cmdErr = commands.CheckLiquidityProvisionSubmission(cmd.LiquidityProvisionSubmission)
	case *walletpb.SubmitTransactionRequest_ProposalSubmission:
		cmdErr = commands.CheckProposalSubmission(cmd.ProposalSubmission)
	case *walletpb.SubmitTransactionRequest_NodeRegistration:
		cmdErr = commands.CheckNodeRegistration(cmd.NodeRegistration)
	case *walletpb.SubmitTransactionRequest_NodeVote:
		cmdErr = commands.CheckNodeVote(cmd.NodeVote)
	case *walletpb.SubmitTransactionRequest_NodeSignature:
		cmdErr = commands.CheckNodeSignature(cmd.NodeSignature)
	case *walletpb.SubmitTransactionRequest_ChainEvent:
		cmdErr = commands.CheckChainEvent(cmd.ChainEvent)
	case *walletpb.SubmitTransactionRequest_OracleDataSubmission:
		cmdErr = commands.CheckOracleDataSubmission(cmd.OracleDataSubmission)
	default:
		errs.AddForProperty("input_data.command", commands.ErrIsNotSupported)
	}

	if cmdErr != nil {
		errs.Merge(toErrors(cmdErr))
	}

	return errs
}

func toErrors(err error) commands.Errors {
	e, ok := err.(commands.Errors)
	if !ok {
		errs := commands.NewErrors()
		return errs.FinalAdd(err)
	}
	return e
}
