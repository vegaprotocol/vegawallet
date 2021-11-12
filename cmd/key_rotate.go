package cmd

import (
	"fmt"
	"io"

	"code.vegaprotocol.io/vegawallet/cmd/cli"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/wallet"
	"code.vegaprotocol.io/vegawallet/wallets"
	"github.com/spf13/cobra"
)

var (
	rotateKeyLong = cli.LongDesc(`
		Get a signed key rotation transaction as a Base64 encoded string by choosing a new public key to rotate to and target block height.

		Later the transaction is applied to Vega protocol through the wallet's "send tx" command.
	`)

	rotateKeyExample = cli.Examples(`
		Given that those variables exists:
			- WALLET - Vega wallet to be used.
			- PUBLIC_KEY - A newly generated public key. Should be generate by wallet's "generate" command.
			- TARGET_HEIGHT - Height of block where the new public key change will take effect
			- TX_HEIGHT - It should be close to the current block height when the transaction is applied, with a threshold of ~ - 150 blocks.

		# Get signed transaction for rotating to new key public key
		vegawallet key rotate --wallet WALLET --tx-height TX_HEIGHT --target-height TARGET_HEIGHT PUBLIC_KEY
	`)
)

type RotateKeyHandler func(*wallet.RotateKeyRequest) (*wallet.RotateKeyResponse, error)

func NewCmdRotateKey(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func(req *wallet.RotateKeyRequest) (*wallet.RotateKeyResponse, error) {
		s, err := wallets.InitialiseStore(rf.Home)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialise wallets store: %w", err)
		}

		return wallet.RotateKey(s, req)
	}

	return BuildCmdRotateKey(w, h, rf)
}

func BuildCmdRotateKey(w io.Writer, handler RotateKeyHandler, rf *RootFlags) *cobra.Command {
	f := RotateKeyFlags{}

	cmd := &cobra.Command{
		Use:     "rotate",
		Short:   "Get signed key rotation transaction",
		Long:    rotateKeyLong,
		Example: rotateKeyExample,
		RunE: func(_ *cobra.Command, args []string) error {
			if aLen := len(args); aLen == 0 {
				return flags.ArgMustBeSpecifiedError("public-key")
			} else if aLen > 1 {
				return fmt.Errorf("too many arguments specified")
			}
			f.NewPublicKey = args[0]

			req, err := f.Validate()
			if err != nil {
				return err
			}

			resp, err := handler(req)
			if err != nil {
				return err
			}

			switch rf.Output {
			case flags.InteractiveOutput:
				PrintRotateKeyResponse(w, resp)
			case flags.JSONOutput:
				return printer.FprintJSON(w, resp)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&f.Wallet,
		"wallet", "w",
		"",
		"Wallet holding the master key and new public key",
	)
	cmd.Flags().StringVarP(&f.PassphraseFile,
		"passphrase-file", "p",
		"",
		"Path to the file containing the wallet's passphrase",
	)
	cmd.Flags().Uint32VarP(&f.TXBlockHeight,
		"tx-height", "xh",
		0,
		"Block height of which the transaction will be applied at",
	)
	cmd.Flags().Uint32VarP(&f.TargetBlockHeight,
		"target-height", "th",
		0,
		"Block height from which the new public key will be used",
	)

	return cmd
}

type RotateKeyFlags struct {
	Wallet            string
	PassphraseFile    string
	NewPublicKey      string
	TXBlockHeight     uint32
	TargetBlockHeight uint32
}

func (f *RotateKeyFlags) Validate() (*wallet.RotateKeyRequest, error) {
	req := &wallet.RotateKeyRequest{
		NewPublicKey: f.NewPublicKey,
	}

	if f.TargetBlockHeight == 0 {
		return nil, flags.FlagMustBeSpecifiedError("target-height")
	}
	req.TargetBlockHeight = f.TargetBlockHeight

	if f.TXBlockHeight == 0 {
		return nil, flags.FlagMustBeSpecifiedError("tx-height")
	}
	req.TXBlockHeight = f.TXBlockHeight

	if req.TargetBlockHeight <= req.TXBlockHeight {
		return nil, fmt.Errorf("--target-height flag must be greater then --tx-height")
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

	return req, nil
}

func PrintRotateKeyResponse(w io.Writer, req *wallet.RotateKeyResponse) {
	p := printer.NewInteractivePrinter(w)
	p.CheckMark().SuccessText("Key rotation succeeded").NextSection()
	p.Text("Base64 encoded transaction:").NextLine()
	p.Text(req.Base64Transaction).NextLine()
	p.Text("New public key:").NextLine()
	p.Text(req.NewPublicKey).NextLine()
	p.Text("Master public key used:").NextLine()
	p.Text(req.MasterPublicKey).NextLine()
}
