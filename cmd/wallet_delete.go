package cmd

import (
	"errors"
	"fmt"
	"io"

	vgterm "code.vegaprotocol.io/shared/libs/term"
	"code.vegaprotocol.io/vegawallet/cmd/cli"
	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/cmd/printer"
	"code.vegaprotocol.io/vegawallet/wallets"
	"github.com/spf13/cobra"
)

var (
	ErrForceFlagIsRequiredWithoutTTY = errors.New("--force is required without TTY")

	deleteWalletLong = cli.LongDesc(`
		Delete the specified wallet and its keys.

		Be sure to have its recovery phrase, otherwise you won't be able to restore it. If you
		lost it, you should transfer your funds, assets, orders, and anything else attached to 
		this wallet to another wallet.

		The deletion removes the file in which the wallet and its keys are stored, meaning you
		can reuse the wallet name, without causing any conflict.
	`)

	deleteWalletExample = cli.Examples(`
		# Delete the specified wallet
		vegawallet delete --wallet WALLET

		# Delete the specified wallet without asking for confirmation
		vegawallet delete --wallet WALLET --force
	`)
)

type DeleteWalletHandler func(wallet string) error

func NewCmdDeleteWallet(w io.Writer, rf *RootFlags) *cobra.Command {
	h := func(wallet string) error {
		s, err := wallets.InitialiseStore(rf.Home)
		if err != nil {
			return fmt.Errorf("couldn't initialise wallets store: %w", err)
		}

		return s.DeleteWallet(wallet)
	}

	return BuildCmdDeleteWallet(w, h, rf)
}

func BuildCmdDeleteWallet(w io.Writer, handler DeleteWalletHandler, rf *RootFlags) *cobra.Command {
	f := &DeleteWalletFlags{}

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete the specified wallet and its keys",
		Long:    deleteWalletLong,
		Example: deleteWalletExample,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := f.Validate(); err != nil {
				return err
			}

			if !f.Force && vgterm.HasTTY() {
				confirm, err := flags.DoYouConfirm()
				if err != nil {
					return err
				}
				if !confirm {
					return nil
				}
			}

			if err := handler(f.Wallet); err != nil {
				return err
			}

			switch rf.Output {
			case flags.InteractiveOutput:
				PrintDeleteWalletResponse(w, f.Wallet)
			case flags.JSONOutput:
				return nil
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&f.Wallet,
		"wallet", "w",
		"",
		"Wallet to delete",
	)
	cmd.Flags().BoolVarP(&f.Force,
		"force", "f",
		false,
		"Do not ask for confirmation",
	)

	autoCompleteWallet(cmd, rf.Home)

	return cmd
}

type DeleteWalletFlags struct {
	Wallet string
	Force  bool
}

func (f *DeleteWalletFlags) Validate() error {
	if len(f.Wallet) == 0 {
		return flags.FlagMustBeSpecifiedError("wallet")
	}

	if !f.Force && vgterm.HasNoTTY() {
		return ErrForceFlagIsRequiredWithoutTTY
	}

	return nil
}

func PrintDeleteWalletResponse(w io.Writer, walletName string) {
	p := printer.NewInteractivePrinter(w)

	p.CheckMark().SuccessText("Wallet ").SuccessBold(walletName).SuccessText(" deleted").NextLine()
}
