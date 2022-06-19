package api

import (
	"fmt"

	"code.vegaprotocol.io/shared/paths"
	"code.vegaprotocol.io/vegawallet/handlers"
	"code.vegaprotocol.io/vegawallet/libs/jsonrpc"
	permissionsstore "code.vegaprotocol.io/vegawallet/permissions/store/v1"
)

// RestrictedAPI builds a JSON-RPC API of the wallet with a subset of the requests
// that are intended to be exposed to external services, such as bots, apps,
// scripts.
func RestrictedAPI(vegaPath paths.Paths) (*jsonrpc.API, error) {
	walletAPI := jsonrpc.New()

	permStore, err := permissionsstore.InitialiseStore(vegaPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialise permissions store: %w", err)
	}

	walletAPI.RegisterCommand("get_permissions", handlers.NewGetPermissions(permStore))
	walletAPI.RegisterCommand("request_permissions", handlers.NewRequestPermissions(permStore))

	return walletAPI, nil
}
