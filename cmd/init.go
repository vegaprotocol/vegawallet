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
		NoGenRsaKey bool
	}

	initCmd = &cobra.Command{
		Use:   "init",
		Short: "Generate the configuration",
		Long:  "Generate the configuration for the wallet service",
		RunE:  runServiceInit,
	}
)

func init() {
	serviceCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&initArgs.force, "force", "f", false, "Erase exiting wallet service configuration at the specified path")
	initCmd.Flags().BoolVarP(&initArgs.NoGenRsaKey, "no-genrsakey", "g", false, "Do not generate rsa keys for the jwt tokens by default")
}

func runServiceInit(cmd *cobra.Command, args []string) error {
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

	return config.GenerateConfig(log, store, initArgs.force, !initArgs.NoGenRsaKey)
}
