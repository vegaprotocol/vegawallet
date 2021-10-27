package wallet_test

import (
	"testing"

	"code.vegaprotocol.io/vegawallet/wallet"
	"github.com/stretchr/testify/assert"
)

func TestVersionIsSupported(t *testing.T) {
	tcs := []struct {
		name      string
		version   uint32
		supported bool
	}{
		{
			name:      "version 0",
			version:   0,
			supported: false,
		}, {
			name:      "version 1",
			version:   1,
			supported: true,
		}, {
			name:      "version 2",
			version:   2,
			supported: true,
		}, {
			name:      "version 3",
			version:   3,
			supported: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			// when
			supported := wallet.IsVersionSupported(tc.version)

			assert.Equal(tt, tc.supported, supported)
		})
	}
}
