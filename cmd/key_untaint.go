package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	keyUntaintArgs struct {
		name           string
		passphrase     string
		passphraseFile string
		pubKey         string
	}

	keyUntaintCmd = &cobra.Command{
		Use:   "untaint",
		Short: "Untaint a public key",
		Long:  "Untaint a public key",
		RunE:  runKeyUntaint,
	}
)

func init() {
	keyCmd.AddCommand(keyUntaintCmd)
	keyUntaintCmd.Flags().StringVarP(&keyUntaintArgs.name, "name", "n", "", "Name of the wallet to use")
	keyUntaintCmd.Flags().StringVar(&keyUntaintArgs.passphrase, "passphrase", "", "Passphrase to access the wallet")
	keyUntaintCmd.Flags().StringVar(&keyUntaintArgs.passphraseFile, "passphrase-file", "", "Path of the file containing the passphrase to access the wallet")
	keyUntaintCmd.Flags().StringVarP(&keyUntaintArgs.pubKey, "pubkey", "k", "", "Public key to be used (hex)")
}

func runKeyUntaint(cmd *cobra.Command, args []string) error {
	handler, err := newWalletHandler(rootArgs.rootPath)
	if err != nil {
		return err
	}

	if len(keyUntaintArgs.name) == 0 {
		return errors.New("wallet name is required")
	}

	passphrase, err := getPassphrase(keyUntaintArgs.passphrase, keyUntaintArgs.passphraseFile, false)
	if err != nil {
		return err
	}

	err = handler.UntaintKey(keyUntaintArgs.name, keyUntaintArgs.pubKey, passphrase)
	if err != nil {
		return fmt.Errorf("could not untaint the key: %v", err)
	}

	fmt.Printf("The key has been untainted.\n")
	return nil
}
