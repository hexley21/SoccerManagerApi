package hasher

import "errors"

var ErrPasswordMismatch = errors.New("password does not match")

//go:generate mockgen -destination=mock/mock_hasher.go -package=mock github.com/hexley21/soccer-manager/pkg/hasher Hasher
type Hasher interface {
	HashPassword(password string) (string, error)
	HashPasswordWithSalt(password string, salt string) (string, error)
	VerifyPassword(password string, hash string) error
	GetSalt() ([]byte, error)
}
