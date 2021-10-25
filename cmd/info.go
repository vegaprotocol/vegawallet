package cmd

import (
	"fmt"

	vgjson "code.vegaprotocol.io/shared/libs/json"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/wallets"
	"github.com/spf13/cobra"
)

var (
	infoArgs struct {
		wallet         string
		passphraseFile string
	}

	// infoCmd represents the info command.
	infoCmd = &cobra.Command{
		Use:   "info",
		Short: "Print wallet information",
		Long:  "Print wallet information",
		RunE:  runInfo,
	}
)

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&infoArgs.wallet, "wallet", "w", "", "Name of the wallet to use")
	infoCmd.Flags().StringVarP(&infoArgs.passphraseFile, "passphrase-file", "p", "", "Path of the file containing the passphrase to access the wallet")
	_ = infoCmd.MarkFlagRequired("wallet")
	_ = infoCmd.MarkFlagFilename("passphrase-file")
}

func runInfo(_ *cobra.Command, _ []string) error {
	passphrase, err := getPassphrase(infoArgs.passphraseFile, false)
	if err != nil {
		return err
	}

	store, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	w, err := store.GetWallet(infoArgs.wallet, passphrase)
	if err != nil {
		return fmt.Errorf("couldn't get the wallet %s: %w", infoArgs.wallet, err)
	}

	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		p.Text("Type:").Jump().WarningText(w.Type()).Jump()
		p.Text("Version:").Jump().WarningText(fmt.Sprintf("%d", w.Version())).Jump()
		p.Text("ID:").Jump().WarningText(w.ID()).Jump()
	} else if rootArgs.output == "json" {
		return vgjson.Print(struct {
			Type    string `json:"type"`
			Version uint32 `json:"version"`
			ID      string `json:"id"`
		}{
			Type:    w.Type(),
			Version: w.Version(),
			ID:      w.ID(),
		})
	}

	return nil
}
