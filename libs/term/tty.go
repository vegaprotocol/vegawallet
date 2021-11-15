package term

import (
	"os"

	"github.com/mattn/go-isatty"
)

func HasTTY() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}

func HasNoTTY() bool {
	return !HasTTY()
}
