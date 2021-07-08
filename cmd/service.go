package cmd

import (
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage the service",
	Long:  "Manage the service",
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
