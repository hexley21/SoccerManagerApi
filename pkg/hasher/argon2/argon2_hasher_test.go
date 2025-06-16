package argon2_test

import (
	"testing"

	"github.com/hexley21/soccer-manager/pkg/config"
	"github.com/hexley21/soccer-manager/pkg/hasher"
	"github.com/hexley21/soccer-manager/pkg/hasher/argon2"
	"github.com/stretchr/testify/assert"
)

const (
	hashLen        = 128
	normalPassword = "abcdefghijklmnopqrstuvwxyz123456789"
	crazyPassword  = "abcdefghijklmnopqrstuvwxyz123456789😀😀😀😀無無無無無無無無無"
)

var (
	argon2Hasher = argon2.NewHasher(config.Argon2{
		SaltLen:    16,
		KeyLen:     79,
		Time:       1,
		Memory:     47104,
		Threads:    1,
		Breakpoint: config.Argon2KeylenBreakpoint(79),
	})
)

func Test_HashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "Normal Password",
			password: normalPassword,
		},
		{
			name:     "Crazy Password",
			password: crazyPassword,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := argon2Hasher.HashPassword(tt.password)
			assert.NoError(t, err)
			assert.Len(t, hash, hashLen)
		})
	}
}

func Test_VerifyPassword(t *testing.T) {
	hash, err := argon2Hasher.HashPasswordWithSalt("pwd", "somesalt")
	if err != nil {
		t.Error(err)
	}

	// Take in mind, that this test will fail if hasher's config has incorrect breakpoint
	t.Run("OK", func(t *testing.T) {
		assert.NoError(t, argon2Hasher.VerifyPassword("pwd", hash))
	})

	t.Run("incorrect password", func(t *testing.T) {
		assert.ErrorIs(t, argon2Hasher.VerifyPassword("rand", hash), hasher.ErrPasswordMismatch)
	})
}
