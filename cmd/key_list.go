package cmd

import (
	"errors"
	"fmt"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/wallets"
	vgjson "code.vegaprotocol.io/shared/libs/json"
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
	store, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	handler := wallets.NewHandler(store)

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
	var result []keyGenerateKeyJson

	for _, keyPair := range keys {
		result = append(result,
			keyGenerateKeyJson{
				KeyPair: keyGenerateKeyPairJson{
					PrivateKey: keyPair.PrivateKey(),
					PublicKey:  keyPair.PublicKey(),
				},
				Algorithm: keyGenerateAlgorithmJson{
					Name:    keyPair.AlgorithmName(),
					Version: keyPair.AlgorithmVersion(),
				},
				Meta: keyPair.Meta(),
			},
		)
	}

	return vgjson.Print(keys)
}
