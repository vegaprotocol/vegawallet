package cmd

import (
	"fmt"

	"code.vegaprotocol.io/go-wallet/config"
	storev1 "code.vegaprotocol.io/go-wallet/store/v1"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	serviceInitArgs struct {
		force       bool
	}

	serviceInitCmd = &cobra.Command{
		Use:   "init",
		Short: "Generate the configuration (deprecated)",
		Long:  "DEPRECATED! Use init instead. Generate the configuration for the wallet service.",
		RunE:  runServiceInit,
	}
)

func init() {
	serviceCmd.AddCommand(serviceInitCmd)
	serviceInitCmd.Flags().BoolVarP(&serviceInitArgs.force, "force", "f", false, "Erase exiting wallet service configuration at the specified path")
}

func runServiceInit(cmd *cobra.Command, args []string) error {
	log, err := zap.NewProduction()
	if err != nil {
		return err
	}

	fmt.Println("\n\nDEPRECATION:\nThe command `service init` is deprecated. Use `init` instead.")

	store, err := storev1.NewStore(rootArgs.rootPath)
	if err != nil {
		return err
	}

	if err := store.Initialise(); err != nil {
		return err
	}

	return config.GenerateConfig(log, store, initArgs.force)
}
