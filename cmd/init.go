package cmd

import (
	"code.vegaprotocol.io/go-wallet/config"
	storev1 "code.vegaprotocol.io/go-wallet/store/v1"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	initArgs struct {
		force       bool
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
	log, err := zap.NewProduction()
	if err != nil {
		return err
	}

	store, err := storev1.NewStore(rootArgs.rootPath)
	if err != nil {
		return err
	}

	if err := store.Initialise(); err != nil {
		return err
	}

	return config.GenerateConfig(log, store, initArgs.force)
}
