package service

import (
	"errors"

	"code.vegaprotocol.io/protos/commands"
)

var (
	ErrInvalidToken                 = errors.New("invalid token")
	ErrInvalidClaims                = errors.New("invalid claims")
	ErrInvalidOrMissingToken        = newErrorResponse("invalid or missing token")
	ErrCouldNotReadRequest          = errors.New("couldn't read request")
	ErrCouldNotGetBlockHeight       = errors.New("couldn't get last block height")
	ErrShouldBeBase64Encoded        = errors.New("should be base64 encoded")
	ErrRSAKeysAlreadyExists         = errors.New("RSA keys already exist")
	ErrCouldNotGetPoW               = errors.New("could not get proof of work")
	ErrRejectedSignRequest          = errors.New("user rejected sign request")
	ErrInvalidSignRequestConfirm    = errors.New("invalid sign request confirmation")
	ErrUnexpectedSignRequestConfirm = errors.New("unexpected sign request confirmation")
)

type ErrorsResponse struct {
	Errors commands.Errors `json:"errors"`
}

type ErrorResponse struct { //nolint:errname
	ErrorStr string   `json:"error"`
	Details  []string `json:"details,omitempty"`
}

func (e ErrorResponse) Error() string {
	return e.ErrorStr
}

func newErrorResponse(e string) ErrorResponse {
	return ErrorResponse{
		ErrorStr: e,
	}
}

func newErrorWithDetails(e string, details []string) ErrorResponse {
	return ErrorResponse{
		ErrorStr: e,
		Details:  details,
	}
}
