package commands

import (
	"errors"
	"fmt"

	"code.vegaprotocol.io/protos/commands"
	commandspb "code.vegaprotocol.io/protos/vega/commands/v1"
	walletpb "code.vegaprotocol.io/protos/vega/wallet/v1"
	"github.com/golang/protobuf/proto"
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
	case *walletpb.SubmitTransactionRequest_AnnounceNode:
		cmdErr = commands.CheckAnnounceNode(cmd.AnnounceNode)
	case *walletpb.SubmitTransactionRequest_NodeVote:
		cmdErr = commands.CheckNodeVote(cmd.NodeVote)
	case *walletpb.SubmitTransactionRequest_NodeSignature:
		cmdErr = commands.CheckNodeSignature(cmd.NodeSignature)
	case *walletpb.SubmitTransactionRequest_ChainEvent:
		cmdErr = commands.CheckChainEvent(cmd.ChainEvent)
	case *walletpb.SubmitTransactionRequest_OracleDataSubmission:
		cmdErr = commands.CheckOracleDataSubmission(cmd.OracleDataSubmission)
	case *walletpb.SubmitTransactionRequest_UndelegateSubmission:
		cmdErr = commands.CheckUndelegateSubmission(cmd.UndelegateSubmission)
	case *walletpb.SubmitTransactionRequest_DelegateSubmission:
		cmdErr = commands.CheckDelegateSubmission(cmd.DelegateSubmission)
	case *walletpb.SubmitTransactionRequest_LiquidityProvisionCancellation:
		cmdErr = commands.CheckLiquidityProvisionCancellation(cmd.LiquidityProvisionCancellation)
	case *walletpb.SubmitTransactionRequest_LiquidityProvisionAmendment:
		cmdErr = commands.CheckLiquidityProvisionAmendment(cmd.LiquidityProvisionAmendment)
	case *walletpb.SubmitTransactionRequest_Transfer:
		cmdErr = commands.CheckTransfer(cmd.Transfer)
	case *walletpb.SubmitTransactionRequest_CancelTransfer:
		cmdErr = commands.CheckCancelTransfer(cmd.CancelTransfer)
	case *walletpb.SubmitTransactionRequest_KeyRotateSubmission:
		cmdErr = commands.CheckKeyRotateSubmission(cmd.KeyRotateSubmission)
	case *walletpb.SubmitTransactionRequest_EthereumKeyRotateSubmission:
		cmdErr = commands.CheckEthereumKeyRotateSubmission(cmd.EthereumKeyRotateSubmission)
	default:
		errs.AddForProperty("input_data.command", commands.ErrIsNotSupported)
	}

	if cmdErr != nil {
		errs.Merge(toErrors(cmdErr))
	}

	return errs
}

func ToMarshaledInputData(req *walletpb.SubmitTransactionRequest, height uint64) ([]byte, error) {
	data := commands.NewInputData(height)
	wrapRequestCommandIntoInputData(data, req)
	return proto.Marshal(data)
}

func wrapRequestCommandIntoInputData(data *commandspb.InputData, req *walletpb.SubmitTransactionRequest) {
	switch cmd := req.Command.(type) {
	case *walletpb.SubmitTransactionRequest_OrderSubmission:
		data.Command = &commandspb.InputData_OrderSubmission{
			OrderSubmission: req.GetOrderSubmission(),
		}
	case *walletpb.SubmitTransactionRequest_OrderCancellation:
		data.Command = &commandspb.InputData_OrderCancellation{
			OrderCancellation: req.GetOrderCancellation(),
		}
	case *walletpb.SubmitTransactionRequest_OrderAmendment:
		data.Command = &commandspb.InputData_OrderAmendment{
			OrderAmendment: req.GetOrderAmendment(),
		}
	case *walletpb.SubmitTransactionRequest_VoteSubmission:
		data.Command = &commandspb.InputData_VoteSubmission{
			VoteSubmission: req.GetVoteSubmission(),
		}
	case *walletpb.SubmitTransactionRequest_WithdrawSubmission:
		data.Command = &commandspb.InputData_WithdrawSubmission{
			WithdrawSubmission: req.GetWithdrawSubmission(),
		}
	case *walletpb.SubmitTransactionRequest_LiquidityProvisionSubmission:
		data.Command = &commandspb.InputData_LiquidityProvisionSubmission{
			LiquidityProvisionSubmission: req.GetLiquidityProvisionSubmission(),
		}
	case *walletpb.SubmitTransactionRequest_ProposalSubmission:
		data.Command = &commandspb.InputData_ProposalSubmission{
			ProposalSubmission: req.GetProposalSubmission(),
		}
	case *walletpb.SubmitTransactionRequest_AnnounceNode:
		data.Command = &commandspb.InputData_AnnounceNode{
			AnnounceNode: req.GetAnnounceNode(),
		}
	case *walletpb.SubmitTransactionRequest_NodeVote:
		data.Command = &commandspb.InputData_NodeVote{
			NodeVote: req.GetNodeVote(),
		}
	case *walletpb.SubmitTransactionRequest_NodeSignature:
		data.Command = &commandspb.InputData_NodeSignature{
			NodeSignature: req.GetNodeSignature(),
		}
	case *walletpb.SubmitTransactionRequest_ChainEvent:
		data.Command = &commandspb.InputData_ChainEvent{
			ChainEvent: req.GetChainEvent(),
		}
	case *walletpb.SubmitTransactionRequest_OracleDataSubmission:
		data.Command = &commandspb.InputData_OracleDataSubmission{
			OracleDataSubmission: req.GetOracleDataSubmission(),
		}
	case *walletpb.SubmitTransactionRequest_DelegateSubmission:
		data.Command = &commandspb.InputData_DelegateSubmission{
			DelegateSubmission: req.GetDelegateSubmission(),
		}
	case *walletpb.SubmitTransactionRequest_UndelegateSubmission:
		data.Command = &commandspb.InputData_UndelegateSubmission{
			UndelegateSubmission: req.GetUndelegateSubmission(),
		}
	case *walletpb.SubmitTransactionRequest_LiquidityProvisionCancellation:
		data.Command = &commandspb.InputData_LiquidityProvisionCancellation{
			LiquidityProvisionCancellation: req.GetLiquidityProvisionCancellation(),
		}
	case *walletpb.SubmitTransactionRequest_LiquidityProvisionAmendment:
		data.Command = &commandspb.InputData_LiquidityProvisionAmendment{
			LiquidityProvisionAmendment: req.GetLiquidityProvisionAmendment(),
		}
	case *walletpb.SubmitTransactionRequest_Transfer:
		data.Command = &commandspb.InputData_Transfer{
			Transfer: req.GetTransfer(),
		}
	case *walletpb.SubmitTransactionRequest_CancelTransfer:
		data.Command = &commandspb.InputData_CancelTransfer{
			CancelTransfer: req.GetCancelTransfer(),
		}
	case *walletpb.SubmitTransactionRequest_KeyRotateSubmission:
		data.Command = &commandspb.InputData_KeyRotateSubmission{
			KeyRotateSubmission: req.GetKeyRotateSubmission(),
		}
	case *walletpb.SubmitTransactionRequest_EthereumKeyRotateSubmission:
		data.Command = &commandspb.InputData_EthereumKeyRotateSubmission{
			EthereumKeyRotateSubmission: req.GetEthereumKeyRotateSubmission(),
		}
	default:
		panic(fmt.Sprintf("command %v is not supported", cmd))
	}
}

func toErrors(err error) commands.Errors {
	errs := &commands.Errors{}
	if !errors.As(err, errs) {
		errs := commands.NewErrors()
		return errs.FinalAdd(err)
	}
	return *errs
}
