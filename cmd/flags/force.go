package flags

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func AreYouSure() bool {
	return YesOrNo("Are you sure?")
}

func DoYouApproveTx() bool {
	return YesOrNo("Do you approve this transaction?")
}

func YesOrNo(question string) bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(question + " (y/n) ") //nolint:forbidigo

		answer, err := reader.ReadString('\n')
		if err != nil {
			panic(fmt.Errorf("couldn't read input: %w", err))
		}

		answer = strings.ToLower(strings.Trim(answer, " \r\n\t"))

		switch answer {
		case "yes", "y":
			return true
		case "no", "n":
			return false
		default:
			fmt.Printf("invalid answer \"%s\", expect \"yes\" or \"no\"\n", answer) //nolint:forbidigo
		}
	}
}
