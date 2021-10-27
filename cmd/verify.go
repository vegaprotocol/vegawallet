package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	vgjson "code.vegaprotocol.io/shared/libs/json"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/wallets"
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

	ErrMessageShouldBeBase64   = errors.New("message should be encoded into base64")
	ErrSignatureShouldBeBase64 = errors.New("signature should be encoded into base64")
)

func init() {
	rootCmd.AddCommand(verifyCmd)
	verifyCmd.Flags().StringVarP(&verifyArgs.message, "message", "m", "", "Message to be verified (base64)")
	verifyCmd.Flags().StringVarP(&verifyArgs.sig, "signature", "s", "", "Signature to be verified (base64)")
	verifyCmd.Flags().StringVarP(&verifyArgs.pubkey, "pubkey", "k", "", "Public key to be used (hex)")
	_ = verifyCmd.MarkFlagRequired("pubkey")
	_ = verifyCmd.MarkFlagRequired("message")
	_ = verifyCmd.MarkFlagRequired("signature")
}

func runVerify(_ *cobra.Command, _ []string) error {
	store, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	handler := wallets.NewHandler(store)

	decodedMessage, err := base64.StdEncoding.DecodeString(verifyArgs.message)
	if err != nil {
		return ErrMessageShouldBeBase64
	}

	decodedSig, err := base64.StdEncoding.DecodeString(verifyArgs.sig)
	if err != nil {
		return ErrSignatureShouldBeBase64
	}

	isValid, err := handler.VerifyAny(decodedMessage, decodedSig, verifyArgs.pubkey)
	if err != nil {
		return fmt.Errorf("could not verify the message: %w", err)
	}

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		if isValid {
			p.CheckMark().SuccessText("Valid signature").NextSection()
		} else {
			p.CrossMark().DangerText("Invalid signature").NextSection()
		}

		p.BlueArrow().InfoText("Sign a message").NextLine()
		p.Text("To sign a base-64 encoded message, use the following commands:").NextSection()
		p.Code(fmt.Sprintf("%s sign --wallet \"YOUR_NAME\" --pubkey %s --message \"YOUR_MESSAGE\"", os.Args[0], verifyArgs.pubkey)).NextSection()
		p.Text("For more information, use ").Bold("--help").Text(" flag.").NextLine()
	} else if rootArgs.output == "json" {
		return printVerifyJSON(isValid)
	}

	return nil
}

func printVerifyJSON(isValid bool) error {
	return vgjson.Print(struct {
		IsValid bool `json:"isValid"`
	}{
		IsValid: isValid,
	})
}
