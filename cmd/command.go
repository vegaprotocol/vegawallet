package cmd

import (
	"fmt"
	"io"

	"code.vegaprotocol.io/vegawallet/cmd/cli"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"go.uber.org/zap/zapcore"

	"github.com/spf13/cobra"
)

var (
	commandLong = cli.LongDesc(`
		DEPRECATED: use "send command" instead.

		Send a command to a Vega node via the gRPC API. The command can be sent to 
		any node of a registered network or to a specific node address.

		The command should be a Vega command formatted as a JSON payload, as follows:

		'{"commandName": {"someProperty": "someValue"} }'

		For vote submission, it will look like this:

		'{"voteSubmission": {"proposalId": "some-id", "value": "VALUE_YES"}}'
	`)

	commandExample = cli.Examples(`
		# Send a command to a registered network
		vegawallet command --network NETWORK --wallet WALLET --pubkey PUBKEY COMMAND

		# Send a command to a specific Vega node address
		vegawallet command --node-address ADDRESS --wallet WALLET --pubkey PUBKEY COMMAND

		# Send a command with a log level set to debug
		vegawallet command --network NETWORK --wallet WALLET --pubkey PUBKEY --level debug COMMAND

		# Send a command with a maximum of 10 retry
		vegawallet command --network NETWORK --wallet WALLET --pubkey PUBKEY --retries 10 COMMAND
	`)
)

func NewCmdCommand(w io.Writer, rf *RootFlags) *cobra.Command {
	return BuildCmdCommand(w, SendCommand, rf)
}

func BuildCmdCommand(w io.Writer, handler SendCommandHandler, rf *RootFlags) *cobra.Command {
	f := &SendCommandFlags{}

	cmd := &cobra.Command{
		Use:     "command",
		Short:   "Send a command to a Vega node",
		Long:    commandLong,
		Example: commandExample,
		RunE: func(_ *cobra.Command, args []string) error {
			if rf.Output == flags.InteractiveOutput {
				p := printer.NewInteractivePrinter(w)
				p.BangMark().DangerText("This command is DEPRECATED.").NextLine()
				p.BangMark().DangerText("Use `vegawallet send command` instead").NextSection()
			}

			if aLen := len(args); aLen == 0 {
				return flags.ArgMustBeSpecifiedError("command")
			} else if aLen > 1 {
				return flags.TooManyArgsError("command")
			}
			f.RawCommand = args[0]

			req, err := f.Validate()
			if err != nil {
				return err
			}

			if err := handler(w, rf, req); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&f.Network,
		"network", "n",
		"",
		"Network to which the command is sent",
	)
	cmd.Flags().StringVar(&f.NodeAddress,
		"node-address",
		"",
		"Vega node address to which the command is sent",
	)
	cmd.Flags().StringVarP(&f.Wallet,
		"wallet", "w",
		"",
		"Wallet holding the public key",
	)
	cmd.Flags().StringVarP(&f.PubKey,
		"pubkey", "k",
		"",
		"Public key of the key pair to use for signing (hex-encoded)",
	)
	cmd.Flags().StringVarP(&f.PassphraseFile,
		"passphrase-file", "p",
		"",
		"Path to the file containing the wallet's passphrase",
	)
	cmd.Flags().StringVar(&f.LogLevel,
		"level",
		zapcore.InfoLevel.String(),
		fmt.Sprintf("Set the log level: %v", SupportedLogLevels),
	)
	cmd.Flags().Uint64Var(&f.Retries,
		"retries",
		DefaultForwarderRetryCount,
		"Number of retries when contacting the Vega node",
	)

	return cmd
}
