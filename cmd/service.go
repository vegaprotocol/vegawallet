package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

func NewCmdService(w io.Writer, rf *RootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Manage the Vega wallet's service",
		Long:  "Manage the Vega wallet's service",
	}

	cmd.AddCommand(NewCmdRunService(w, rf))
	cmd.AddCommand(NewCmdListEndpoints(w, rf))
	return cmd
}
