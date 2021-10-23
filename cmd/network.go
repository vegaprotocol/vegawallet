package cmd

import (
	"github.com/spf13/cobra"
)

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Manage networks",
	Long:  "Manage networks",
}

func init() {
	rootCmd.AddCommand(networkCmd)
}
