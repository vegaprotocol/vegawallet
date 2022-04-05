package service

import (
	"crypto/sha256"
	"fmt"

	v1 "code.vegaprotocol.io/protos/vega/wallet/v1"
	"github.com/golang/protobuf/jsonpb"
)

type ConsentConfirmation struct {
	TxStr    string
	Decision bool
}

type ConsentRequest struct {
	Tx            *v1.SubmitTransactionRequest
	Confirmations chan ConsentConfirmation
}

func (r *ConsentRequest) String() (string, error) {
	m := jsonpb.Marshaler{Indent: "    "}
	marshalledRequest, err := m.MarshalToString(r.Tx)
	return marshalledRequest, err
}

func (r *ConsentRequest) TxHash() string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s%v%t", r.Tx.PubKey, r.Tx.Command, r.Tx.Propagate)))

	return fmt.Sprintf("%x", h.Sum(nil))
}

type Policy interface {
	Ask(tx *v1.SubmitTransactionRequest) (bool, error)
}

type AutomaticConsentPolicy struct{}

func NewAutomaticConsentPolicy() Policy {
	return &AutomaticConsentPolicy{}
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
	p.pendingEvents <- ConsentRequest{Tx: tx, Confirmations: confirmations}

	c := <-confirmations
	return c.Decision, nil
}
