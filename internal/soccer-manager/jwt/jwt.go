package jwt

import "time"

//go:generate mockgen -destination=mock/mock_jwt.go -package=mock github.com/hexley21/soccer-manager/internal/soccer-manager/jwt Manager,ManagerWithTTL
type Manager[T any] interface {
	CreateTokenString(data T) (string, error)
	ParseTokenString(token string) (T, error)
}

type ManagerWithTTL[T any] interface {
	Manager[T]
	TTL() time.Duration
}
