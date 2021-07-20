package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	serviceInitArgs struct {
		force bool
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
	fmt.Println("\n\nDEPRECATION:\nThe command `service init` is deprecated. Use `init` instead.")

	initArgs.force = serviceInitArgs.force
	return runInit(cmd, args)
}
