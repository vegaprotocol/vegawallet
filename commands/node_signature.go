package commands

import (
	commandspb "code.vegaprotocol.io/go-wallet/internal/proto/commands/v1"
)

func CheckNodeSignature(cmd *commandspb.NodeSignature) error {
	return checkNodeSignature(cmd).ErrorOrNil()
}

func checkNodeSignature(cmd *commandspb.NodeSignature) Errors {
	errs := NewErrors()

	if cmd == nil {
		return errs.FinalAddForProperty("node_signature", ErrIsRequired)
	}

	if len(cmd.Id) == 0 {
		errs.AddForProperty("node_signature.id", ErrIsRequired)
	}

	if len(cmd.Sig) == 0 {
		errs.AddForProperty("node_signature.sig", ErrIsRequired)
	}

	if cmd.Kind == commandspb.NodeSignatureKind_NODE_SIGNATURE_KIND_UNSPECIFIED {
		errs.AddForProperty("node_signature.kind", ErrIsRequired)
	}
	if _, ok := commandspb.NodeSignatureKind_name[int32(cmd.Kind)]; !ok {
		errs.AddForProperty("node_signature.kind", ErrIsNotValid)
	}

	return errs
}
