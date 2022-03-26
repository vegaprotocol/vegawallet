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
	Ask(tx *v1.SubmitTransactionRequest) bool
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

func (p *AutomaticConsentPolicy) Ask(tx *v1.SubmitTransactionRequest) bool {
	return true
}

func (p *AutomaticConsentPolicy) GetSignRequestsConfirmations(hash string) chan ConsentConfirmation {
	return nil
}

type ExplicitConsentPolicy struct {
	pendingEvents    chan ConsentRequest
	AllConfirmations map[string]chan ConsentConfirmation
}

func NewExplicitConsentPolicy(pending chan ConsentRequest) Policy {
	return &ExplicitConsentPolicy{
		pendingEvents:    pending,
		AllConfirmations: make(map[string]chan ConsentConfirmation),
	}
}

func (p *ExplicitConsentPolicy) GetSignRequestsConfirmations(hash string) chan ConsentConfirmation {
	return p.AllConfirmations[hash]
}

func (p *ExplicitConsentPolicy) Ask(tx *v1.SubmitTransactionRequest) bool {
	txHash := crypto.AsSha256(tx)
	p.AllConfirmations[txHash] = make(chan ConsentConfirmation)

	p.pendingEvents <- ConsentRequest{tx: tx, Confirmations: p.AllConfirmations[txHash]}

	for c := range p.AllConfirmations[txHash] {
		req := &v1.SubmitTransactionRequest{}
		if err := jsonpb.UnmarshalString(c.TxStr, req); err != nil {
			continue
		}
		if crypto.AsSha256(req) == txHash {
			return c.Decision
		}
	}
	return true
}
