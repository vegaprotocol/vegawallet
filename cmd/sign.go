package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
	"code.vegaprotocol.io/go-wallet/wallets"
	vgjson "code.vegaprotocol.io/shared/libs/json"
	"github.com/spf13/cobra"
)

var (
	signArgs struct {
		name           string
		passphraseFile string
		message        string
		pubkey         string
	}

	// signCmd represents the sign command
	signCmd = &cobra.Command{
		Use:   "sign",
		Short: "Sign a blob of data",
		Long:  "Sign a blob of dara base64 encoded",
		RunE:  runSign,
	}
)

func init() {
	rootCmd.AddCommand(signCmd)
	signCmd.Flags().StringVarP(&signArgs.name, "name", "n", "", "Name of the wallet to use")
	signCmd.Flags().StringVarP(&signArgs.passphraseFile, "passphrase-file", "p", "", "Path of the file containing the passphrase to access the wallet")
	signCmd.Flags().StringVarP(&signArgs.message, "message", "m", "", "Message to be signed (base64)")
	signCmd.Flags().StringVarP(&signArgs.pubkey, "pubkey", "k", "", "Public key to be used (hex)")
}

func runSign(_ *cobra.Command, _ []string) error {
	store, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	handler := wallets.NewHandler(store)

	if len(signArgs.name) == 0 {
		return errors.New("wallet name is required")
	}
	if len(signArgs.pubkey) == 0 {
		return errors.New("pubkey is required")
	}
	if len(signArgs.message) == 0 {
		return errors.New("message is required")
	}
	decodedMessage, err := base64.StdEncoding.DecodeString(signArgs.message)
	if err != nil {
		return errors.New("message should be encoded into base64")
	}

	passphrase, err := getPassphrase(signArgs.passphraseFile, false)
	if err != nil {
		return err
	}

	err = handler.LoginWallet(signArgs.name, passphrase)
	if err != nil {
		return fmt.Errorf("could not login to the wallet: %w", err)
	}

	sig, err := handler.SignAny(signArgs.name, decodedMessage, signArgs.pubkey)
	if err != nil {
		return err
	}

	encodedSig := base64.StdEncoding.EncodeToString(sig)

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		p.CheckMark().SuccessText("Message signature successful").NJump(2)
		p.Text("Signature (base64):").Jump().WarningText(encodedSig).NJump(2)

		p.BlueArrow().InfoText("Verify a signature").Jump()
		p.Text("To verify a base-64 encoded message, use the following commands:").NJump(2)
		p.Code(fmt.Sprintf("%s verify --pubkey %s --message \"%s\" --signature %s", os.Args[0], signArgs.pubkey, signArgs.message, encodedSig)).NJump(2)
		p.Text("For more information, use ").Bold("--help").Text(" flag.").Jump()
	} else if rootArgs.output == "json" {
		return printSignJson(encodedSig)
	}

	return nil
}

func printSignJson(sig string) error {
	return vgjson.Print(struct {
		Signature string
	}{
		Signature: sig,
	})
}
