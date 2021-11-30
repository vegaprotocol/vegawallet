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
	generateKeyLong = cli.LongDesc(`
		Generate a new Ed25519 key pair in a given wallet.

		If the targeted wallet doesn't exist, it will be automatically generated.
	`)

	generateKeyExample = cli.Examples(`
		# Generate a key pair
		vegawallet key generate --wallet WALLET

		# Generate a key pair with additional metadata (name = my-wallet and type = validation)
		vegawallet key generate --wallet WALLET --meta "name:my-wallet,type:validation"

		# Generate a key pair with custom name
		vegawallet key generate --wallet WALLET --meta "name:my-wallet"
	`)
)

type GenerateKeyHandler func(flags.PassphraseGetterWithOps, *wallet.GenerateKeyRequest) (*wallet.GenerateKeyResponse, error)

func NewCmdGenerateKey(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func(passphraseGetter flags.PassphraseGetterWithOps, req *wallet.GenerateKeyRequest) (*wallet.GenerateKeyResponse, error) {
		s, err := wallets.InitialiseStore(rf.Home)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialise wallets store: %w", err)
		}

		// Because the passphrase needs to be retrieved based on the state of the
		// specified wallet, we get it here and fill the request with it.
		passphrase, err := passphraseGetter(!s.WalletExists(req.Wallet))
		if err != nil {
			return nil, err
		}
		req.Passphrase = passphrase

		return wallet.GenerateKey(s, req)
	}

	return BuildCmdGenerateKey(w, h, rf)
}

func BuildCmdGenerateKey(w io.Writer, handler GenerateKeyHandler, rf *RootFlags) *cobra.Command {
	f := &GenerateKeyFlags{}

	cmd := &cobra.Command{
		Use:     "generate",
		Short:   "Generate a new key pair in a given wallet",
		Long:    generateKeyLong,
		Example: generateKeyExample,
		RunE: func(_ *cobra.Command, _ []string) error {
			req, err := f.Validate()
			if err != nil {
				return err
			}

			// If the specified wallet doesn't exist yet, we need to ask for
			// passphrase confirmation. However, this information is only available
			// to the handler, which means we can't retrieve the passphrase during
			// the flag validation step.
			//
			// As a result, we need to delegate the retrieval of the passphrase
			// to the handler. This is why we build a function that takes care
			// of this task based on the flags set, and pass it to the handler.
			//
			// With this method, the handler can get the passphrase in isolation,
			// without knowledge of the command line flags nor the retrieval
			// process.
			pg := flags.BuildPassphraseGetterWithOps(f.PassphraseFile)

			resp, err := handler(pg, req)
			if err != nil {
				return err
			}

			switch rf.Output {
			case flags.InteractiveOutput:
				PrintGenerateKeyResponse(w, resp)
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
	cmd.Flags().StringSliceVarP(&f.RawMetadata,
		"meta", "m",
		[]string{},
		`Metadata to add to the generated key-pair: "my-key1:my-value1,my-key2:my-value2"`,
	)

	autoCompleteWallet(cmd, rf.Home)

	return cmd
}

type GenerateKeyFlags struct {
	Wallet         string
	PassphraseFile string
	RawMetadata    []string
}

func (f *GenerateKeyFlags) Validate() (*wallet.GenerateKeyRequest, error) {
	req := &wallet.GenerateKeyRequest{}

	if len(f.Wallet) == 0 {
		return nil, flags.FlagMustBeSpecifiedError("wallet")
	}
	req.Wallet = f.Wallet

	metadata, err := cli.ParseMetadata(f.RawMetadata)
	if err != nil {
		return nil, err
	}
	req.Metadata = metadata

	return req, nil
}

func PrintGenerateKeyResponse(w io.Writer, resp *wallet.GenerateKeyResponse) {
	walletHasBeenCreated := len(resp.Wallet.Mnemonic) != 0

	p := printer.NewInteractivePrinter(w)

	if walletHasBeenCreated {
		p.CheckMark().Text("Wallet ").Bold(resp.Wallet.Name).Text(" has been created at: ").SuccessText(resp.Wallet.FilePath).NextLine()
	}
	p.CheckMark().Text("Key pair has been generated for wallet ").Bold(resp.Wallet.Name).Text(" at: ").SuccessText(resp.Wallet.FilePath).NextLine()
	p.CheckMark().SuccessText("Generating a key pair succeeded").NextSection()

	if walletHasBeenCreated {
		p.Text("Wallet mnemonic:").NextLine()
		p.WarningText(resp.Wallet.Mnemonic).NextLine()
		p.Text("Wallet version:").NextLine()
		p.WarningText(fmt.Sprintf("%d", resp.Wallet.Version)).NextLine()
	}
	p.Text("Public key:").NextLine()
	p.WarningText(resp.Key.PublicKey).NextLine()
	p.Text("Metadata:").NextLine()
	printMeta(p, resp.Key.Meta)
	p.NextSection()

	if walletHasBeenCreated {
		p.RedArrow().DangerText("Important").NextLine()
		p.Text("Write down the ").Bold("mnemonic").Text(" and the ").Bold("wallet's version").Text(", and store it somewhere safe and secure, now.").NextLine()
		p.DangerText("The mnemonic will not be displayed ever again, nor will you be able to retrieve it!").NextSection()
	}

	p.BlueArrow().InfoText("Run the service").NextLine()
	p.Text("Now, you can run the service. See the following command:").NextSection()
	p.Code(fmt.Sprintf("%s service run --help", os.Args[0])).NextSection()
}
