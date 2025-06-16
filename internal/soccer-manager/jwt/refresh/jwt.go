package refresh

import (
	"fmt"
	"time"

	jwtgo "github.com/golang-jwt/jwt/v5"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/jwt"
	"github.com/hexley21/soccer-manager/pkg/config"
)

type Manager = jwt.ManagerWithTTL[Data]

type Data struct {
	UserID int64 `json:"uid"`
}

func NewData(UserID int64) Data {
	return Data{UserID: UserID}
}

type Claims struct {
	Data
	jwtgo.RegisteredClaims
}

type tokenManager struct {
	secret string
	ttl    time.Duration
}

func NewManager(cfg config.TokenParams) *tokenManager {
	return &tokenManager{
		secret: cfg.Secret,
		ttl:    cfg.TTL,
	}
}

func (m *tokenManager) CreateTokenString(data Data) (string, error) {
	keyByte := []byte(m.secret)

	currTime := time.Now()

	claims := Claims{
		Data: data,
		RegisteredClaims: jwtgo.RegisteredClaims{
			ExpiresAt: jwtgo.NewNumericDate(currTime.Add(m.ttl)),
			NotBefore: jwtgo.NewNumericDate(currTime),
		},
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	ss, err := token.SignedString(keyByte)
	if err != nil {
		return "", fmt.Errorf("%w: %w", jwt.ErrErrorSigningToken, err)
	}

	return ss, nil
}

func (m *tokenManager) ParseTokenString(tokenString string) (Data, error) {
	token, err := jwtgo.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwtgo.Token) (interface{}, error) {
			return []byte(m.secret), nil
		},
	)
	if err != nil {
		return Data{}, fmt.Errorf("%w: %w", jwt.ErrErrorParsingToken, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return Data{}, jwt.ErrInvalidToken
	}

	return claims.Data, nil
}

func (m *tokenManager) TTL() time.Duration {
	return m.ttl
}
