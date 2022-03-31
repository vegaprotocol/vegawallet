package cmd

import (
	"fmt"
	"io"
	"os"

	walletpb "code.vegaprotocol.io/protos/vega/wallet/v1"
	"code.vegaprotocol.io/vegawallet/cmd/cli"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	wcommands "code.vegaprotocol.io/vegawallet/commands"
	"code.vegaprotocol.io/vegawallet/wallet"
	"code.vegaprotocol.io/vegawallet/wallets"
	"github.com/golang/protobuf/jsonpb"
	"github.com/spf13/cobra"
)

var (
	signCommandLong = cli.LongDesc(`
		Sign a command using the specified wallet and public key and bundle it as a
		transaction ready to be sent. The resulting transaction is base64-encoded and
		can be sent using the command "tx send".

		The command should be a Vega command formatted as a JSON payload, as follows:

		'{"commandName": {"someProperty": "someValue"} }'

		For vote submission, it will look like this:

		'{"voteSubmission": {"proposalId": "some-id", "value": "VALUE_YES"}}'
	`)

	signCommandExample = cli.Examples(`
		# Sign a command
		vegawallet command sign --wallet WALLET --pubkey PUBKEY --tx-height TX_HEIGHT COMMAND

		# To decode the result, save the result in a file and use the command 
		# "base64"
		vegawallet command sign --wallet WALLET --pubkey PUBKEY --tx-height TX_HEIGHT COMMAND > result.txt
		base64 --decode --input result.txt
	`)
)

type SignCommandHandler func(*wallet.SignCommandRequest) (*wallet.SignCommandResponse, error)

func NewCmdCommandSign(w io.Writer, rf *RootFlags) *cobra.Command {
	handler := func(req *wallet.SignCommandRequest) (*wallet.SignCommandResponse, error) {
		store, err := wallets.InitialiseStore(rf.Home)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialise wallets store: %w", err)
		}

		return wallet.SignCommand(store, req)
	}

	return BuildCmdCommandSign(w, handler, rf)
}

func BuildCmdCommandSign(w io.Writer, handler SignCommandHandler, rf *RootFlags) *cobra.Command {
	f := &SignCommandFlags{}

	cmd := &cobra.Command{
		Use:     "sign",
		Short:   "Sign a command for offline use",
		Long:    signCommandLong,
		Example: signCommandExample,
		RunE: func(_ *cobra.Command, args []string) error {
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

			resp, err := handler(req)
			if err != nil {
				return err
			}

			switch f.Output {
			case flags.InteractiveOutput:
				PrintSignCommandResponse(w, resp)
			case flags.JSONOutput:
				return printer.FprintJSON(w, resp)
			}

			return nil
		},
	}

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
	cmd.Flags().Uint64Var(&f.TxBlockHeight,
		"tx-height",
		0,
		"It should be close to the current block height when the transaction is applied, with a threshold of ~ - 150 blocks.",
	)

	addOutputFlag(cmd, &f.Output)

	autoCompleteWallet(cmd, rf.Home)

	return cmd
}

type SignCommandFlags struct {
	Wallet         string
	PubKey         string
	PassphraseFile string
	RawCommand     string
	TxBlockHeight  uint64
	Output         string
}

func (f *SignCommandFlags) Validate() (*wallet.SignCommandRequest, error) {
	req := &wallet.SignCommandRequest{}

	if err := flags.ValidateOutput(f.Output); err != nil {
		return nil, err
	}

	if len(f.Wallet) == 0 {
		return nil, flags.FlagMustBeSpecifiedError("wallet")
	}
	req.Wallet = f.Wallet

	passphrase, err := flags.GetPassphrase(f.PassphraseFile)
	if err != nil {
		return nil, err
	}
	req.Passphrase = passphrase

	if f.TxBlockHeight == 0 {
		return nil, flags.FlagMustBeSpecifiedError("tx-height")
	}
	req.TxBlockHeight = f.TxBlockHeight

	if len(f.PubKey) == 0 {
		return nil, flags.FlagMustBeSpecifiedError("pubkey")
	}
	if len(f.RawCommand) == 0 {
		return nil, flags.ArgMustBeSpecifiedError("command")
	}
	request := &walletpb.SubmitTransactionRequest{}
	if err := jsonpb.UnmarshalString(f.RawCommand, request); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal command as request: %w", err)
	}
	if len(request.PubKey) != 0 {
		return nil, ErrDoNotSetPubKeyInCommand
	}
	request.PubKey = f.PubKey
	request.Propagate = true
	req.Request = request
	if errs := wcommands.CheckSubmitTransactionRequest(req.Request); !errs.Empty() {
		return nil, fmt.Errorf("invalid request: %w", errs)
	}

	return req, nil
}

func PrintSignCommandResponse(w io.Writer, req *wallet.SignCommandResponse) {
	p := printer.NewInteractivePrinter(w)

	p.CheckMark().SuccessText("Command signature successful").NextSection()
	p.Text("Transaction (base64-encoded):").NextLine().WarningText(req.Base64Transaction).NextSection()

	p.BlueArrow().InfoText("Send a transaction").NextLine()
	p.Text("To send a raw transaction, see the following command:").NextSection()
	p.Code(fmt.Sprintf("%s tx send --help", os.Args[0])).NextSection()
}
