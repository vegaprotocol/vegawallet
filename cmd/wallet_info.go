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
	infoLong = cli.LongDesc(`
		Get wallet information such as wallet ID, version and type.
	`)

	infoExample = cli.Examples(`
		# Get the wallet information
		vegawallet info --wallet WALLET
	`)
)

type GetInfoWalletHandler func(*wallet.GetWalletInfoRequest) (*wallet.GetWalletInfoResponse, error)

func NewCmdGetInfoWallet(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func(req *wallet.GetWalletInfoRequest) (*wallet.GetWalletInfoResponse, error) {
		s, err := wallets.InitialiseStore(rf.Home)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialise wallets store: %w", err)
		}

		return wallet.GetWalletInfo(s, req)
	}

	return BuildCmdGetInfoWallet(w, h, rf)
}

func BuildCmdGetInfoWallet(w io.Writer, handler GetInfoWalletHandler, rf *RootFlags) *cobra.Command {
	f := &GetWalletInfoFlags{}

	cmd := &cobra.Command{
		Use:     "info",
		Short:   "Get wallet information",
		Long:    infoLong,
		Example: infoExample,
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
				PrintGetWalletInfoResponse(w, resp)
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

	autoCompleteWallet(cmd, rf.Home)

	return cmd
}

type GetWalletInfoFlags struct {
	Wallet         string
	PassphraseFile string
}

func (f *GetWalletInfoFlags) Validate() (*wallet.GetWalletInfoRequest, error) {
	req := &wallet.GetWalletInfoRequest{}

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

func PrintGetWalletInfoResponse(w io.Writer, resp *wallet.GetWalletInfoResponse) {
	p := printer.NewInteractivePrinter(w)

	p.Text("Type:").NextLine().WarningText(resp.Type).NextLine()
	p.Text("Version:").NextLine().WarningText(fmt.Sprintf("%d", resp.Version)).NextLine()
	p.Text("ID:").NextLine().WarningText(resp.ID).NextLine()
}
