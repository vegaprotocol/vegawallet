package flags

const (
	InteractiveOutput = "human"
	JSONOutput        = "json"
)

var AvailableOutputs = []string{
	InteractiveOutput,
	JSONOutput,
}

func ValidateOutput(output string) error {
	if len(output) == 0 {
		return FlagMustBeSpecifiedError("output")
	}

	for _, o := range AvailableOutputs {
		if output == o {
			return nil
		}
	}

	availableOutputs := make([]interface{}, len(AvailableOutputs))
	for i := range AvailableOutputs {
		availableOutputs[i] = AvailableOutputs[i]
	}

	return UnsupportedFlagValueError("output", output, availableOutputs)
}
