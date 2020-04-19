package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	force     bool
	genRsaKey bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
	},
}

func init() {
	serviceCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&force, "force", "f", false, "Erase exiting wallet service configuration at the specified path")
	initCmd.Flags().BoolVarP(&genRsaKey, "genrsakey", "g", false, "Generate rsa keys for the jwt tokens")
}
