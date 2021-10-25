package cmd

import (
	vgjson "code.vegaprotocol.io/shared/libs/json"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/version"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of the Vega wallet",
	Long:  "Show the version of the Vega wallet",
	RunE:  runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(_ *cobra.Command, _ []string) error {
	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		p.Text("Version:").NextLine().WarningText(version.Version).NextSection()
		p.Text("Git hash:").NextLine().WarningText(version.VersionHash).NextSection()
	} else if rootArgs.output == "json" {
		return printVersionJSON()
	}
	return nil
}

func printVersionJSON() error {
	return vgjson.Print(struct {
		Version string `json:"version"`
		GitHash string `json:"gitHash"`
	}{
		Version: version.Version,
		GitHash: version.VersionHash,
	})
}
