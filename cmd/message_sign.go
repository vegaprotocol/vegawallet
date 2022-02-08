package cmd

import (
	"encoding/base64"
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
	signMessageLong = cli.LongDesc(`
		Sign any message using a Vega wallet key.
	`)

	signMessageExample = cli.Examples(`
		# Sign a message
		vegawallet message sign --message MESSAGE --wallet WALLET --pubkey PUBKEY
	`)
)

type SignMessageHandler func(*wallet.SignMessageRequest) (*wallet.SignMessageResponse, error)

func NewCmdSignMessage(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func(req *wallet.SignMessageRequest) (*wallet.SignMessageResponse, error) {
		s, err := wallets.InitialiseStore(rf.Home)
		if err != nil {
			return nil, fmt.Errorf("couldn't initialise wallets store: %w", err)
		}

		return wallet.SignMessage(s, req)
	}
	return BuildCmdSignMessage(w, h, rf)
}

func BuildCmdSignMessage(w io.Writer, handler SignMessageHandler, rf *RootFlags) *cobra.Command {
	f := &SignMessageFlags{}

	cmd := &cobra.Command{
		Use:     "sign",
		Short:   "Sign a message using a Vega wallet key",
		Long:    signMessageLong,
		Example: signMessageExample,
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
				PrintSignMessageResponse(w, resp)
			case flags.JSONOutput:
				return printer.FprintJSON(w, struct {
					Signature string `json:"signature"`
				}{
					Signature: resp.Base64,
				})
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
		"Public key to use to the sign the message (hex-encoded)",
	)
	cmd.Flags().StringVarP(&f.Message,
		"message", "m",
		"",
		"Message to be verified (base64-encoded)",
	)
	cmd.Flags().StringVarP(&f.PassphraseFile,
		"passphrase-file", "p",
		"",
		"Path to the file containing the wallet's passphrase",
	)

	autoCompleteWallet(cmd, rf.Home)

	return cmd
}

type SignMessageFlags struct {
	Wallet         string
	PubKey         string
	Message        string
	PassphraseFile string
}

func (f *SignMessageFlags) Validate() (*wallet.SignMessageRequest, error) {
	req := &wallet.SignMessageRequest{}

	if len(f.Wallet) == 0 {
		return nil, flags.FlagMustBeSpecifiedError("wallet")
	}
	req.Wallet = f.Wallet

	if len(f.PubKey) == 0 {
		return nil, flags.FlagMustBeSpecifiedError("pubkey")
	}
	req.PubKey = f.PubKey

	if len(f.Message) == 0 {
		return nil, flags.FlagMustBeSpecifiedError("message")
	}
	decodedMessage, err := base64.StdEncoding.DecodeString(f.Message)
	if err != nil {
		return nil, flags.MustBase64EncodedError("message")
	}
	req.Message = decodedMessage

	passphrase, err := flags.GetPassphrase(f.PassphraseFile)
	if err != nil {
		return nil, err
	}
	req.Passphrase = passphrase

	return req, nil
}

func PrintSignMessageResponse(w io.Writer, req *wallet.SignMessageResponse) {
	p := printer.NewInteractivePrinter(w)

	p.CheckMark().SuccessText("Message signature successful").NextSection()
	p.Text("Signature (base64-encoded):").NextLine().WarningText(req.Base64).NextSection()

	p.BlueArrow().InfoText("Sign a message").NextLine()
	p.Text("To verify a message, see the following command:").NextSection()
	p.Code(fmt.Sprintf("%s verify --help", os.Args[0])).NextSection()
}
