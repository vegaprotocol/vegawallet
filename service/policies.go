package service

import (
	"time"

	commandspb "code.vegaprotocol.io/protos/vega/commands/v1"
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

type SentTransaction struct {
	TxHash       string
	TxID         string
	ReceivedAt   time.Time
	Tx           *commandspb.Transaction
	Error        error
	ErrorDetails []string
}

type Policy interface {
	Ask(tx *v1.SubmitTransactionRequest, txID string, receivedAt time.Time) (bool, error)
	Report(tx SentTransaction)
	NeedsInteractiveOutput() bool
}

type AutomaticConsentPolicy struct{}

func NewAutomaticConsentPolicy() Policy {
	return &AutomaticConsentPolicy{}
}

func (p *AutomaticConsentPolicy) Ask(_ *v1.SubmitTransactionRequest, txID string, receivedAt time.Time) (bool, error) {
	return true, nil
}

func (p *AutomaticConsentPolicy) Report(_ SentTransaction) {
	// Nothing to report as we expect this policy to be non-interactive.
}

func (p *AutomaticConsentPolicy) NeedsInteractiveOutput() bool {
	return false
}

type ExplicitConsentPolicy struct {
	pendingEvents chan ConsentRequest
	sentTxs       chan SentTransaction
}

func NewExplicitConsentPolicy(pending chan ConsentRequest, sentTxs chan SentTransaction) Policy {
	return &ExplicitConsentPolicy{
		pendingEvents: pending,
		sentTxs:       sentTxs,
	}
}

func (p *ExplicitConsentPolicy) Ask(tx *v1.SubmitTransactionRequest, txID string, receivedAt time.Time) (bool, error) {
	confirmations := make(chan ConsentConfirmation)
	consentReq := ConsentRequest{Tx: tx, Confirmations: confirmations, ReceivedAt: receivedAt}
	consentReq.TxID = txID
	p.pendingEvents <- consentReq

	c := <-confirmations
	return c.Decision, nil
}

func (p *ExplicitConsentPolicy) Report(tx SentTransaction) {
	p.sentTxs <- tx
}

func (p *ExplicitConsentPolicy) NeedsInteractiveOutput() bool {
	return true
}
