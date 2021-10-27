package network

import "fmt"

type NetworkAlreadyExistsError struct {
	Name string
}

func NewNetworkAlreadyExistsError(n string) NetworkAlreadyExistsError {
	return NetworkAlreadyExistsError{
		Name: n,
	}
}

func (e NetworkAlreadyExistsError) Error() string {
	return fmt.Sprintf("network \"%s\" already exists", e.Name)
}

type NetworkDoesNotExistError struct {
	Name string
}

func NewNetworkDoesNotExistError(n string) NetworkDoesNotExistError {
	return NetworkDoesNotExistError{
		Name: n,
	}
}

func (e NetworkDoesNotExistError) Error() string {
	return fmt.Sprintf("network \"%s\" doesn't exist", e.Name)
}
