package service

import (
	v1 "code.vegaprotocol.io/protos/vega/wallet/v1"
	"code.vegaprotocol.io/vegawallet/crypto"
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

func (p *AutomaticConsentPolicy) Ask(tx *v1.SubmitTransactionRequest) (bool, error) {
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
	req := &v1.SubmitTransactionRequest{}
	if err := jsonpb.UnmarshalString(c.TxStr, req); err != nil {
		return false, ErrInvalidSignRequestConfirm
	}
	if crypto.AsSha256(req) != crypto.AsSha256(tx) {
		return false, ErrUnexpectedSignRequestConfirm
	}

	return c.Decision, nil
}
