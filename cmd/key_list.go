package cmd

import (
	"errors"
	"fmt"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
	vgjson "code.vegaprotocol.io/go-wallet/libs/json"
	"code.vegaprotocol.io/go-wallet/wallet"
	"github.com/spf13/cobra"
)

var (
	keyListArgs struct {
		name           string
		passphraseFile string
	}

	keyListCmd = &cobra.Command{
		Use:   "list",
		Short: "List keys of a wallet",
		Long:  "List all the keys for a given wallet",
		RunE:  runKeyList,
	}
)

func init() {
	keyCmd.AddCommand(keyListCmd)
	keyListCmd.Flags().StringVarP(&keyListArgs.name, "name", "n", "", "Name of the wallet to use")
	keyListCmd.Flags().StringVarP(&keyListArgs.passphraseFile, "passphrase-file", "p", "", "Path of the file containing the passphrase to access the wallet")
}

func runKeyList(_ *cobra.Command, _ []string) error {
	handler, err := newWalletHandler(rootArgs.vegaHome)
	if err != nil {
		return err
	}

	if len(keyListArgs.name) == 0 {
		return errors.New("wallet name is required")
	}

	passphrase, err := getPassphrase(keyListArgs.passphraseFile, false)
	if err != nil {
		return err
	}

	err = handler.LoginWallet(keyListArgs.name, passphrase)
	if err != nil {
		return fmt.Errorf("could not login to the wallet: %w", err)
	}

	keys, err := handler.ListKeyPairs(keyListArgs.name)
	if err != nil {
		return fmt.Errorf("could not list the public keys: %w", err)
	}

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		for i, keyPair := range keys {
			p.InfoText(fmt.Sprintf("# Key %d", i+1)).Jump()
			printKeyPair(p, keyPair)
			p.Jump()
		}
	} else if rootArgs.output == "json" {
		return printJsonKeyPairs(keys)
	}

	return nil
}

func printJsonKeyPairs(keys []wallet.KeyPair) error {
	var result []struct {
		PrivateKey       string
		PublicKey        string
		AlgorithmName    string
		AlgorithmVersion uint32
		Meta             []wallet.Meta
	}

	for _, keyPair := range keys {
		result = append(result,
			struct {
				PrivateKey       string
				PublicKey        string
				AlgorithmName    string
				AlgorithmVersion uint32
				Meta             []wallet.Meta
			}{
				PrivateKey:       keyPair.PrivateKey(),
				PublicKey:        keyPair.PublicKey(),
				AlgorithmName:    keyPair.AlgorithmName(),
				AlgorithmVersion: keyPair.AlgorithmVersion(),
				Meta:             keyPair.Meta(),
			})
	}

	return vgjson.Print(keys)
}
