package service

import (
	"context"
	"time"

	commandspb "code.vegaprotocol.io/protos/vega/commands/v1"
	v1 "code.vegaprotocol.io/protos/vega/wallet/v1"
)

type ConsentConfirmation struct {
	TxID     string
	Decision bool
}

type ConsentRequest struct {
	TxID         string
	Tx           *v1.SubmitTransactionRequest
	ReceivedAt   time.Time
	Confirmation chan ConsentConfirmation
}

type SentTransaction struct {
	TxHash string
	TxID   string
	Tx     *commandspb.Transaction
	Error  error
	SentAt time.Time
}

type Policy interface {
	Ask(tx *v1.SubmitTransactionRequest, txID string, receivedAt time.Time) (bool, error)
	Report(tx SentTransaction)
}

type AutomaticConsentPolicy struct{}

func NewAutomaticConsentPolicy() Policy {
	return &AutomaticConsentPolicy{}
}

func (p *AutomaticConsentPolicy) Ask(_ *v1.SubmitTransactionRequest, _ string, _ time.Time) (bool, error) {
	return true, nil
}

func (p *AutomaticConsentPolicy) Report(_ SentTransaction) {
	// Nothing to report as we expect this policy to be non-interactive.
}

type ExplicitConsentPolicy struct {
	// ctx is used to interrupt the wait for consent confirmation
	ctx context.Context

	consentRequestsChan  chan ConsentRequest
	sentTransactionsChan chan SentTransaction
}

func NewExplicitConsentPolicy(ctx context.Context, consentRequests chan ConsentRequest, sentTransactions chan SentTransaction) Policy {
	return &ExplicitConsentPolicy{
		ctx:                  ctx,
		consentRequestsChan:  consentRequests,
		sentTransactionsChan: sentTransactions,
	}
}

func (p *ExplicitConsentPolicy) Ask(tx *v1.SubmitTransactionRequest, txID string, receivedAt time.Time) (bool, error) {
	confirmationChan := make(chan ConsentConfirmation, 1)
	defer close(confirmationChan)

	consentRequest := ConsentRequest{
		TxID:         txID,
		Tx:           tx,
		ReceivedAt:   receivedAt,
		Confirmation: confirmationChan,
	}

	if err := p.sendConsentRequest(consentRequest); err != nil {
		return false, err
	}

	return p.receiveConsentConfirmation(consentRequest)
}

func (p *ExplicitConsentPolicy) receiveConsentConfirmation(consentRequest ConsentRequest) (bool, error) {
	for {
		select {
		case <-p.ctx.Done():
			return false, ErrInterruptedConsentRequest
		case decision := <-consentRequest.Confirmation:
			return decision.Decision, nil
		}
	}
}

func (p *ExplicitConsentPolicy) sendConsentRequest(consentRequest ConsentRequest) error {
	for {
		select {
		case <-p.ctx.Done():
			return ErrInterruptedConsentRequest
		case p.consentRequestsChan <- consentRequest:
			return nil
		}
	}
}

func (p *ExplicitConsentPolicy) Report(tx SentTransaction) {
	p.sentTransactionsChan <- tx
}
