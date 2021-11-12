package flags

import (
	"errors"
)

const (
	InteractiveOutput = "interactive"
	JSONOutput        = "json"
)

var (
	ErrUnsupportedOutput = errors.New("unsupported output")

	AvailableOutputs = []string{
		InteractiveOutput,
		JSONOutput,
	}
)

func ValidateOutput(output string) error {
	if len(output) == 0 {
		return FlagMustBeSpecifiedError("output")
	}

	for _, o := range AvailableOutputs {
		if output == o {
			return nil
		}
	}

	// The output flag has special treatment because error reporting depends on
	// it, and we need to differentiate output errors from the rest to select
	// the right way to print the data.
	// As a result, we return a specific error, instead of a generic
	// UnsupportedFlagValueError.
	return ErrUnsupportedOutput
}
