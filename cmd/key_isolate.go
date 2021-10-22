package cmd

import (
	"fmt"

	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/wallets"
	vgjson "code.vegaprotocol.io/shared/libs/json"
	"github.com/spf13/cobra"
)

var (
	keyIsolateArgs struct {
		name           string
		passphraseFile string
		pubkey         string
	}

	keyIsolateCmd = &cobra.Command{
		Use:   "isolate",
		Short: "Isolate a wallet with the specified key pair",
		Long:  "Isolate a wallet with the specified key pair, without the root node responsible of key pairs generation. This allows extra layer of security.",
		RunE:  runKeyIsolate,
	}
)

func init() {
	keyCmd.AddCommand(keyIsolateCmd)
	keyIsolateCmd.Flags().StringVarP(&keyIsolateArgs.name, "wallet", "w", "", "Name of the wallet to use")
	keyIsolateCmd.Flags().StringVarP(&keyIsolateArgs.passphraseFile, "passphrase-file", "p", "", "Path of the file containing the passphrase to access the wallet")
	keyIsolateCmd.Flags().StringVarP(&keyIsolateArgs.pubkey, "pubkey", "k", "", "Public key to be used (hex)")
	_ = keyIsolateCmd.MarkFlagRequired("wallet")
	_ = keyIsolateCmd.MarkFlagRequired("pubkey")
}

func runKeyIsolate(_ *cobra.Command, _ []string) error {
	store, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	passphrase, err := getPassphrase(keyIsolateArgs.passphraseFile, false)
	if err != nil {
		return err
	}

	w, err := store.GetWallet(keyIsolateArgs.name, passphrase)
	if err != nil {
		return fmt.Errorf("couldn't get wallet %s: %w", keyIsolateArgs.name, err)
	}

	isolatedWallet, err := w.IsolateWithKey(keyIsolateArgs.pubkey)
	if err != nil {
		return fmt.Errorf("couldn't isolate wallet %s: %w", keyIsolateArgs.name, err)
	}

	if err := store.SaveWallet(isolatedWallet, passphrase); err != nil {
		return fmt.Errorf("couldn't save isolated wallet %s: %w", isolatedWallet.Name(), err)
	}

	walletPath := store.GetWalletPath(isolatedWallet.Name())

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		p.CheckMark().Text("Key pair has been isolated in wallet ").Bold(isolatedWallet.Name()).Text(" at: ").SuccessText(walletPath).Jump()
		p.CheckMark().SuccessText("Key isolation succeeded").NJump(2)
	} else if rootArgs.output == "json" {
		return vgjson.Print(keyIsolateJson{
			Wallet:   isolatedWallet.Name(),
			FilePath: walletPath,
		})
	}
	return nil
}

type keyIsolateJson struct {
	Wallet   string `json:"wallet"`
	FilePath string `json:"filePath"`
}
