package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	signArgs struct {
		name       string
		passphrase string
		message    string
		pubkey     string
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
	signCmd.Flags().StringVarP(&signArgs.passphrase, "passphrase", "p", "", "Passphrase to access the wallet")
	signCmd.Flags().StringVarP(&signArgs.message, "message", "m", "", "Message to be signed (base64)")
	signCmd.Flags().StringVarP(&signArgs.pubkey, "pubkey", "k", "", "Public key to be used (hex)")
}

func runSign(cmd *cobra.Command, args []string) error {
	handler, err := newWalletHandler(rootArgs.rootPath)
	if err != nil {
		return err
	}

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

	if len(signArgs.passphrase) <= 0 {
		var err error
		signArgs.passphrase, err = promptForPassphrase()
		if err != nil {
			return fmt.Errorf("could not get passphrase: %v", err)
		}
	}

	err = handler.LoginWallet(signArgs.name, signArgs.passphrase)
	if err != nil {
		return fmt.Errorf("could not login to the wallet: %v", err)
	}
	sig, err := handler.SignAny(signArgs.name, decodedMessage, signArgs.pubkey)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", base64.StdEncoding.EncodeToString(sig))
	return nil
}
