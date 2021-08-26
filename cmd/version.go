package cmd

import (
	"code.vegaprotocol.io/go-wallet/cmd/printer"
	vgjson "code.vegaprotocol.io/go-wallet/libs/json"
	"code.vegaprotocol.io/go-wallet/version"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of the vega wallet",
	Long:  "Show the version of the vega wallet",
	RunE:  runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(_ *cobra.Command, _ []string) error {
	if rootArgs.output == "human" {
		p := printer.NewHumanPrinter()
		p.Text("Version:").Jump().WarningText(version.Version).NJump(2)
		p.Text("Git hash:").Jump().WarningText(version.VersionHash).NJump(2)
	} else if rootArgs.output == "json" {
		return printVersionJson()
	}
	return nil
}

func printVersionJson() error {
	return vgjson.Print(struct {
		Version string
		GitHash string
	}{
		Version: version.Version,
		GitHash: version.VersionHash,
	})
}
