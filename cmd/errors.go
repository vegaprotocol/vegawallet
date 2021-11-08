package cmd

type InvalidOutputError struct {
	err error
}

func NewInvalidOutputError(err error) InvalidOutputError {
	return InvalidOutputError{
		err: err,
	}
}

func (e InvalidOutputError) Error() string {
	return e.err.Error()
}
