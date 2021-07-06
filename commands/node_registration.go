package commands

import commandspb "github.com/vegaprotocol/api/grpc/clients/go/generated/code.vegaprotocol.io/vega/proto/commands/v1"

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
