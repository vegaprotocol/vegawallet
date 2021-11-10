package cmd

import (
	"fmt"
	"io"

	"code.vegaprotocol.io/vegawallet/cmd/cli"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/wallet"
	"code.vegaprotocol.io/vegawallet/wallets"
	"github.com/spf13/cobra"
)

var (
	rotateKeyLong = cli.LongDesc(`
		Get signed key rotation transaction.
		This transaction can be applied to Vega protocol trough wallet command.
		A new public key to rotate to is selected.
		
		The transaction is outputted as a base64 string.
	`)

	rotateKeyExample = cli.Examples(`
		Given that a new public key NEW_PUBLIC_KEY has been previously generated in the wallet:

		# Get signed transaction for rotating new key
		vegawallet key rotate --wallet WALLET --newpubkey NEW_PUBLIC_KEY
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
		RunE: func(_ *cobra.Command, _ []string) error {
			req, err := f.Validate()
			if err != nil {
				return err
			}

			res, err := handler(req)
			if err != nil {
				return err
			}

			switch rf.Output {
			case flags.InteractiveOutput:
				PrintRotateKeyResponse(w, res)
			case flags.JSONOutput:
				return nil
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&f.Wallet,
		"wallet", "w",
		"",
		"Wallet holding the master public key",
	)
	cmd.Flags().StringVarP(&f.PassphraseFile,
		"passphrase-file", "p",
		"",
		"Path to the file containing the wallet's passphrase",
	)
	cmd.Flags().StringVarP(&f.NewPubKey,
		"new-pubkey", "nk",
		"",
		"New public key to rotate to",
	)

	return cmd
}

type RotateKeyFlags struct {
	Wallet         string
	NewPubKey      string
	PassphraseFile string
}

func (f *RotateKeyFlags) Validate() (*wallet.RotateKeyRequest, error) {
	req := &wallet.RotateKeyRequest{}

	if len(f.Wallet) == 0 {
		return nil, flags.FlagMustBeSpecifiedError("wallet")
	}
	req.Wallet = f.Wallet

	if len(f.NewPubKey) == 0 {
		return nil, flags.FlagMustBeSpecifiedError("pubkey")
	}
	req.NewPubKey = f.NewPubKey

	passphrase, err := flags.GetPassphrase(f.PassphraseFile)
	if err != nil {
		return nil, err
	}
	req.Passphrase = passphrase

	return req, nil
}

func PrintRotateKeyResponse(w io.Writer, req *wallet.RotateKeyResponse) {
	// @TODO implement printer

	// metadataHaveBeenCleared := len(req.Metadata) == 0

	// p := printer.NewInteractivePrinter(w)
	// if metadataHaveBeenCleared {
	// 	p.CheckMark().SuccessText("Annotation cleared").NextLine()
	// 	return
	// }
	// p.CheckMark().SuccessText("Annotation succeeded").NextSection()
	// p.Text("New metadata:").NextLine()
	// printMeta(p, req.Metadata)
}
