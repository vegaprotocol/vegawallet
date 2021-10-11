package service

import (
	"errors"

	"code.vegaprotocol.io/protos/commands"
)

var (
	ErrInvalidOrMissingToken  = newErrorResponse("invalid or missing token")
	ErrCouldNotReadRequest    = errors.New("could not read request")
	ErrCouldNotGetBlockHeight = errors.New("could not get last block height")
	ErrShouldBeBase64Encoded  = errors.New("should be base64 encoded")
	ErrCouldNotMarshalTxResponse = errors.New("could not marshal transaction response")
)

type ErrorsResponse struct {
	Errors commands.Errors `json:"errors"`
}

type ErrorResponse struct {
	ErrorStr string   `json:"error"`
	Details  []string `json:"details"`
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
