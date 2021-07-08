package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	defaultVersionHash = "unknown"
	defaultVersion     = "unknown"
)

var (
	// VersionHash specifies the git commit used to build the application. See VERSION_HASH in Makefile for details.
	VersionHash = defaultVersionHash

	// Version specifies the version used to build the application. See VERSION in Makefile for details.
	Version = defaultVersion
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
	fmt.Printf("vegawallet version %v (%v)\n", Version, VersionHash)
	return nil
}
