package service

import (
	v1 "code.vegaprotocol.io/protos/vega/wallet/v1"
)

type ConsentConfirmation struct {
	TxStr    string
	Decision bool
}

type ConsentRequest struct {
	tx *v1.SubmitTransactionRequest
}

func (r *ConsentRequest) String() string {
	return r.tx.String()
}

type Policy interface {
	Ask(tx *v1.SubmitTransactionRequest) bool
}

type AutomaticConsentPolicy struct {
	pendingEvents chan ConsentRequest
	confirmations chan ConsentConfirmation
}

func NewAutomaticConsentPolicy(pending chan ConsentRequest, response chan ConsentConfirmation) Policy {
	return &AutomaticConsentPolicy{
		pendingEvents: pending,
		confirmations: response,
	}
}

func (p *AutomaticConsentPolicy) Ask(tx *v1.SubmitTransactionRequest) bool {
	return true
}

type ExplicitConsentPolicy struct {
	pendingEvents chan ConsentRequest
	confirmations chan ConsentConfirmation
}

func NewExplicitConsentPolicy(pending chan ConsentRequest, response chan ConsentConfirmation) Policy {
	return &ExplicitConsentPolicy{
		pendingEvents: pending,
		confirmations: response,
	}
}

func (p *ExplicitConsentPolicy) Ask(tx *v1.SubmitTransactionRequest) bool {
	p.pendingEvents <- ConsentRequest{tx}
	txStr := tx.String()

	for c := range p.confirmations {
		if c.TxStr == txStr {
			return c.Decision
		}
	}

	return true
}
