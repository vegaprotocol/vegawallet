package wallet

import (
	"encoding/base64"
)

func (h *Handler) SignTx(name, tx, pubKey string, blockHeight uint64) (SignedBundle, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	rawTx, err := base64.StdEncoding.DecodeString(tx)
	if err != nil {
		return SignedBundle{}, err
	}

	w, err := h.loggedWallets.Get(name)
	if err != nil {
		return SignedBundle{}, err
	}

	return w.SignTxV1(pubKey, rawTx, blockHeight)
}
