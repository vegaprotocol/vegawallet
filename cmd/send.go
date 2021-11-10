package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

func NewCmdSend(w io.Writer, rf *RootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send data to the Vega network",
		Long:  "Send data to the Vega network",
	}

	// create subcommands
	cmd.AddCommand(NewCmdSendCommand(w, rf))
	return cmd
}
