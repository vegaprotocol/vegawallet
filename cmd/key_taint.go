package cmd

import (
	"fmt"
	"io"
	"os"

	"code.vegaprotocol.io/vegawallet/cmd/cli"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/wallet"
	"code.vegaprotocol.io/vegawallet/wallets"
	"github.com/spf13/cobra"
)

var (
	taintKeyLong = cli.LongDesc(`
		Tainting a key pair marks it as unsafe to use and ensure it will not be
		used to sign transactions.

		This mechanism is useful when the key pair has been compromised.
	`)

	taintKeyExample = cli.Examples(`
		# Taint a key pair
		vegawallet key taint --wallet WALLET --pubkey PUBKEY
	`)
)

type TaintKeyHandler func(*wallet.TaintKeyRequest) error

func NewCmdTaintKey(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func(req *wallet.TaintKeyRequest) error {
		s, err := wallets.InitialiseStore(rf.Home)
		if err != nil {
			return fmt.Errorf("couldn't initialise wallets store: %w", err)
		}

		return wallet.TaintKey(s, req)
	}

	return BuildCmdTaintKey(w, h, rf)
}

func BuildCmdTaintKey(w io.Writer, handler TaintKeyHandler, rf *RootFlags) *cobra.Command {
	f := &TaintKeyFlags{}

	cmd := &cobra.Command{
		Use:     "taint",
		Short:   "Mark a key pair as tainted",
		Long:    taintKeyLong,
		Example: taintKeyExample,
		RunE: func(_ *cobra.Command, _ []string) error {
			req, err := f.Validate()
			if err != nil {
				return err
			}

			if err := handler(req); err != nil {
				return err
			}

			switch rf.Output {
			case flags.InteractiveOutput:
				PrintTaintKeyResponse(w)
			case flags.JSONOutput:
				return nil
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
		"Public key to taint (hex-encoded)",
	)
	cmd.Flags().StringVarP(&f.PassphraseFile,
		"passphrase-file", "p",
		"",
		"Path to the file containing the wallet's passphrase",
	)

	return cmd
}

type TaintKeyFlags struct {
	Wallet         string
	PubKey         string
	PassphraseFile string
}

func (f *TaintKeyFlags) Validate() (*wallet.TaintKeyRequest, error) {
	req := &wallet.TaintKeyRequest{}

	if len(f.Wallet) == 0 {
		return nil, flags.FlagMustBeSpecifiedError("wallet")
	}
	req.Wallet = f.Wallet

	if len(f.PubKey) == 0 {
		return nil, flags.FlagMustBeSpecifiedError("pubkey")
	}
	req.PubKey = f.PubKey

	passphrase, err := flags.GetPassphrase(f.PassphraseFile)
	if err != nil {
		return nil, err
	}
	req.Passphrase = passphrase

	return req, nil
}

func PrintTaintKeyResponse(w io.Writer) {
	p := printer.NewInteractivePrinter(w)

	p.CheckMark().SuccessText("Tainting succeeded").NextSection()

	p.RedArrow().DangerText("Important").NextLine()
	p.Text("If you tainted a key for security reasons, you should not untaint it.").NextSection()

	p.BlueArrow().InfoText("Untaint a key").NextLine()
	p.Text("You may have tainted a key pair by mistake. If you want to untaint it, see the following command:").NextSection()
	p.Code(fmt.Sprintf("%s key untaint --help", os.Args[0])).NextSection()
}
