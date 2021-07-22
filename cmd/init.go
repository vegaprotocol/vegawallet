package cmd

import (
	"code.vegaprotocol.io/go-wallet/logger"
	"code.vegaprotocol.io/go-wallet/service"
	"code.vegaprotocol.io/go-wallet/service/store/v1"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
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

func runInit(cmd *cobra.Command, args []string) error {
	log, err := logger.New(zapcore.InfoLevel)
	if err != nil {
		return err
	}
	defer log.Sync()

	wStore, err := newWalletsStore(rootArgs.rootPath)
	if err != nil {
		return err
	}

	if err := wStore.Initialise(); err != nil {
		return err
	}

	svcStore, err := v1.NewStore(rootArgs.rootPath)
	if err != nil {
		return err
	}

	if err := svcStore.Initialise(); err != nil {
		return err
	}

	return service.GenerateConfig(log, svcStore, initArgs.force)
}
