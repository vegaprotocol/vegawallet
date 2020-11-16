package cmd

import (
	"fmt"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/wallet"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	initArgs struct {
		force     bool
		genRsaKey bool
	}

	// initCmd represents the init command
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
	initCmd.Flags().BoolVarP(&initArgs.genRsaKey, "no-genrsakey", "g", false, "Do not generate rsa keys for the jwt tokens by default")
}

func runServiceInit(cmd *cobra.Command, args []string) error {
	if ok, err := fsutil.PathExists(rootArgs.rootPath); !ok {
		if _, ok := err.(*fsutil.PathNotFound); !ok {
			return fmt.Errorf("invalid root directory path: %v", err)
		}
		// create the folder
		if err := fsutil.EnsureDir(rootArgs.rootPath); err != nil {
			return fmt.Errorf("error creating root directory: %v", err)
		}
	}

	log, err := zap.NewProduction()
	if err != nil {
		return err
	}

	return wallet.GenConfig(log, rootArgs.rootPath, initArgs.force, !initArgs.genRsaKey)
}
