package wallet

import "errors"

var (
	ErrInvalidMnemonic      = errors.New("mnemonic is not valid")
	ErrPubKeyAlreadyTainted = errors.New("public key is already tainted")
	ErrPubKeyIsTainted      = errors.New("public key is tainted")
	ErrPubKeyNotTainted     = errors.New("public key is not tainted")
	ErrPubKeyDoesNotExist   = errors.New("public key does not exist")
	ErrWalletAlreadyExists  = errors.New("a wallet with the same name already exists")
	ErrWalletDoesNotExists  = errors.New("wallet does not exist")
	ErrWalletNotLoggedIn    = errors.New("wallet is not logged in")
)
