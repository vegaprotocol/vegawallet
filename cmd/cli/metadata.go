package cli

import (
	"strings"

	"code.vegaprotocol.io/vegawallet/cmd/flags"
	"code.vegaprotocol.io/vegawallet/wallet"
)

func ParseMetadata(rawMetadata []string) ([]wallet.Meta, error) {
	if len(rawMetadata) == 0 {
		return nil, nil
	}

	metadata := make([]wallet.Meta, 0, len(rawMetadata))
	for _, v := range rawMetadata {
		rawMeta := strings.Split(v, ":")
		if len(rawMeta) != 2 { //nolint:gomnd
			return nil, flags.InvalidFlagFormatError("meta")
		}
		metadata = append(metadata, wallet.Meta{Key: rawMeta[0], Value: rawMeta[1]})
	}

	return metadata, nil
}
