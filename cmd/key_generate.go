package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"code.vegaprotocol.io/go-wallet/wallet"
	"github.com/spf13/cobra"
)

var (
	keyGenerateArgs struct {
		name           string
		passphrase     string
		passphraseFile string
		metas          string
	}

	keyGenerateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate a new key pair for a wallet",
		Long:  "Generate a new key pair for a wallet, this will implicitly generate a new wallet if none exist for the given name",
		RunE:  runKeyGenerate,
	}
)

func init() {
	keyCmd.AddCommand(keyGenerateCmd)
	keyGenerateCmd.Flags().StringVarP(&keyGenerateArgs.name, "name", "n", "", "Name of the wallet to use")
	keyGenerateCmd.Flags().StringVarP(&keyGenerateArgs.passphrase, "passphrase", "p", "", "Passphrase to access the wallet")
	keyGenerateCmd.Flags().StringVar(&keyGenerateArgs.passphraseFile, "passphrase-file", "", "Path of the file containing the passphrase to access the wallet")
	keyGenerateCmd.Flags().StringVarP(&keyGenerateArgs.metas, "meta", "m", "", `A list of metadata e.g: "primary:true;asset:BTC"`)
}

func runKeyGenerate(cmd *cobra.Command, args []string) error {
	handler, err := newWalletHandler(rootArgs.rootPath)
	if err != nil {
		return err
	}

	if len(keyGenerateArgs.name) == 0 {
		return errors.New("wallet name is required")
	}

	walletExists := handler.WalletExists(keyGenerateArgs.name)

	passphrase, err := getPassphrase(keyGenerateArgs.passphrase, keyGenerateArgs.passphraseFile, !walletExists)
	if err != nil {
		return err
	}

	metas, err := parseMeta(keyGenerateArgs.metas)
	if err != nil {
		return err
	}

	if !walletExists {
		mnemonic, err := handler.CreateWallet(keyGenerateArgs.name, passphrase)
		if err != nil {
			return fmt.Errorf("couldn't create wallet: %v", err)
		}
		fmt.Printf("new mnemonic:\n")
		fmt.Printf("%s\n", mnemonic)
	}

	keyPair, err := handler.GenerateKeyPair(keyGenerateArgs.name, passphrase)
	if err != nil {
		return fmt.Errorf("could not generate a key pair: %v", err)
	}

	err = handler.UpdateMeta(keyGenerateArgs.name, keyPair.PublicKey(), passphrase, metas)
	if err != nil {
		return fmt.Errorf("could not update the meta: %v", err)
	}

	buf, err := json.MarshalIndent(keyPair, " ", " ")
	if err != nil {
		return fmt.Errorf("unable to marshal message: %v", err)
	}
	fmt.Printf("new generated keys:\n")
	fmt.Printf("%s\n", string(buf))

	return nil
}

func parseMeta(metaStr string) ([]wallet.Meta, error) {
	var metas []wallet.Meta

	if len(metaStr) == 0 {
		return metas, nil
	}

	rawMetas := strings.Split(metaStr, ";")
	for _, v := range rawMetas {
		rawMeta := strings.Split(v, ":")
		if len(rawMeta) != 2 {
			return nil, fmt.Errorf("invalid meta format")
		}
		metas = append(metas, wallet.Meta{Key: rawMeta[0], Value: rawMeta[1]})
	}

	return metas, nil
}
