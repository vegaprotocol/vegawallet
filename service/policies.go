package service

import (
	"crypto/sha256"
	"fmt"
	"time"

	v1 "code.vegaprotocol.io/protos/vega/wallet/v1"
	"github.com/golang/protobuf/jsonpb"
)

type ConsentConfirmation struct {
	TxStr    string
	Decision bool
}

type ConsentRequest struct {
	TxID          string
	Tx            *v1.SubmitTransactionRequest
	ReceivedAt    time.Time
	Confirmations chan ConsentConfirmation
}

func (r *ConsentRequest) String() (string, error) {
	m := jsonpb.Marshaler{Indent: "    "}
	marshalledRequest, err := m.MarshalToString(r.Tx)
	return marshalledRequest, err
}

func (r *ConsentRequest) GetTxID() string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s%v%t", r.Tx.PubKey, r.Tx.Command, r.Tx.Propagate)))

	return fmt.Sprintf("%x", h.Sum(nil))
}

type Policy interface {
	Ask(tx *v1.SubmitTransactionRequest) (bool, error)
	NeedsInteractiveOutput() bool
}

type AutomaticConsentPolicy struct{}

func NewAutomaticConsentPolicy() Policy {
	return &AutomaticConsentPolicy{}
}

func (p *AutomaticConsentPolicy) Ask(_ *v1.SubmitTransactionRequest) (bool, error) {
	return true, nil
}

func (p *AutomaticConsentPolicy) NeedsInteractiveOutput() bool {
	return false
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
	consentReq := ConsentRequest{Tx: tx, Confirmations: confirmations, ReceivedAt: time.Now()}
	consentReq.TxID = consentReq.GetTxID()
	p.pendingEvents <- consentReq

	c := <-confirmations
	return c.Decision, nil
}

func (p *ExplicitConsentPolicy) NeedsInteractiveOutput() bool {
	return true
}
