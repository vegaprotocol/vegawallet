package wallet

import "errors"

var (
	ErrWalletDoesNotExists = errors.New("wallet does not exist")
	ErrWalletAlreadyExists = errors.New("a wallet with the same name already exists")
)
