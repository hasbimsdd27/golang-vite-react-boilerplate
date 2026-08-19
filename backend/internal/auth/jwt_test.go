package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTService(t *testing.T) {
	secret := "test-secret-key"
	service := NewJWTService(secret)

	t.Run("generates valid token", func(t *testing.T) {
		claims := Claims{
			UserID:   1,
			Username: "testuser",
		}

		token, err := service.GenerateToken(claims)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("validates valid token", func(t *testing.T) {
		claims := Claims{
			UserID:   1,
			Username: "testuser",
		}

		token, _ := service.GenerateToken(claims)
		parsedClaims, err := service.ValidateToken(token)

		assert.NoError(t, err)
		assert.Equal(t, uint(1), parsedClaims.UserID)
		assert.Equal(t, "testuser", parsedClaims.Username)
	})

	t.Run("rejects invalid token", func(t *testing.T) {
		_, err := service.ValidateToken("invalid-token")

		assert.Error(t, err)
	})

	t.Run("rejects expired token", func(t *testing.T) {
		expiredService := &JWTService{
			secret:     []byte(secret),
			expiration: -1 * time.Hour,
		}

		claims := Claims{
			UserID:   1,
			Username: "testuser",
		}

		token, _ := expiredService.GenerateToken(claims)
		_, err := service.ValidateToken(token)

		assert.Error(t, err)
	})
}
