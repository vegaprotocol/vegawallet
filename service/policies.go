package service

import (
	v1 "code.vegaprotocol.io/protos/vega/wallet/v1"
	"github.com/golang/protobuf/jsonpb"
)

type ConsentConfirmation struct {
	TxStr    string
	Decision bool
}

type ConsentRequest struct {
	tx            *v1.SubmitTransactionRequest
	Confirmations chan ConsentConfirmation
}

func (r *ConsentRequest) String() (string, error) {
	m := jsonpb.Marshaler{Indent: "    "}
	marshalledRequest, err := m.MarshalToString(r.tx)
	return marshalledRequest, err
}

type Policy interface {
	Ask(tx *v1.SubmitTransactionRequest) (bool, error)
}

type AutomaticConsentPolicy struct {
	pendingEvents chan ConsentRequest
}

func NewAutomaticConsentPolicy(pending chan ConsentRequest) Policy {
	return &AutomaticConsentPolicy{
		pendingEvents: pending,
	}
}

func (p *AutomaticConsentPolicy) Ask(_ *v1.SubmitTransactionRequest) (bool, error) {
	return true, nil
}

type ExplicitConsentPolicy struct {
	pendingEvents chan ConsentRequest
}

func NewExplicitConsentPolicy(pending chan ConsentRequest) Policy {
	return &ExplicitConsentPolicy{
		pendingEvents: pending,
	}
}

func (p *ExplicitConsentPolicy) Ask(tx *v1.SubmitTransactionRequest) (bool, error) {
	confirmations := make(chan ConsentConfirmation)
	p.pendingEvents <- ConsentRequest{tx: tx, Confirmations: confirmations}

	c := <-confirmations
	return c.Decision, nil
}
