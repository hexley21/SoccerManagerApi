package service

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrIncorrectPassword = errors.New("incorrect user password")
	ErrUsernameTaken = errors.New("username is taken")

	ErrTeamNotFound = errors.New("team not found")
	
	ErrPlayerNotFound = errors.New("player not found")

	ErrTransferNotFound = errors.New("transfer not found")
	ErrPlayerAlreadyInTransfers = errors.New("player is already in transfers")
	ErrCantBuyFromYourself = errors.New("can't buy from yourself")
	ErrNotEnoughFunds = errors.New("not enough funds")

	ErrTransferRecordNotFound = errors.New("transfer record not found")

	ErrNonexistentCode = errors.New("nonexistent code or key")
	
	ErrTranslationNotFound = errors.New("translation not found")
	ErrTranslationExists = errors.New("translation already exists")

	ErrInvalidArguments = errors.New("invalid arguments")
)
