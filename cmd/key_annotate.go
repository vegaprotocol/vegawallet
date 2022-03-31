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
	annotateKeyLong = cli.LongDesc(`
		Add metadata to a key pair. All existing metadata is removed and replaced
		by the specified new metadata.

		The metadata is a list of key-value pairs. A key-value is colon-separated, 
		and the key-values are comma-separated.

		It is possible to give a name to a key pair, that is recognised by user
		interfaces, by setting the metadata "name".
	`)

	annotateKeyExample = cli.Examples(`
		Given the following metadata to be added:
			- name: my-wallet
			- type: validation

		# Annotate a key pair
		vegawallet key annotate --wallet WALLET --pubkey PUBKEY --meta "name:my-wallet,type:validation"

		# Remove all metadata from a key pair
		vegawallet key annotate --wallet WALLET --pubkey PUBKEY --clear
	`)
)

type AnnotateKeyHandler func(*wallet.AnnotateKeyRequest) error

func NewCmdAnnotateKey(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func(req *wallet.AnnotateKeyRequest) error {
		s, err := wallets.InitialiseStore(rf.Home)
		if err != nil {
			return fmt.Errorf("couldn't initialise wallets store: %w", err)
		}

		return wallet.AnnotateKey(s, req)
	}

	return BuildCmdAnnotateKey(w, h, rf)
}

func BuildCmdAnnotateKey(w io.Writer, handler AnnotateKeyHandler, rf *RootFlags) *cobra.Command {
	f := AnnotateKeyFlags{}

	cmd := &cobra.Command{
		Use:     "annotate",
		Short:   "Add metadata to a key pair",
		Long:    annotateKeyLong,
		Example: annotateKeyExample,
		RunE: func(_ *cobra.Command, _ []string) error {
			req, err := f.Validate()
			if err != nil {
				return err
			}

			if err := handler(req); err != nil {
				return err
			}

			switch f.Output {
			case flags.InteractiveOutput:
				PrintAnnotateKeyResponse(w, req)
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
	cmd.Flags().StringVarP(&f.PassphraseFile,
		"passphrase-file", "p",
		"",
		"Path to the file containing the wallet's passphrase",
	)
	cmd.Flags().StringVarP(&f.PubKey,
		"pubkey", "k",
		"",
		"Public key to annotate (hex-encoded)",
	)
	cmd.Flags().StringSliceVarP(&f.RawMetadata,
		"meta", "m",
		[]string{},
		`A list of metadata e.g: "my-key1:my-value1,my-key2:my-value2"`,
	)
	cmd.Flags().BoolVar(&f.Clear,
		"clear",
		false,
		"Clear the metadata",
	)

	addOutputFlag(cmd, &f.Output)

	autoCompleteWallet(cmd, rf.Home)

	return cmd
}

type AnnotateKeyFlags struct {
	Wallet         string
	PubKey         string
	PassphraseFile string
	Clear          bool
	RawMetadata    []string
	Output         string
}

func (f *AnnotateKeyFlags) Validate() (*wallet.AnnotateKeyRequest, error) {
	req := &wallet.AnnotateKeyRequest{}

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

	if len(f.RawMetadata) == 0 && !f.Clear {
		return nil, flags.OneOfFlagsMustBeSpecifiedError("meta", "clear")
	}
	if len(f.RawMetadata) != 0 && f.Clear {
		return nil, flags.FlagsMutuallyExclusiveError("meta", "clear")
	}

	metadata, err := cli.ParseMetadata(f.RawMetadata)
	if err != nil {
		return nil, err
	}
	req.Metadata = metadata

	passphrase, err := flags.GetPassphrase(f.PassphraseFile)
	if err != nil {
		return nil, err
	}
	req.Passphrase = passphrase

	return req, nil
}

func PrintAnnotateKeyResponse(w io.Writer, req *wallet.AnnotateKeyRequest) {
	metadataHaveBeenCleared := len(req.Metadata) == 0

	p := printer.NewInteractivePrinter(w)
	if metadataHaveBeenCleared {
		p.CheckMark().SuccessText("Annotation cleared").NextLine()
		return
	}
	p.CheckMark().SuccessText("Annotation succeeded").NextSection()
	p.Text("New metadata:").NextLine()
	printMeta(p, req.Metadata)
}
