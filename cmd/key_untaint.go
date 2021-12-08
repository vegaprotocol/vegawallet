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
	untaintKeyLong = cli.LongDesc(`
		Remove the taint from a key pair.

		If you tainted a key for security reasons, you should not untaint it.
	`)

	untaintKeyExample = cli.Examples(`
		# Untaint a key pair
		vegawallet key untaint --wallet WALLET --pubkey PUBKEY
	`)
)

type UntaintKeyHandler func(*wallet.UntaintKeyRequest) error

func NewCmdUntaintKey(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func(req *wallet.UntaintKeyRequest) error {
		s, err := wallets.InitialiseStore(rf.Home)
		if err != nil {
			return fmt.Errorf("couldn't initialise wallets store: %w", err)
		}

		return wallet.UntaintKey(s, req)
	}

	return BuildCmdUntaintKey(w, h, rf)
}

func BuildCmdUntaintKey(w io.Writer, handler UntaintKeyHandler, rf *RootFlags) *cobra.Command {
	f := &UntaintKeyFlags{}

	cmd := &cobra.Command{
		Use:     "untaint",
		Short:   "Remove the taint from a key pair",
		Long:    untaintKeyLong,
		Example: untaintKeyExample,
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
				PrintUntaintKeyResponse(w)
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
		"Public key to untaint (hex-encoded)",
	)
	cmd.Flags().StringVarP(&f.PassphraseFile,
		"passphrase-file", "p",
		"",
		"Path to the file containing the wallet's passphrase",
	)

	autoCompleteWallet(cmd, rf.Home)

	return cmd
}

type UntaintKeyFlags struct {
	Wallet         string
	PubKey         string
	PassphraseFile string
}

func (f *UntaintKeyFlags) Validate() (*wallet.UntaintKeyRequest, error) {
	req := &wallet.UntaintKeyRequest{}

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

func PrintUntaintKeyResponse(w io.Writer) {
	p := printer.NewInteractivePrinter(w)

	p.CheckMark().SuccessText("Untainting succeeded").NextSection()

	p.RedArrow().DangerText("Important").NextLine()
	p.Text("If you tainted a key for security reasons, you should not use it.").NextLine()

	p.BlueArrow().InfoText("Taint a key").NextLine()
	p.Text("To taint a key pair, see the following command:").NextSection()
	p.Code(fmt.Sprintf("%s key taint --help", os.Args[0])).NextLine()
}
