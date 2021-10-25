package cmd

import (
	"fmt"
	"os"

	vgjson "code.vegaprotocol.io/shared/libs/json"
	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/network"
	netstore "code.vegaprotocol.io/vegawallet/network/store/v1"
	"code.vegaprotocol.io/vegawallet/service"
	svcstore "code.vegaprotocol.io/vegawallet/service/store/v1"
	"code.vegaprotocol.io/vegawallet/wallets"

	"github.com/spf13/cobra"
)

var (
	initArgs struct {
		force bool
	}

	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialise the wallet",
		Long:  "Creates the folders, the configuration file and RSA keys needed by the service to operate.",
		RunE:  runInit,
	}
)

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&initArgs.force, "force", "f", false, "Erase exiting wallet service configuration at the specified path")
}

func runInit(_ *cobra.Command, _ []string) error {
	_, err := wallets.InitialiseStore(rootArgs.home)
	if err != nil {
		return fmt.Errorf("couldn't initialise wallets store: %w", err)
	}

	svcStore, err := svcstore.InitialiseStore(paths.New(rootArgs.home))
	if err != nil {
		return fmt.Errorf("couldn't initialise service store: %w", err)
	}

	if err = service.InitialiseService(svcStore, initArgs.force); err != nil {
		return fmt.Errorf("couldn't initialise the service: %w", err)
	}

	netStore, err := netstore.InitialiseStore(paths.New(rootArgs.home))
	if err != nil {
		return fmt.Errorf("couldn't initialise service store: %w", err)
	}

	if err = network.InitialiseNetworks(netStore, initArgs.force); err != nil {
		return fmt.Errorf("couldn't initialise the networks: %w", err)
	}

	if rootArgs.output == "human" {
		printInitHuman(svcStore, netStore)
	} else if rootArgs.output == "json" {
		return printInitJson(svcStore, netStore)
	} else {
		return NewUnsupportedCommandOutputError(rootArgs.output)
	}

	return nil
}

func printInitHuman(svcStore *svcstore.Store, netStore *netstore.Store) {
	p := printer.NewHumanPrinter()
	p.CheckMark().Text("Networks configurations created at: ").SuccessText(netStore.GetNetworksPath()).Jump()
	pubRSAKeysPath, privRSAKeysPath := svcStore.GetRSAKeysPath()
	p.CheckMark().Text("Service public RSA keys created at: ").SuccessText(pubRSAKeysPath).Jump()
	p.CheckMark().Text("Service private RSA keys created at: ").SuccessText(privRSAKeysPath).Jump()
	p.CheckMark().SuccessText("Initialisation succeeded").NJump(2)

	p.BlueArrow().InfoText("Create a wallet").Jump()
	p.Text("To create a wallet, generate your first key pair using the following command:").NJump(2)
	p.Code(fmt.Sprintf("%s key generate --wallet \"YOUR_USERNAME\"", os.Args[0])).NJump(2)
	p.Text("The ").Bold("--wallet").Text(" flag sets the wallet of your wallet and will be used to login to Vega Console.").NJump(2)
	p.Text("For more information, use ").Bold("--help").Text(" flag.").Jump()
}

type initJson struct {
	RSAKeys      initRsaKeysJson `json:"rsaKeys"`
	NetworksHome string          `json:"networksHome"`
}

type initRsaKeysJson struct {
	PublicKeyFilePath  string `json:"publicKeyFilePath"`
	PrivateKeyFilePath string `json:"privateKeyFilePath"`
}

func printInitJson(svcStore *svcstore.Store, netStore *netstore.Store) error {
	pubRSAKeysPath, privRSAKeysPath := svcStore.GetRSAKeysPath()
	result := initJson{
		RSAKeys: initRsaKeysJson{
			PublicKeyFilePath:  pubRSAKeysPath,
			PrivateKeyFilePath: privRSAKeysPath,
		},
		NetworksHome: netStore.GetNetworksPath(),
	}
	return vgjson.Print(result)
}
