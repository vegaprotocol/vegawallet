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
	verifyArgs struct {
		sig     string
		message string
		pubkey  string
	}

	verifyCmd = &cobra.Command{
		Use:   "verify",
		Short: "Verify the signature",
		Long:  "Verify the signature for a blob of data",
		RunE:  runVerify,
	}
)

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().StringVarP(&verifyArgs.message, "message", "m", "", "Message to be verified (base64)")
	verifyCmd.Flags().StringVarP(&verifyArgs.sig, "signature", "s", "", "Signature to be verified (base64)")
	verifyCmd.Flags().StringVarP(&verifyArgs.pubkey, "pubkey", "k", "", "Public key to be used (hex)")
}

func runVerify(_ *cobra.Command, _ []string) error {
	store, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	handler := wallets.NewHandler(store)

	if len(verifyArgs.pubkey) == 0 {
		return errors.New("pubkey is required")
	}
	if len(verifyArgs.message) == 0 {
		return errors.New("message is required")
	}
	decodedMessage, err := base64.StdEncoding.DecodeString(verifyArgs.message)
	if err != nil {
		return errors.New("message should be encoded into base64")
	}
	if len(verifyArgs.sig) == 0 {
		return errors.New("signature is required")
	}
	decodedSig, err := base64.StdEncoding.DecodeString(verifyArgs.sig)
	if err != nil {
		return errors.New("signature should be encoded into base64")
	}

	isValid, err := handler.VerifyAny(decodedMessage, decodedSig, verifyArgs.pubkey)
	if err != nil {
		return fmt.Errorf("could not verify the message: %w", err)
	}

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		if isValid {
			p.CheckMark().SuccessText("Valid signature").NJump(2)
		} else {
			p.CrossMark().DangerText("Invalid signature").NJump(2)
		}

		p.BlueArrow().InfoText("Sign a message").Jump()
		p.Text("To sign a base-64 encoded message, use the following commands:").NJump(2)
		p.Code(fmt.Sprintf("%s sign --wallet \"YOUR_NAME\" --pubkey %s --message \"YOUR_MESSAGE\"", os.Args[0], verifyArgs.pubkey)).NJump(2)
		p.Text("For more information, use ").Bold("--help").Text(" flag.").Jump()
	} else if rootArgs.output == "json" {
		return printVerifyJson(isValid)
	}

	return nil
}

func printVerifyJson(isValid bool) error {
	return vgjson.Print(struct {
		IsValid bool `json:"isValid"`
	}{
		IsValid: isValid,
	})
}
