package env

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPostgresConfig(t *testing.T) {
	t.Setenv("POSTGRES_HOST", "localhost")
	t.Setenv("POSTGRES_PORT", "5435")
	t.Setenv("POSTGRES_USER", "order_user")
	t.Setenv("POSTGRES_PASSWORD", "order_password")
	t.Setenv("POSTGRES_DB", "order")
	t.Setenv("POSTGRES_SSL_MODE", "disable")
	t.Setenv("MIGRATION_DIRECTORY", "./order/migrations")

	cfg, err := NewPostgresConfig()
	require.NoError(t, err)

	require.Equal(t,
		"postgres://order_user:order_password@localhost:5435/order?sslmode=disable",
		cfg.URI(),
	)
	require.Equal(t, "./order/migrations", cfg.MigrationsDir())
}

func TestOrderHTTPConfig(t *testing.T) {
	t.Setenv("HTTP_HOST", "localhost")
	t.Setenv("HTTP_PORT", "8080")
	t.Setenv("HTTP_READ_TIMEOUT", "5s")

	cfg, err := NewOrderHTTPConfig()
	require.NoError(t, err)

	require.Equal(t, "localhost:8080", cfg.Address())
	require.Equal(t, 5*time.Second, cfg.ReadTimeout())
}

func TestOrderHTTPConfig_BadTimeout(t *testing.T) {
	t.Setenv("HTTP_HOST", "localhost")
	t.Setenv("HTTP_PORT", "8080")
	t.Setenv("HTTP_READ_TIMEOUT", "5")

	_, err := NewOrderHTTPConfig()
	require.Error(t, err)
}
