package service

import (
	"sync"

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
	GetSignRequestsConfirmations(hash string) chan ConsentConfirmation
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

func (p *AutomaticConsentPolicy) GetSignRequestsConfirmations(hash string) chan ConsentConfirmation {
	return nil
}

type ExplicitConsentPolicy struct {
	pendingEvents    chan ConsentRequest
	AllConfirmations sync.Map
}

func NewExplicitConsentPolicy(pending chan ConsentRequest) Policy {
	return &ExplicitConsentPolicy{
		pendingEvents:    pending,
		AllConfirmations: sync.Map{},
	}
}

func (p *ExplicitConsentPolicy) GetSignRequestsConfirmations(hash string) chan ConsentConfirmation {
	confirmations, ok := p.AllConfirmations.Load(hash)
	if ok {
		return confirmations.(chan ConsentConfirmation)
	}
	return nil
}

func (p *ExplicitConsentPolicy) Ask(tx *v1.SubmitTransactionRequest) (bool, error) {
	txHash := crypto.AsSha256(tx)
	confirmations := make(chan ConsentConfirmation)
	p.AllConfirmations.Store(txHash, confirmations)
	p.pendingEvents <- ConsentRequest{tx: tx, Confirmations: confirmations}

	c := <-confirmations
	req := &v1.SubmitTransactionRequest{}
	if err := jsonpb.UnmarshalString(c.TxStr, req); err != nil {
		return false, ErrInvalidSignRequestConfirm
	}
	if crypto.AsSha256(req) != txHash {
		return false, ErrUnexpectedSignRequestConfirm
	}

	return c.Decision, nil
}
