package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
	"code.vegaprotocol.io/go-wallet/wallet"
	"code.vegaprotocol.io/go-wallet/wallets"
	vgjson "code.vegaprotocol.io/shared/libs/json"
	"github.com/spf13/cobra"
)

var (
	keyGenerateArgs struct {
		wallet         string
		passphraseFile string
		metas          string
	}

	keyGenerateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate a new key pair for a wallet",
		Long:  "Generate a new key pair for a wallet, this will implicitly generate a new wallet if none exist for the given wallet",
		RunE:  runKeyGenerate,
	}
)

func init() {
	keyCmd.AddCommand(keyGenerateCmd)
	keyGenerateCmd.Flags().StringVarP(&keyGenerateArgs.wallet, "wallet", "w", "", "Name of the wallet to use")
	keyGenerateCmd.Flags().StringVarP(&keyGenerateArgs.passphraseFile, "passphrase-file", "p", "", "Path of the file containing the passphrase to access the wallet")
	keyGenerateCmd.Flags().StringVarP(&keyGenerateArgs.metas, "meta", "m", "", `A list of metadata e.g: "primary:true;asset:BTC"`)
}

func runKeyGenerate(_ *cobra.Command, _ []string) error {
	p := printer.NewHumanPrinter()

	store, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	handler := wallets.NewHandler(store)

	if len(keyGenerateArgs.wallet) == 0 {
		return errors.New("wallet is required")
	}

	walletExists := handler.WalletExists(keyGenerateArgs.wallet)

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
			p.BangMark().Text("Wallet ").Bold(keyGenerateArgs.wallet).Text(" does not exist yet").Jump()
		}

		mnemonic, err = handler.CreateWallet(keyGenerateArgs.wallet, passphrase)
		if err != nil {
			return fmt.Errorf("couldn't create wallet: %w", err)
		}

		if rootArgs.output == "human" {
			p.CheckMark().Text("Wallet ").Bold(keyGenerateArgs.wallet).Text(" has been created at: ").SuccessText(store.GetWalletPath(keyGenerateArgs.wallet)).Jump()
		}
	}

	keyPair, err := handler.GenerateKeyPair(keyGenerateArgs.wallet, passphrase, metas)
	if err != nil {
		return fmt.Errorf("could not generate a key pair: %w", err)
	}

	if rootArgs.output == "human" {
		printHuman(p, mnemonic, keyPair, store.GetWalletPath(keyGenerateArgs.wallet))
	} else if rootArgs.output == "json" {
		return printKeyGenerateJSON(mnemonic, keyPair, store.GetWalletPath(keyGenerateArgs.wallet))
	} else {
		return fmt.Errorf("output \"%s\" is not supported for this command", rootArgs.output)
	}
	return nil
}

func printHuman(p *printer.HumanPrinter, mnemonic string, keyPair wallet.KeyPair, walletPath string) {
	p.CheckMark().Text("Key pair has been generated for wallet ").Bold(keyGenerateArgs.wallet).Text(" at: ").SuccessText(walletPath).Jump()
	p.CheckMark().SuccessText("Generating a key pair succeeded").NJump(2)
	if len(mnemonic) != 0 {
		p.Text("Wallet mnemonic:").Jump().WarningText(mnemonic).Jump()
	}
	p.Text("Public key:").Jump().WarningText(keyPair.PublicKey()).Jump()
	p.Text("Metadata:").Jump()
	printMeta(p, keyPair.Meta())
	p.Jump()

	p.RedArrow().DangerText("Important").Jump()
	if len(mnemonic) != 0 {
		p.DangerText("1. ").Text("Write down the mnemonic and store it somewhere safe and secure, now, as it will ").Underline("not").Text(" be displayed ever again!").Jump()
		p.DangerText("2. ").Text("Do not share the mnemonic nor the private key.").NJump(2)
	} else {
		p.Text("Do not share the mnemonic nor the private key.").NJump(2)
	}

	p.BlueArrow().InfoText("Run the service").Jump()
	p.Text("Once you have a key pair generated, you can run the service with the following command:").NJump(2)
	p.Code(fmt.Sprintf("%s service run", os.Args[0])).NJump(2)
	p.Text("If you want to open up a local version of Vega Console alongside the service, use the following command:").NJump(2)
	p.Code(fmt.Sprintf("%s service run --console-proxy", os.Args[0])).NJump(2)
	p.Text("To terminate the process, hit ").Bold("ctrl+c").NJump(2)
	p.Text("For more information, use ").Bold("--help").Text(" flag.").Jump()
}

type keyGenerateJson struct {
	Wallet keyGenerateWalletJson `json:"wallet"`
	Key    keyGenerateKeyJson    `json:"key"`
}

type keyGenerateWalletJson struct {
	FilePath string `json:"filePath"`
	Mnemonic string `json:"mnemonic,omitempty"`
}

type keyGenerateKeyJson struct {
	KeyPair   keyGenerateKeyPairJson   `json:"keyPair"`
	Algorithm keyGenerateAlgorithmJson `json:"algorithm"`
	Meta      []wallet.Meta            `json:"meta"`
}

type keyGenerateKeyPairJson struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

type keyGenerateAlgorithmJson struct {
	Name    string `json:"wallet"`
	Version uint32 `json:"version"`
}

func printKeyGenerateJSON(mnemonic string, keyPair wallet.KeyPair, walletPath string) error {
	result := keyGenerateJson{
		Wallet: keyGenerateWalletJson{
			FilePath: walletPath,
			Mnemonic: mnemonic,
		},
		Key: keyGenerateKeyJson{
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
