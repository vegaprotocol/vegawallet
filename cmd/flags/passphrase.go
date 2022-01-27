package flags

import (
	"errors"
	"fmt"
	"os"
	"strings"

	vgfs "code.vegaprotocol.io/shared/libs/fs"
	vgterm "code.vegaprotocol.io/shared/libs/term"
	"golang.org/x/term"
)

var (
	ErrPassphraseRequiredWithoutTTY = errors.New("passphrase flag is required without TTY")
	ErrPassphraseDoNotMatch         = errors.New("passphrases do not match")
	ErrPassphraseMustBeSpecified    = errors.New("passphrase must be specified")
	ErrMsysPasswordInput            = errors.New("password input is not supported on msys (use --passphrase-file or a standard windows terminal)")
)

type PassphraseGetterWithOps func(bool) (string, error)

// BuildPassphraseGetterWithOps builds a function that returns a passphrase.
// If passphraseFile is set, the returned function is built to read a file. If
// it's not set, the returned function is built to read from user input.
// The one based on the user input takes an argument withConfirmation that
// asks for passphrase confirmation base on its value.
func BuildPassphraseGetterWithOps(passphraseFile string) PassphraseGetterWithOps {
	if len(passphraseFile) != 0 {
		return func(_ bool) (string, error) {
			return ReadPassphraseFile(passphraseFile)
		}
	}

	return ReadPassphraseInputWithOpts
}

func GetPassphrase(passphraseFile string) (string, error) {
	if len(passphraseFile) != 0 {
		return ReadPassphraseFile(passphraseFile)
	}

	return ReadPassphraseInput()
}

func GetConfirmedPassphrase(passphraseFile string) (string, error) {
	if len(passphraseFile) != 0 {
		return ReadPassphraseFile(passphraseFile)
	}

	return ReadConfirmedPassphraseInput()
}

func ReadPassphraseFile(passphraseFilePath string) (string, error) {
	rawPassphrase, err := vgfs.ReadFile(passphraseFilePath)
	if err != nil {
		return "", fmt.Errorf("couldn't read passphrase file: %w", err)
	}

	// user might have added a newline at the end of the line, let's remove it,
	// remembering Windows does things differently
	cleanupPassphrase := strings.Trim(string(rawPassphrase), "\r\n")
	if len(cleanupPassphrase) == 0 {
		return "", ErrPassphraseMustBeSpecified
	}

	return cleanupPassphrase, nil
}

func ReadPassphraseInput() (string, error) {
	return ReadPassphraseInputWithOpts(false)
}

func ReadConfirmedPassphraseInput() (string, error) {
	return ReadPassphraseInputWithOpts(true)
}

func ReadPassphraseInputWithOpts(withConfirmation bool) (string, error) {
	if vgterm.HasNoTTY() {
		return "", ErrPassphraseRequiredWithoutTTY
	}

	passphrase, err := promptForPassphrase("Enter passphrase: ")
	if err != nil {
		return "", fmt.Errorf("couldn't get passphrase: %w", err)
	}
	if len(passphrase) == 0 {
		return "", ErrPassphraseMustBeSpecified
	}

	if withConfirmation {
		confirmation, err := promptForPassphrase("Confirm passphrase: ")
		if err != nil {
			return "", fmt.Errorf("couldn't get passphrase confirmation: %w", err)
		}

		if passphrase != confirmation {
			return "", ErrPassphraseDoNotMatch
		}
	}
	fmt.Println() //nolint:forbidigo

	return passphrase, nil
}

func runningInMsys() bool {
	ms := os.Getenv("MSYSTEM")
	return ms != ""
}

func promptForPassphrase(msg ...string) (string, error) {
	if runningInMsys() {
		return "", ErrMsysPasswordInput
	}
	fmt.Print(msg[0]) //nolint:forbidigo
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("couldn't read password input: %w", err)
	}
	fmt.Println() //nolint:forbidigo

	return string(password), nil
}
