package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

func NewCmdCommand(w io.Writer, rf *RootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "command",
		Short: "Provides utilities for interacting with commands",
		Long:  "Provides utilities for interacting with commands",
	}

	// create subcommands
	cmd.AddCommand(NewCmdCommandSend(w, rf))
	return cmd
}
