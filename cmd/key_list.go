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
	listKeysLong = cli.LongDesc(`
		List the keys of a given wallet.
	`)

	listKeysExample = cli.Examples(`
		# List all keys
		vegawallet key list --wallet WALLET
	`)
)

type ListKeysHandler func(*wallet.ListKeysRequest) (*wallet.ListKeysResponse, error)

func NewCmdListKeys(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func(req *wallet.ListKeysRequest) (*wallet.ListKeysResponse, error) {
		s, err := wallets.InitialiseStore(rf.Home)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialise wallets store: %w", err)
		}

		return wallet.ListKeys(s, req)
	}

	return BuildCmdListKeys(w, h, rf)
}

func BuildCmdListKeys(w io.Writer, handler ListKeysHandler, rf *RootFlags) *cobra.Command {
	f := &ListKeysFlags{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List the keys of a given wallet",
		Long:    listKeysLong,
		Example: listKeysExample,
		RunE: func(_ *cobra.Command, _ []string) error {
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
				PrintListKeysResponse(w, resp)
			case flags.JSONOutput:
				return printer.FprintJSON(w, resp)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&f.Wallet,
		"wallet", "w",
		"",
		"Name of the wallet to use",
	)
	cmd.Flags().StringVarP(&f.PassphraseFile,
		"passphrase-file", "p",
		"",
		"Path to the file containing the wallet's passphrase",
	)

	addOutputFlag(cmd, &f.Output)

	autoCompleteWallet(cmd, rf.Home)

	return cmd
}

type ListKeysFlags struct {
	Wallet         string
	PassphraseFile string
	Output         string
}

func (f *ListKeysFlags) Validate() (*wallet.ListKeysRequest, error) {
	req := &wallet.ListKeysRequest{}

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

	return req, nil
}

func PrintListKeysResponse(w io.Writer, resp *wallet.ListKeysResponse) {
	p := printer.NewInteractivePrinter(w)

	for i, key := range resp.Keys {
		if i != 0 {
			p.NextLine()
		}
		p.Text("Name:       ").WarningText(key.Name).NextLine()
		p.Text("Public key: ").WarningText(key.PublicKey).NextLine()
	}
}
