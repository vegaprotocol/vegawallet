package cmd

import (
	"fmt"
	"os"

	"code.vegaprotocol.io/go-wallet/cmd/printer"
	"code.vegaprotocol.io/go-wallet/service"
	"code.vegaprotocol.io/go-wallet/service/store/v1"
	"code.vegaprotocol.io/go-wallet/wallets"
	vgjson "code.vegaprotocol.io/shared/libs/json"
	"code.vegaprotocol.io/shared/paths"
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

	svcStore, err := v1.InitialiseStore(paths.NewPaths(rootArgs.home))
	if err != nil {
		return err
	}

	if err = service.GenerateConfig(svcStore, initArgs.force); err != nil {
		return err
	}

	if rootArgs.output == "human" {
		printInitHuman(svcStore)
	} else if rootArgs.output == "json" {
		return printInitJson(svcStore)
	} else {
		return fmt.Errorf("output \"%s\" is not supported for this command", rootArgs.output)
	}

	return nil
}

func printInitHuman(svcStore *v1.Store) {
	p := printer.NewHumanPrinter()
	p.CheckMark().Text("Service configuration created at: ").SuccessText(svcStore.GetConfigPath()).Jump()
	rsaKeysPath := svcStore.GetRSAKeysPath()
	p.CheckMark().Text("Service public RSA keys created at: ").SuccessText(rsaKeysPath["public"]).Jump()
	p.CheckMark().Text("Service private RSA keys created at: ").SuccessText(rsaKeysPath["private"]).Jump()
	p.CheckMark().SuccessText("Initialisation succeeded").NJump(2)

	p.BlueArrow().InfoText("Create a wallet").Jump()
	p.Text("To create a wallet, generate your first key pair using the following command:").NJump(2)
	p.Code(fmt.Sprintf("%s key generate --name \"YOUR_USERNAME\"", os.Args[0])).NJump(2)
	p.Text("The ").Bold("--name").Text(" flag sets the name of your wallet and will be used to login to Vega Console.").NJump(2)
	p.Text("For more information, use ").Bold("--help").Text(" flag.").Jump()
}

func printInitJson(svcStore *v1.Store) error {
	rsaKeysPath := svcStore.GetRSAKeysPath()
	result := struct {
		PublicRSAKeyFile  string
		PrivateRSAKeyFile string
		ServiceConfigFile string
	}{
		PublicRSAKeyFile:  rsaKeysPath["public"],
		PrivateRSAKeyFile: rsaKeysPath["private"],
		ServiceConfigFile: svcStore.GetConfigPath(),
	}
	return vgjson.Print(result)
}
