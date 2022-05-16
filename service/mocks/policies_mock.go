package mocks

import (
	"time"

	v1 "code.vegaprotocol.io/protos/vega/wallet/v1"
	"code.vegaprotocol.io/vegawallet/service"
)

type MockConsentPolicy struct {
	pendingEvents chan service.ConsentRequest
	sentTxs       chan service.SentTransaction
}

func (p *MockConsentPolicy) Report(tx service.SentTransaction) {
	p.sentTxs <- tx
}

func NewMockConsentPolicy(pending chan service.ConsentRequest, sentTxs chan service.SentTransaction) *MockConsentPolicy {
	return &MockConsentPolicy{
		pendingEvents: pending,
		sentTxs:       sentTxs,
	}
}

func (p *MockConsentPolicy) Ask(tx *v1.SubmitTransactionRequest, txID string, receivedAt time.Time) (bool, error) {
	if tx.PubKey == "toBeDeclined" {
		return false, nil
	}
	return true, nil
}

func (p *MockConsentPolicy) NeedsInteractiveOutput() bool {
	return true
}
