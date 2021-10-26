package cmd

import (
	"errors"
	"fmt"
)

var (
	ErrShouldSetNodeAddressOrNetworkFlag           = errors.New("should set node address or network flag")
	ErrCanNotHaveBothNodeAddressAndNetworkFlagsSet = errors.New("can't have both node address and network flag set")
	ErrInvalidMetadataFormat                       = errors.New("invalid metadata format")
	ErrUseJSONOutputInScript                       = errors.New("output \"human\" is not script-friendly, use \"json\" instead")
	ErrMetaOrClearIsRequired                       = errors.New("`--meta` is required or use `--clear` flag")
	ErrCanNotHaveBothMetadataAndClearFlagsSet      = errors.New("can't have `--meta` and `--clear` flags at the same time")
)

type UnsupportedCommandOutputError struct {
	UnsupportedOutput string
}

func NewUnsupportedCommandOutputError(o string) UnsupportedCommandOutputError {
	return UnsupportedCommandOutputError{
		UnsupportedOutput: o,
	}
}

func (e UnsupportedCommandOutputError) Error() string {
	return fmt.Sprintf("output \"%s\" is not supported for this command", e.UnsupportedOutput)
}

type UnsupportedOutputError struct {
	UnsupportedOutput string
}

func NewUnsupportedOutputError(o string) UnsupportedOutputError {
	return UnsupportedOutputError{
		UnsupportedOutput: o,
	}
}

func (e UnsupportedOutputError) Error() string {
	return fmt.Sprintf("unsupported output \"%s\"", e.UnsupportedOutput)
}
