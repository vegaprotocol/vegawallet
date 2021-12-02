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
	createWalletLong = cli.LongDesc(`
		Create a wallet and generate the first Ed25519 key pair.
	`)

	createWalletExample = cli.Examples(`
		# Creating a wallet
		vegawallet create --wallet WALLET
	`)
)

type CreateWalletHandler func(*wallet.CreateWalletRequest) (*wallet.CreateWalletResponse, error)

func NewCmdCreateWallet(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func(req *wallet.CreateWalletRequest) (*wallet.CreateWalletResponse, error) {
		s, err := wallets.InitialiseStore(rf.Home)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialise wallets store: %w", err)
		}

		return wallet.CreateWallet(s, req)
	}

	return BuildCmdCreateWallet(w, h, rf)
}

func BuildCmdCreateWallet(w io.Writer, handler CreateWalletHandler, rf *RootFlags) *cobra.Command {
	f := &CreateWalletFlags{}

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a wallet",
		Long:    createWalletLong,
		Example: createWalletExample,
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
				PrintCreateWalletResponse(w, resp)
			case flags.JSONOutput:
				return printer.FprintJSON(w, resp)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&f.Wallet,
		"wallet", "w",
		"",
		"The wallet where the key is generated in",
	)
	cmd.Flags().StringVarP(&f.PassphraseFile,
		"passphrase-file", "p",
		"",
		"Path to the file containing the wallet's passphrase",
	)

	return cmd
}

type CreateWalletFlags struct {
	Wallet         string
	PassphraseFile string
}

func (f *CreateWalletFlags) Validate() (*wallet.CreateWalletRequest, error) {
	req := &wallet.CreateWalletRequest{}

	if len(f.Wallet) == 0 {
		return nil, flags.FlagMustBeSpecifiedError("wallet")
	}
	req.Wallet = f.Wallet

	passphrase, err := flags.GetConfirmedPassphrase(f.PassphraseFile)
	if err != nil {
		return nil, err
	}
	req.Passphrase = passphrase

	return req, nil
}

func PrintCreateWalletResponse(w io.Writer, resp *wallet.CreateWalletResponse) {
	p := printer.NewInteractivePrinter(w)

	p.CheckMark().Text("Wallet ").Bold(resp.Wallet.Name).Text(" has been created at: ").SuccessText(resp.Wallet.FilePath).NextLine()
	p.CheckMark().Text("First key pair has been generated for wallet ").Bold(resp.Wallet.Name).Text(" at: ").SuccessText(resp.Wallet.FilePath).NextLine()
	p.CheckMark().SuccessText("Creating wallet succeeded").NextSection()

	p.Text("Wallet mnemonic:").NextLine()
	p.WarningText(resp.Wallet.Mnemonic).NextLine()
	p.Text("Wallet version:").NextLine()
	p.WarningText(fmt.Sprintf("%d", resp.Wallet.Version)).NextLine()
	p.Text("First public key:").NextLine()
	p.WarningText(resp.Key.PublicKey).NextLine()
	p.NextSection()

	p.RedArrow().DangerText("Important").NextLine()
	p.Text("Write down the ").Bold("mnemonic").Text(" and the ").Bold("wallet's version").Text(", and store it somewhere safe and secure, now.").NextLine()
	p.DangerText("The mnemonic will not be displayed ever again, nor will you be able to retrieve it!").NextSection()

	p.BlueArrow().InfoText("Run the service").NextLine()
	p.Text("Now, you can run the service. See the following command:").NextSection()
	p.Code(fmt.Sprintf("%s service run --help", os.Args[0])).NextSection()
}
