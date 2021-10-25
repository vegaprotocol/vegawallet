package cmd

import (
	"encoding/base64"
	"fmt"
	"os"

	vgjson "code.vegaprotocol.io/shared/libs/json"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/wallets"
	"github.com/spf13/cobra"
)

var (
	signArgs struct {
		wallet         string
		passphraseFile string
		message        string
		pubKey         string
	}

	// signCmd represents the sign command.
	signCmd = &cobra.Command{
		Use:   "sign",
		Short: "Sign a blob of data",
		Long:  "Sign a blob of dara base64 encoded",
		RunE:  runSign,
	}
)

func init() {
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().StringVarP(&signArgs.wallet, "wallet", "w", "", "Name of the wallet to use")
	signCmd.Flags().StringVarP(&signArgs.passphraseFile, "passphrase-file", "p", "", "Path of the file containing the passphrase to access the wallet")
	signCmd.Flags().StringVarP(&signArgs.message, "message", "m", "", "Message to be signed (base64)")
	signCmd.Flags().StringVarP(&signArgs.pubKey, "pubkey", "k", "", "Public key to be used (hex)")
	_ = signCmd.MarkFlagRequired("network")
	_ = signCmd.MarkFlagRequired("pubkey")
	_ = signCmd.MarkFlagRequired("message")
}

func runSign(_ *cobra.Command, _ []string) error {
	store, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	handler := wallets.NewHandler(store)

	decodedMessage, err := base64.StdEncoding.DecodeString(signArgs.message)
	if err != nil {
		return ErrMessageShouldBeBase64
	}

	passphrase, err := getPassphrase(signArgs.passphraseFile, false)
	if err != nil {
		return err
	}

	err = handler.LoginWallet(signArgs.wallet, passphrase)
	if err != nil {
		return fmt.Errorf("could not login to the wallet: %w", err)
	}

	sig, err := handler.SignAny(signArgs.wallet, decodedMessage, signArgs.pubKey)
	if err != nil {
		return err
	}

	encodedSig := base64.StdEncoding.EncodeToString(sig)

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		p.CheckMark().SuccessText("Message signature successful").NextSection()
		p.Text("Signature (base64):").NextLine().WarningText(encodedSig).NextSection()

		p.BlueArrow().InfoText("Verify a signature").NextLine()
		p.Text("To verify a base-64 encoded message, use the following commands:").NextSection()
		p.Code(fmt.Sprintf("%s verify --pubkey %s --message \"%s\" --signature %s", os.Args[0], signArgs.pubKey, signArgs.message, encodedSig)).NextSection()
		p.Text("For more information, use ").Bold("--help").Text(" flag.").NextLine()
	} else if rootArgs.output == "json" {
		return printSignJSON(encodedSig)
	}

	return nil
}

func printSignJSON(sig string) error {
	return vgjson.Print(struct {
		Signature string `json:"signature"`
	}{
		Signature: sig,
	})
}
