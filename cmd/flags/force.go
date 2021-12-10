package flags

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var ErrExpectYesOrNo = errors.New("invalid answer, expect \"yes\" or \"no\"")

func DoYouConfirm() (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Are you sure? (y/n) ") //nolint:forbidigo
	answer, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("couldn't read password input: %w", err)
	}
	fmt.Println() //nolint:forbidigo

	answer = strings.ToLower(strings.Trim(answer, " \n\t"))

	switch answer {
	case "yes", "y":
		return true, nil
	case "no", "n":
		return false, nil
	default:
		return false, ErrExpectYesOrNo
	}
}
