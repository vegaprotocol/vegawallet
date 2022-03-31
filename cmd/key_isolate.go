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
	isolateKeyLong = cli.LongDesc(`
		Extract the specified key pair into an isolated wallet.
	
		An isolated wallet is a wallet that contains a single key pair and that
		has been stripped from its cryptographic node.
		
		Removing the cryptographic node from the wallet minimizes the impact of a
		stolen wallet, as it makes it impossible to retrieve or generate keys out
		of it.

		This creates a wallet that is only able to sign and verify transactions.

		This adds an extra layer of security.
	`)

	isolateKeyExample = cli.Examples(`
		# Isolate a key pair
		vegawallet key isolate --wallet WALLET --pubkey PUBKEY
	`)
)

type IsolateKeyHandler func(*wallet.IsolateKeyRequest) (*wallet.IsolateKeyResponse, error)

func NewCmdIsolateKey(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func(req *wallet.IsolateKeyRequest) (*wallet.IsolateKeyResponse, error) {
		s, err := wallets.InitialiseStore(rf.Home)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialise wallets store: %w", err)
		}

		return wallet.IsolateKey(s, req)
	}

	return BuildCmdIsolateKey(w, h, rf)
}

func BuildCmdIsolateKey(w io.Writer, handler IsolateKeyHandler, rf *RootFlags) *cobra.Command {
	f := &IsolateKeyFlags{}

	cmd := &cobra.Command{
		Use:     "isolate",
		Short:   "Extract the specified key pair into an isolated wallet",
		Long:    isolateKeyLong,
		Example: isolateKeyExample,
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
				PrintIsolateKeyResponse(w, resp)
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
		"Public key to isolate (hex-encoded)",
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

type IsolateKeyFlags struct {
	Wallet         string
	PubKey         string
	PassphraseFile string
	Output         string
}

func (f *IsolateKeyFlags) Validate() (*wallet.IsolateKeyRequest, error) {
	req := &wallet.IsolateKeyRequest{}

	if err := flags.ValidateOutput(f.Output); err != nil {
		return nil, err
	}

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

func PrintIsolateKeyResponse(w io.Writer, resp *wallet.IsolateKeyResponse) {
	p := printer.NewInteractivePrinter(w)

	p.CheckMark().Text("Key pair has been isolated in wallet ").Bold(resp.Wallet).Text(" at: ").SuccessText(resp.FilePath).NextLine()
	p.CheckMark().SuccessText("Key isolation succeeded").NextLine()
}
