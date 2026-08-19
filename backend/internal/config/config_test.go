package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("loads config from environment variables", func(t *testing.T) {
		os.Setenv("DB_HOST", "testhost")
		os.Setenv("DB_USER", "testuser")
		os.Setenv("DB_PASSWORD", "testpass")
		os.Setenv("DB_NAME", "testdb")
		os.Setenv("DB_PORT", "5433")
		os.Setenv("PORT", "8080")
		os.Setenv("JWT_SECRET", "test-secret-key")
		os.Setenv("STATIC_DIR", "/custom/static/dir")
		defer func() {
			os.Unsetenv("DB_HOST")
			os.Unsetenv("DB_USER")
			os.Unsetenv("DB_PASSWORD")
			os.Unsetenv("DB_NAME")
			os.Unsetenv("DB_PORT")
			os.Unsetenv("PORT")
			os.Unsetenv("JWT_SECRET")
			os.Unsetenv("STATIC_DIR")
		}()

		cfg := Load()

		assert.Equal(t, "testhost", cfg.DBHost)
		assert.Equal(t, "testuser", cfg.DBUser)
		assert.Equal(t, "testpass", cfg.DBPassword)
		assert.Equal(t, "testdb", cfg.DBName)
		assert.Equal(t, "5433", cfg.DBPort)
		assert.Equal(t, "8080", cfg.Port)
		assert.Equal(t, "test-secret-key", cfg.JWTSecret)
		assert.Equal(t, "/custom/static/dir", cfg.StaticDir)
	})

	t.Run("uses default values when env vars not set", func(t *testing.T) {
		os.Clearenv()

		cfg := Load()

		assert.Equal(t, "localhost", cfg.DBHost)
		assert.Equal(t, "postgres", cfg.DBUser)
		assert.Equal(t, "", cfg.DBPassword)
		assert.Equal(t, "inventory", cfg.DBName)
		assert.Equal(t, "5432", cfg.DBPort)
		assert.Equal(t, "8080", cfg.Port)
		assert.Equal(t, "../frontend/dist", cfg.StaticDir)
	})
}
