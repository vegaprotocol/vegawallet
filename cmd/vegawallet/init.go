package main

import (
	"fmt"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/wallet"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	force     bool
	genRsaKey bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate the configuration",
	Long:  "Generate the configuration for the wallet service",
	RunE:  runServiceInit,
}

func init() {
	serviceCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&force, "force", "f", false, "Erase exiting wallet service configuration at the specified path")
	initCmd.Flags().BoolVarP(&genRsaKey, "genrsakey", "g", false, "Generate rsa keys for the jwt tokens")
}

func runServiceInit(cmd *cobra.Command, args []string) error {
	if ok, err := fsutil.PathExists(rootPath); !ok {
		if _, ok := err.(*fsutil.PathNotFound); !ok {
			return fmt.Errorf("invalid root directory path: %v", err)
		}
		// create the folder
		if err := fsutil.EnsureDir(rootPath); err != nil {
			return fmt.Errorf("error creating root directory: %v", err)
		}
	}

	log, err := zap.NewProduction()
	if err != nil {
		return err
	}

	return wallet.GenConfig(log, rootPath, force, genRsaKey)
}
