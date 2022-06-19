package handlers

//
// import (
// 	"errors"
// 	"fmt"
// 	"time"
//
// 	"code.vegaprotocol.io/vegawallet/libs/jsonrpc"
// )
//
// var (
// 	ErrUserRejectedCommand = errors.New("user rejected the command")
// )
//
// type CommandToReview struct {
// 	TraceID    string
// 	ReceivedAt time.Time
// 	Command    SendTransactionParams
// }
//
// type ReviewCommandFn func(CommandToReview) (bool, error)
//
// type TransactionStatus struct {
// 	TraceID string
// 	TxHash  string
// 	Tx      string
// 	Error   error
// 	SentAt  time.Time
// }
//
// type ReportTransactionStatusFn func(TransactionStatus)
//
// type SendTransaction struct {
// 	reviewCommandFn ReviewCommandFn
// 	reportStatusFn  ReportTransactionStatusFn
// }
//
// func NewSendTransaction(reviewCommandFn ReviewCommandFn, reportStatusFn ReportTransactionStatusFn) *SendTransaction {
// 	if reviewCommandFn == nil {
// 		panic("function to review commands can't be nil")
// 	}
//
// 	if reportStatusFn == nil {
// 		panic("function to report status can't be nil")
// 	}
//
// 	return &SendTransaction{
// 		reviewCommandFn: reviewCommandFn,
// 		reportStatusFn:  reportStatusFn,
// 	}
// }
//
// func (s *SendTransaction) Handle(request *jsonrpc.Request) *jsonrpc.Response {
// 	params, err := validateSendTransactionParams(request)
// 	if err != nil {
// 		return invalidParamsResponse(request.ID, err)
// 	}
//
// 	cmdToReview := bundleCommandToReview(params)
//
// 	_, err = s.reviewCommandFn(cmdToReview)
// 	if err != nil {
// 		return userErrorResponse(request.ID, fmt.Errorf("couldn't review command: %w", err))
// 	}
//
// 	return jsonrpc.NewSuccessfulResponse(request.ID, struct {
// 		TxHash     string    `json:"txHash"`
// 		ReceivedAt time.Time `json:"receivedAt"`
// 		SentAt     time.Time `json:"sentAt"`
// 		TxID       string    `json:"txId"`
// 		Tx         string    `json:"tx"`
// 	}{
// 		// TxHash:     txHash,
// 		// ReceivedAt: receivedAt,
// 		// SentAt:     sentAt,
// 		// TxID:       txID,
// 		// Tx:         tx,
// 	})
// }
//
// func validateSendTransactionParams(request *jsonrpc.Request) (SendTransactionParams, error) {
// 	if request.Params == nil {
// 		return SendTransactionParams{}, ErrParamsRequired
// 	}
//
// 	params, ok := request.Params.(SendTransactionParams)
// 	if !ok {
// 		return SendTransactionParams{}, ErrParamsDoNotMatch
// 	}
//
// 	// TODO Add params validation
//
// 	return params, nil
// }
//
// func bundleCommandToReview(params SendTransactionParams) CommandToReview {
// 	return CommandToReview{
// 		TraceID:    "123456789",
// 		ReceivedAt: time.Now(),
// 		Command:    params,
// 	}
// }
//
// type SendTransactionParams struct {
// 	PublicKey string `json:"pubkey"`
// 	Tx        string `json:"tx"`
// }
