package commands

import (
	commandspb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/commands/v1"
)

func CheckChainEvent(cmd *commandspb.ChainEvent) error {
	return checkChainEvent(cmd).ErrorOrNil()
}

func checkChainEvent(cmd *commandspb.ChainEvent) Errors {
	errs := NewErrors()

	if cmd == nil {
		return errs.FinalAddForProperty("chain_event", ErrIsRequired)
	}

	if cmd.Event != nil && isBuiltInEvent(cmd) {
		return errs
	}

	if cmd.Event == nil {
		errs.AddForProperty("chain_event.event", ErrIsRequired)
	}

	if len(cmd.TxId) == 0 {
		errs.AddForProperty("chain_event.tx_id", ErrIsRequired)
	}

	return errs
}

func isBuiltInEvent(cmd *commandspb.ChainEvent) bool {
	switch cmd.Event.(type) {
	case *commandspb.ChainEvent_Builtin:
		return true
	default:
		return false
	}
}
