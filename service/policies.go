package service

import (
	"reflect"

	v1 "code.vegaprotocol.io/protos/vega/wallet/v1"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/protobuf/encoding/protojson"
)

type ConsentConfirmation struct {
	TxStr    string
	Decision bool
}

type ConsentRequest struct {
	tx *v1.SubmitTransactionRequest
}

func (r *ConsentRequest) String() (string, error) {
	data, err := protojson.Marshal(r.tx)
	return string(data), err
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

	for c := range p.confirmations {
		req := &v1.SubmitTransactionRequest{}
		if err := jsonpb.UnmarshalString(c.TxStr, req); err != nil {
			continue
		}

		if reflect.DeepEqual(req, tx) {
			return c.Decision
		}
	}
	return true
}
