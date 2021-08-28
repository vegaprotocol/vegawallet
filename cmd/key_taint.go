package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	keyTaintArgs struct {
		name           string
		passphrase     string
		passphraseFile string
		pubkey         string
	}

	keyTaintCmd = &cobra.Command{
		Use:   "taint",
		Short: "Taint a public key",
		Long:  "Taint a public key",
		RunE:  runKeyTaint,
	}
)

func init() {
	keyCmd.AddCommand(keyTaintCmd)
	keyTaintCmd.Flags().StringVarP(&keyTaintArgs.name, "name", "n", "", "Name of the wallet to use")
	keyTaintCmd.Flags().StringVar(&keyTaintArgs.passphrase, "passphrase", "", "Passphrase to access the wallet")
	keyTaintCmd.Flags().StringVar(&keyTaintArgs.passphraseFile, "passphrase-file", "", "Path of the file containing the passphrase to access the wallet")
	keyTaintCmd.Flags().StringVarP(&keyTaintArgs.pubkey, "pubkey", "k", "", "Public key to be used (hex)")
}

func runKeyTaint(cmd *cobra.Command, args []string) error {
	handler, err := newWalletHandler(rootArgs.rootPath)
	if err != nil {
		return err
	}

	if len(keyTaintArgs.name) == 0 {
		return errors.New("wallet name is required")
	}

	passphrase, err := getPassphrase(keyTaintArgs.passphrase, keyTaintArgs.passphraseFile, false)
	if err != nil {
		return err
	}

	err = handler.TaintKey(keyTaintArgs.name, keyTaintArgs.pubkey, passphrase)
	if err != nil {
		return fmt.Errorf("could not taint the key: %w", err)
	}

	fmt.Printf("The key has been tainted.\n")
	return nil
}
