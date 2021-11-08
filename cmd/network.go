package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

func NewCmdNetwork(w io.Writer, rf *RootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Manage networks",
		Long:  "Manage networks",
	}

	// create subcommands
	cmd.AddCommand(NewCmdListNetworks(w, rf))
	cmd.AddCommand(NewCmdImportNetwork(w, rf))
	return cmd
}
