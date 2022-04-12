package mocks

import (
	v1 "code.vegaprotocol.io/protos/vega/wallet/v1"
	"code.vegaprotocol.io/vegawallet/service"
)

type MockConsentPolicy struct {
	pendingEvents chan service.ConsentRequest
}

func NewMockConsentPolicy(pending chan service.ConsentRequest) service.Policy {
	return &MockConsentPolicy{
		pendingEvents: pending,
	}
}

func (p *MockConsentPolicy) Ask(tx *v1.SubmitTransactionRequest) (bool, error) {
	if tx.PubKey == "toBeDeclined" {
		return false, nil
	}
	return true, nil
}
