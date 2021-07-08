package commands

import commandspb "code.vegaprotocol.io/go-wallet/internal/proto/commands/v1"

func CheckNodeRegistration(cmd *commandspb.NodeRegistration) error {
	return checkNodeRegistration(cmd).ErrorOrNil()
}

func checkNodeRegistration(cmd *commandspb.NodeRegistration) Errors {
	errs := NewErrors()

	if cmd == nil {
		return errs.FinalAddForProperty("node_registration", ErrIsRequired)
	}

	if len(cmd.PubKey) == 0 {
		errs.AddForProperty("node_registration.pub_key", ErrIsRequired)
	}

	if len(cmd.ChainPubKey) == 0 {
		errs.AddForProperty("node_registration.chain_pub_key", ErrIsRequired)
	}

	return errs
}
