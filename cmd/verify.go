package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	verifyArgs struct {
		name       string
		passphrase string
		sig        string
		message    string
		pubkey     string
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
	verifyCmd.Flags().StringVarP(&verifyArgs.name, "name", "n", "", "Name of the wallet to use")
	verifyCmd.Flags().StringVarP(&verifyArgs.passphrase, "passphrase", "p", "", "Passphrase to access the wallet")
	verifyCmd.Flags().StringVarP(&verifyArgs.message, "message", "m", "", "Message to be verified (base64)")
	verifyCmd.Flags().StringVarP(&verifyArgs.sig, "signature", "s", "", "Signature to be verified (base64)")
	verifyCmd.Flags().StringVarP(&verifyArgs.pubkey, "pubkey", "k", "", "Public key to be used (hex)")
}

func runVerify(cmd *cobra.Command, args []string) error {
	handler, err := newWalletHandler(rootArgs.rootPath)
	if err != nil {
		return err
	}

	if len(verifyArgs.name) == 0 {
		return errors.New("wallet name is required")
	}
	if len(verifyArgs.pubkey) == 0 {
		return errors.New("pubkey is required")
	}
	if len(verifyArgs.message) == 0 {
		return errors.New("message is required")
	}
	if len(verifyArgs.sig) == 0 {
		return errors.New("data is required")
	}
	if len(verifyArgs.passphrase) == 0 {
		var err error
		verifyArgs.passphrase, err = promptForPassphrase()
		if err != nil {
			return fmt.Errorf("could not get passphrase: %v", err)
		}
	}

	err = handler.LoginWallet(verifyArgs.name, verifyArgs.passphrase)
	if err != nil {
		return fmt.Errorf("could not login to the wallet: %v", err)
	}

	verified, err := handler.VerifyAny(verifyArgs.name, verifyArgs.message, verifyArgs.sig, verifyArgs.pubkey)
	if err != nil {
		return fmt.Errorf("could not verify the message: %v", err)
	}

	fmt.Printf("%v\n", verified)

	return nil
}
