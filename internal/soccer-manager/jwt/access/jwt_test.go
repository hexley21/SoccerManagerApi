package access_test

import (
	"testing"
	"time"

	jwtgo "github.com/golang-jwt/jwt/v5"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/domain"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/jwt"
	"github.com/hexley21/soccer-manager/internal/soccer-manager/jwt/access"
	"github.com/hexley21/soccer-manager/pkg/config"
	"github.com/stretchr/testify/assert"
)

func Test_CreateTokenString(t *testing.T) {
	manager := access.NewManager(config.TokenParams{Secret: "secret", TTL: time.Hour})

	t.Run("OK", func(t *testing.T) {
		token, err := manager.CreateTokenString(access.NewData(1234, domain.UserRoleADMIN))

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func Test_ParseTokenString(t *testing.T) {
	manager := access.NewManager(config.TokenParams{Secret: "secret", TTL: time.Hour})

	t.Run("OK", func(t *testing.T) {
		data := access.NewData(1234, domain.UserRoleADMIN)

		token, err := manager.CreateTokenString(data)
		assert.NoError(t, err)

		parsedData, err := manager.ParseTokenString(token)
		assert.NoError(t, err)
		assert.Equal(t, data, parsedData)
	})

	t.Run("invalid signKey", func(t *testing.T) {
		token, err := manager.CreateTokenString(access.NewData(1234, domain.UserRoleADMIN))
		assert.NoError(t, err)

		managerInvalidKey := access.NewManager(
			config.TokenParams{Secret: "invalidKey", TTL: time.Hour},
		)

		parsedData, err := managerInvalidKey.ParseTokenString(token)
		assert.Error(t, err)
		assert.Empty(t, parsedData)
		assert.ErrorIs(t, err, jwt.ErrErrorParsingToken)
	})

	t.Run("invalid token", func(t *testing.T) {
		parsedData, err := manager.ParseTokenString("invalidToken")
		assert.Error(t, err)
		assert.Empty(t, parsedData)

		parsedData, err = manager.ParseTokenString("")
		assert.Error(t, err)
		assert.Empty(t, parsedData)
		assert.ErrorIs(t, err, jwt.ErrErrorParsingToken)
	})

	t.Run("expired token", func(t *testing.T) {
		keyByte := []byte("secret")

		currTime := time.Now()

		claims := access.Claims{
			Data: access.NewData(1234, domain.UserRoleADMIN),
			RegisteredClaims: jwtgo.RegisteredClaims{
				ExpiresAt: jwtgo.NewNumericDate(currTime.Add(-time.Hour)),
				IssuedAt:  jwtgo.NewNumericDate(currTime),
			},
		}

		token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
		ss, err := token.SignedString(keyByte)
		assert.NoError(t, err)

		parsedData, err := manager.ParseTokenString(ss)
		assert.Error(t, err)
		assert.Empty(t, parsedData)
		assert.ErrorIs(t, err, jwt.ErrErrorParsingToken)
	})
}

func Test_NewManager(t *testing.T) {
	manager := access.NewManager(config.TokenParams{Secret: "secret", TTL: time.Hour})

	assert.Equal(t, manager.TTL(), time.Hour)
	assert.NotNil(t, manager.CreateTokenString)
	assert.NotNil(t, manager.ParseTokenString)
}
