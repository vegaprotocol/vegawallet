package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

var completionLong = `To load completions:

Bash:  To load completions for each session, execute once:
-----  Linux:
       $ vegawallet completion bash > /etc/bash_completion.d/vegawallet
       MacOS:
       $ vegawallet completion bash > /usr/local/etc/bash_completion.d/vegawallet


Zsh:   If shell completion is not already enabled in your environment you will need
----   to enable it.  You can execute the following once:
       $ echo "autoload -U compinit; compinit" >> ~/.zshrc

       To load completions for each session, execute once:
       $ vegawallet completion zsh > "${fpath[1]}/_vegawallet"

       You will need to start a new shell for this setup to take effect.


Fish:  To load completions for each session, execute once:
-----  $ vegawallet completion fish > ~/.config/fish/completions/vegawallet.fish
`

func NewCmdCompletion(w io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:                   "completion [bash|zsh|fish|powershell]",
		Short:                 "Generate completion script",
		Long:                  completionLong,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				_ = cmd.Root().GenBashCompletion(w)
			case "zsh":
				_ = cmd.Root().GenZshCompletion(w)
			case "fish":
				_ = cmd.Root().GenFishCompletion(w, true)
			case "powershell":
				_ = cmd.Root().GenPowerShellCompletion(w)
			}
		},
	}
}
