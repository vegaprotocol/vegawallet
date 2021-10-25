package v1

import "fmt"

type DifferentNetworkNamesError struct {
	fileName   string
	configName string
}

func NewDifferentNetworkNamesError(f, c string) DifferentNetworkNamesError {
	return DifferentNetworkNamesError{
		fileName:   f,
		configName: c,
	}
}

func (e DifferentNetworkNamesError) Error() string {
	return fmt.Sprintf("file name (%s) and configuration name (%s) are different", e.fileName, e.configName)
}
