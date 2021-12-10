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
	describeKeyLong = cli.LongDesc(`
		Describe all known information about the specified key pair
	`)

	describeKeyExample = cli.Examples(`
		# Describe a key
		vegawallet key describe --wallet WALLET --pubkey PUBKEY
	`)
)

type DescribeKeyHandler func(*wallet.DescribeKeyRequest) (*wallet.DescribeKeyResponse, error)

func NewCmdDescribeKey(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func(req *wallet.DescribeKeyRequest) (*wallet.DescribeKeyResponse, error) {
		s, err := wallets.InitialiseStore(rf.Home)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialise wallets store: %w", err)
		}

		return wallet.DescribeKey(s, req)
	}

	return BuildCmdDescribeKey(w, h, rf)
}

func BuildCmdDescribeKey(w io.Writer, handler DescribeKeyHandler, rf *RootFlags) *cobra.Command {
	f := &DescribeKeyFlags{}

	cmd := &cobra.Command{
		Use:     "describe",
		Short:   "Describe the specified key pair",
		Long:    describeKeyLong,
		Example: describeKeyExample,
		RunE: func(_ *cobra.Command, _ []string) error {
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
				PrintDescribeKeyResponse(w, resp)
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
	cmd.Flags().StringVarP(&f.PubKey,
		"pubkey", "k",
		"",
		"Public key to describe (hex-encoded)",
	)
	cmd.Flags().StringVarP(&f.PassphraseFile,
		"passphrase-file", "p",
		"",
		"Path to the file containing the wallet's passphrase",
	)

	autoCompleteWallet(cmd, rf.Home)

	return cmd
}

type DescribeKeyFlags struct {
	Wallet         string
	PassphraseFile string
	PubKey         string
}

func (f *DescribeKeyFlags) Validate() (*wallet.DescribeKeyRequest, error) {
	req := &wallet.DescribeKeyRequest{}

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

func PrintDescribeKeyResponse(w io.Writer, resp *wallet.DescribeKeyResponse) {
	p := printer.NewInteractivePrinter(w).NextLine()

	p.Text("Name:              ").WarningText(wallet.GetKeyName(resp.Meta)).NextLine()
	p.Text("Public key:        ").WarningText(resp.PublicKey).NextLine()
	p.Text("Algorithm Name:    ").WarningText(resp.Algorithm.Name).NextLine()
	p.Text("Algorithm Version: ").WarningText(fmt.Sprint(resp.Algorithm.Version)).NextSection()

	p.Text("Key pair is: ")
	switch resp.IsTainted {
	case true:
		p.DangerText("tainted").NextLine()
	case false:
		p.SuccessText("not tainted").NextLine()
	}
	p.Text("Tainting a key pair marks it as unsafe to use and ensures it will not be used to sign transactions.").NextLine()
	p.Text("This mechanism is useful when the key pair has been compromised.").NextSection()

	p.Text("Metadata:").NextLine()
	printMeta(p, resp.Meta)
}
