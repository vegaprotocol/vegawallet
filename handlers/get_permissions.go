package handlers

import (
	"errors"
	"fmt"
	"net/url"

	"code.vegaprotocol.io/vegawallet/libs/jsonrpc"
	"code.vegaprotocol.io/vegawallet/permissions"
)

var (
	ErrRefererIsRequired    = errors.New("referer is required")
	ErrRefererIsNotValidURL = errors.New("referer is not a valid URL")
)

//go:generate go run github.com/golang/mock/mockgen -destination mocks/permissions_store_mock.go -package mocks code.vegaprotocol.io/vegawallet/handlers PermissionsStore
type PermissionsStore interface {
	PermissionsForHostname(string) (permissions.Permissions, error)
	SavePermissions(string, permissions.Permissions) error
}

type GetPermissions struct {
	store PermissionsStore
}

func NewGetPermissions(store PermissionsStore) *GetPermissions {
	return &GetPermissions{
		store: store,
	}
}

func (s *GetPermissions) Handle(rawParams jsonrpc.Params) (jsonrpc.Result, *jsonrpc.ErrorDetails) {
	params, err := validateGetPermissionsParams(rawParams)
	if err != nil {
		return nil, invalidParams(err)
	}

	hostname, err := toHostname(params.Referer)
	if err != nil {
		return nil, invalidParams(err)
	}

	detailedPerms, err := s.store.PermissionsForHostname(hostname)
	if err != nil {
		return nil, internalError(fmt.Errorf("couldn't retrieve permissions: %w", err))
	}

	return GetPermissionsResult{
		Permissions: detailedPerms.Summary(),
	}, nil
}

func validateGetPermissionsParams(rawParams jsonrpc.Params) (GetPermissionsParams, error) {
	if rawParams == nil {
		return GetPermissionsParams{}, ErrParamsRequired
	}

	params, ok := rawParams.(GetPermissionsParams)
	if !ok {
		return GetPermissionsParams{}, ErrParamsDoNotMatch
	}

	if params.Referer == "" {
		return GetPermissionsParams{}, ErrRefererIsRequired
	}

	return params, nil
}

type GetPermissionsParams struct {
	Referer string `json:"referer"`
}

type GetPermissionsResult struct {
	Permissions map[string]string `json:"permissions"`
}

func toHostname(referer string) (string, error) {
	parsedReferer, err := url.Parse(referer)
	if err != nil {
		return "", ErrRefererIsNotValidURL
	}

	if parsedReferer.Hostname() == "" {
		return "", ErrRefererIsNotValidURL
	}

	return parsedReferer.Hostname(), nil
}
