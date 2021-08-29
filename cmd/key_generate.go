package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
	vgjson "code.vegaprotocol.io/go-wallet/libs/json"
	"code.vegaprotocol.io/go-wallet/wallet"
	"github.com/spf13/cobra"
)

var (
	keyGenerateArgs struct {
		name           string
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
	keyGenerateCmd.Flags().StringVarP(&keyGenerateArgs.passphraseFile, "passphrase-file", "p", "", "Path of the file containing the passphrase to access the wallet")
	keyGenerateCmd.Flags().StringVarP(&keyGenerateArgs.metas, "meta", "m", "", `A list of metadata e.g: "primary:true;asset:BTC"`)
}

func runKeyGenerate(_ *cobra.Command, _ []string) error {
	p := printer.NewHumanPrinter()
	store, err := newWalletsStore(rootArgs.rootPath)
	if err != nil {
		return err
	}

	handler := wallet.NewHandler(store)
	if err != nil {
		return err
	}

	if len(keyGenerateArgs.name) == 0 {
		return errors.New("wallet name is required")
	}

	walletExists := handler.WalletExists(keyGenerateArgs.name)

	passphrase, err := getPassphrase(keyGenerateArgs.passphraseFile, !walletExists)
	if err != nil {
		return err
	}

	metas, err := parseMeta(keyGenerateArgs.metas)
	if err != nil {
		return err
	}

	var mnemonic string
	if !walletExists {
		if rootArgs.output == "human" {
			p.BangMark().Text("Wallet ").Bold(keyGenerateArgs.name).Text(" does not exist yet").Jump()
		}

		mnemonic, err = handler.CreateWallet(keyGenerateArgs.name, passphrase)
		if err != nil {
			return fmt.Errorf("couldn't create wallet: %w", err)
		}

		if rootArgs.output == "human" {
			p.CheckMark().Text("Wallet ").Bold(keyGenerateArgs.name).Text(" has been created at: ").SuccessText(store.GetWalletPath(keyGenerateArgs.name)).Jump()
		}
	}

	keyPair, err := handler.GenerateKeyPair(keyGenerateArgs.name, passphrase, metas)
	if err != nil {
		return fmt.Errorf("could not generate a key pair: %w", err)
	}

	if rootArgs.output == "human" {
		printHuman(p, mnemonic, keyPair)
	} else if rootArgs.output == "json" {
		return printJSON(mnemonic, keyPair)
	} else {
		return fmt.Errorf("output \"%s\" is not supported for this command", rootArgs.output)
	}
	return nil
}

func printHuman(p *printer.HumanPrinter, mnemonic string, keyPair wallet.KeyPair) {
	p.CheckMark().Text("Key pair has been generated for wallet ").Bold(keyGenerateArgs.name).Jump()
	p.CheckMark().SuccessText("Generating a key pair succeeded").NJump(2)
	if len(mnemonic) != 0 {
		p.Text("Wallet mnemonic:").Jump().WarningText(mnemonic).Jump()
	}
	printKeyPair(p, keyPair)
	p.Jump()

	p.RedArrow().DangerText("Important").Jump()
	if len(mnemonic) != 0 {
		p.DangerText("1. ").Text("Save the mnemonic somewhere safe and secure, now. It won’t be displayed ever again!").Jump()
		p.DangerText("2. ").Text("Do not share your private key.").NJump(2)
	} else {
		p.Text("Do not share your private key.").NJump(2)
	}

	p.BlueArrow().InfoText("Run the service").Jump()
	p.Text("Once you have a key pair generated, you can run the service with the following command:").NJump(2)
	p.Code(fmt.Sprintf("%s service run", os.Args[0])).NJump(2)
	p.Text("If you want to open up a local version of Vega Console alongside the service, use the following command:").NJump(2)
	p.Code(fmt.Sprintf("%s service run --console-proxy", os.Args[0])).NJump(2)
	p.Text("To terminate the process, hit ").Bold("ctrl+c").NJump(2)
	p.Text("For more information, use ").Bold("--help").Text(" flag.").Jump()
}

func printJSON(mnemonic string, keyPair wallet.KeyPair) error {
	result := struct {
		WalletMnemonic   string `json:",omitempty"`
		PrivateKey       string
		PublicKey        string
		AlgorithmName    string
		AlgorithmVersion uint32
		Meta             []wallet.Meta
	}{
		WalletMnemonic:   mnemonic,
		PrivateKey:       keyPair.PrivateKey(),
		PublicKey:        keyPair.PublicKey(),
		AlgorithmName:    keyPair.AlgorithmName(),
		AlgorithmVersion: keyPair.AlgorithmVersion(),
		Meta:             keyPair.Meta(),
	}
	return vgjson.Print(result)
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
			return nil, fmt.Errorf("invalid metadata format")
		}
		metas = append(metas, wallet.Meta{Key: rawMeta[0], Value: rawMeta[1]})
	}

	return metas, nil
}
