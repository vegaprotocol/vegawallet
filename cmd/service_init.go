package cmd

import (
	"code.vegaprotocol.io/go-wallet/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	serviceInitArgs struct {
		force       bool
		NoGenRsaKey bool
	}

	serviceInitCmd = &cobra.Command{
		Use:   "init",
		Short: "Generate the configuration",
		Long:  "Generate the configuration for the wallet service",
		RunE:  runServiceInit,
	}
)

func init() {
	serviceCmd.AddCommand(serviceInitCmd)
	serviceInitCmd.Flags().BoolVarP(&serviceInitArgs.force, "force", "f", false, "Erase exiting wallet service configuration at the specified path")
	serviceInitCmd.Flags().BoolVarP(&serviceInitArgs.NoGenRsaKey, "no-genrsakey", "g", false, "Do not generate rsa keys for the jwt tokens by default")
}

func runServiceInit(cmd *cobra.Command, args []string) error {
	log, err := zap.NewProduction()
	if err != nil {
		return err
	}

	store, err := getStore()
	if err != nil {
		return err
	}

	return config.GenerateConfig(log, store, serviceInitArgs.force, !serviceInitArgs.NoGenRsaKey)
}
