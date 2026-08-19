package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	t.Run("creates config with provided values", func(t *testing.T) {
		cfg := NewConfig("host", "user", "password", "dbname", "5432")

		assert.Equal(t, "host", cfg.Host)
		assert.Equal(t, "user", cfg.User)
		assert.Equal(t, "password", cfg.Password)
		assert.Equal(t, "dbname", cfg.DBName)
		assert.Equal(t, "5432", cfg.Port)
	})
}

func TestConfigDSN(t *testing.T) {
	t.Run("returns correct DSN string", func(t *testing.T) {
		cfg := NewConfig("localhost", "postgres", "secret", "inventory", "5432")

		expected := "host=localhost user=postgres password=secret dbname=inventory port=5432 sslmode=disable"
		assert.Equal(t, expected, cfg.DSN())
	})
}
