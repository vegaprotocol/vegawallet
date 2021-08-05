package cmd

import (
	"fmt"

	"code.vegaprotocol.io/go-wallet/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of the vega wallet",
	Long:  "Show the version of the vega wallet",
	RunE:  runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) error {
	fmt.Printf("vegawallet version %v (%v)\n", version.Version, version.VersionHash)
	return nil
}
