package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

func NewCmdMessage(w io.Writer, rf *RootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "message",
		Short: "Sign and verify messages",
		Long:  "Sign and verify messages",
	}

	cmd.AddCommand(NewCmdSignMessage(w, rf))
	cmd.AddCommand(NewCmdVerifyMessage(w, rf))
	return cmd
}
