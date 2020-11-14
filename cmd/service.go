package cmd

import (
	"github.com/spf13/cobra"
)

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "The wallet service",
	Long:  "Run or initialize the wallet service",
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
