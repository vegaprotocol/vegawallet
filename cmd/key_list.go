package cmd

import (
	"fmt"

	vgjson "code.vegaprotocol.io/shared/libs/json"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/wallet"
	"code.vegaprotocol.io/vegawallet/wallets"
	"github.com/spf13/cobra"
)

var (
	keyListArgs struct {
		wallet         string
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
	keyListCmd.Flags().StringVarP(&keyListArgs.wallet, "wallet", "w", "", "Name of the wallet to use")
	keyListCmd.Flags().StringVarP(&keyListArgs.passphraseFile, "passphrase-file", "p", "", "Path of the file containing the passphrase to access the wallet")
	_ = keyListCmd.MarkFlagRequired("wallet")
}

func runKeyList(_ *cobra.Command, _ []string) error {
	store, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	handler := wallets.NewHandler(store)

	passphrase, err := getPassphrase(keyListArgs.passphraseFile, false)
	if err != nil {
		return err
	}

	err = handler.LoginWallet(keyListArgs.wallet, passphrase)
	if err != nil {
		return fmt.Errorf("could not login to the wallet: %w", err)
	}

	keys, err := handler.ListKeyPairs(keyListArgs.wallet)
	if err != nil {
		return fmt.Errorf("could not list the public keys: %w", err)
	}

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		for i, keyPair := range keys {
			p.InfoText(fmt.Sprintf("# Key %d", i+1)).NextLine()
			printKeyPair(p, keyPair)
			p.NextLine()
		}
	} else if rootArgs.output == "json" {
		return printJSONKeyPairs(keys)
	}

	return nil
}

func printJSONKeyPairs(keys []wallet.KeyPair) error {
	result := make([]keyDescriptionKeyJSON, 0, len(keys))

	for _, keyPair := range keys {
		result = append(result,
			keyDescriptionKeyJSON{
				KeyPair: keyDescriptionKeyPairJSON{
					PrivateKey: keyPair.PrivateKey(),
					PublicKey:  keyPair.PublicKey(),
				},
				Algorithm: keyDescriptionAlgorithmJSON{
					Name:    keyPair.AlgorithmName(),
					Version: keyPair.AlgorithmVersion(),
				},
				Meta:    keyPair.Meta(),
				Tainted: keyPair.IsTainted(),
			},
		)
	}

	return vgjson.Print(result)
}

type keyDescriptionKeyJSON struct {
	KeyPair   keyDescriptionKeyPairJSON   `json:"keyPair"`
	Algorithm keyDescriptionAlgorithmJSON `json:"algorithm"`
	Meta      []wallet.Meta               `json:"meta"`
	Tainted   bool                        `json:"tainted"`
}

type keyDescriptionKeyPairJSON struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

type keyDescriptionAlgorithmJSON struct {
	Name    string `json:"name"`
	Version uint32 `json:"version"`
}
