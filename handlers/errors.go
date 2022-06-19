package handlers

import (
	"errors"

	"code.vegaprotocol.io/vegawallet/libs/jsonrpc"
)

const (
	GeneralUserErrorCode = 1
	UserRejectionCode    = 9000
)

var (
	ErrParamsRequired   = errors.New("params is required")
	ErrParamsDoNotMatch = errors.New("params do not match expected ones")
)

func invalidParams(err error) *jsonrpc.ErrorDetails {
	return jsonrpc.NewInvalidParams(err)
}

func applicationError(err error) *jsonrpc.ErrorDetails {
	return jsonrpc.NewServerError(GeneralUserErrorCode, err)
}

func rejectionError(err error) *jsonrpc.ErrorDetails {
	return jsonrpc.NewServerError(UserRejectionCode, err)
}

func internalError(err error) *jsonrpc.ErrorDetails {
	return jsonrpc.NewInternalError(err)
}
