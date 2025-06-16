package jwt

import "errors"

var (
	ErrErrorSigningToken = errors.New("error signing token")
	ErrErrorParsingToken = errors.New("error parsing token")
	ErrInvalidToken      = errors.New("invalid token")
)
