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

	cmd.AddCommand(NewCmdCommandSend(w, rf))
	cmd.AddCommand(NewCmdCommandSign(w, rf))
	return cmd
}
