package wallets

import (
	"fmt"

	wstorev1 "code.vegaprotocol.io/vegawallet/wallet/store/v1"
	"code.vegaprotocol.io/shared/paths"
)

// InitialiseStore builds a wallet Store specifically for users wallets.
func InitialiseStore(vegaHome string) (*wstorev1.Store, error) {
	p := paths.New(vegaHome)
	walletsHome, err := p.CreateDataPathFor(paths.WalletsDataHome)
	if err != nil {
		return nil, fmt.Errorf("couldn't get wallets data home path: %w", err)
	}
	return wstorev1.InitialiseStore(walletsHome)
}
