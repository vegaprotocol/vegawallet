package handlers

import (
	"fmt"

	"code.vegaprotocol.io/vegawallet/libs/jsonrpc"
	"code.vegaprotocol.io/vegawallet/permissions"
)

var ErrUserRejectedPermissionsRequest =

type RequestPermissions struct {
	store PermissionsStore
}

func NewRequestPermissions(store PermissionsStore) *RequestPermissions {
	return &RequestPermissions{
		store: store,
	}
}

func (s *RequestPermissions) Handle(rawParams jsonrpc.Params) (jsonrpc.Result, *jsonrpc.ErrorDetails) {
	params, err := validateRequestPermissionsParams(rawParams)
	if err != nil {
		return nil, invalidParams(err)
	}

	hostname, err := toHostname(params.Referer)
	if err != nil {
		return nil, invalidParams(err)
	}

	// TODO Request a review

	approved := true
	if !approved {
		return nil, applicationError(ErrUserRejectedPermissionsRequest)
	}

	detailedPerms := permissions.Permissions{}

	if access, ok := params.RequestedPermissions["public_keys"]; ok {
		detailedPerms.PublicKeys = &permissions.PublicKeysPermissions{
			Access:         permissions.AccessMode(access),
			RestrictedKeys: nil,
		}
	}

	if err := s.store.SavePermissions(hostname, detailedPerms); err != nil {
		return nil, internalError(fmt.Errorf("couldn't save permissions: %w", err))
	}

	return RequestPermissionsResult{
		Permissions: detailedPerms.Summary(),
	}, nil
}

func validateRequestPermissionsParams(rawParams jsonrpc.Params) (RequestPermissionsParams, error) {
	if rawParams == nil {
		return RequestPermissionsParams{}, ErrParamsRequired
	}

	params, ok := rawParams.(RequestPermissionsParams)
	if !ok {
		return RequestPermissionsParams{}, ErrParamsDoNotMatch
	}

	if params.Referer == "" {
		return RequestPermissionsParams{}, ErrRefererIsRequired
	}

	// TODO Validate requested permissions

	return params, nil
}

type RequestPermissionsParams struct {
	Referer              string            `json:"referer"`
	RequestedPermissions map[string]string `json:"requestedPermissions"`
}

type RequestPermissionsResult struct {
	Permissions map[string]string `json:"permissions"`
}
