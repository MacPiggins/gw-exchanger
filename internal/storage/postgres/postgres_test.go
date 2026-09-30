package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func newTestDatabase(t *testing.T) (*PostgresDB, func()) {
	t.Helper()
	ctx := context.Background()

	container, err := postgrescontainer.Run(ctx, "postgres:16-alpine",
		postgrescontainer.WithDatabase("test"),
		postgrescontainer.WithUsername("test"),
		postgrescontainer.WithPassword("test"),
		postgrescontainer.WithInitScripts(
			"testdata/test.sql",
		),
		postgrescontainer.BasicWaitStrategies(),
	)
	require.NoError(t, err)

	connstr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	db, err := New(ctx, connstr)
	require.NoError(t, err)
	return db, func() {
		db.Close()
		require.NoError(t, container.Terminate(context.Background()))
	}
}

func TestPostgresDB(t *testing.T) {
	db, cleanup := newTestDatabase(t)
	defer cleanup()

	t.Run("GetExchangeRate", func(t *testing.T) {
		ctx := context.Background()
		rate, err := db.GetExchangeRate(ctx, "USD", "EUR")
		require.NoError(t, err)
		require.InDelta(t, 0.92, rate, 0.0001)
	})

	t.Run("GetExchangeRates", func(t *testing.T) {
		ctx := context.Background()
		rates, err := db.GetExchangeRates(ctx)
		require.NoError(t, err)
		require.Equal(t, map[string]float32{"USDtoEUR": 0.92, "EURtoUSD": 1.08}, rates)
	})

	t.Run("GetExchangeRateMissing", func(t *testing.T) {
		ctx := context.Background()
		_, err := db.GetExchangeRate(ctx, "GBP", "JPY")
		require.Error(t, err)
	})
}

func TestPostgresDBAutoMigrate(t *testing.T) {
	ctx := context.Background()
	container, err := postgrescontainer.Run(ctx, "postgres:16-alpine",
		postgrescontainer.WithDatabase("test"),
		postgrescontainer.WithUsername("test"),
		postgrescontainer.WithPassword("test"),
		postgrescontainer.BasicWaitStrategies(),
	)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, container.Terminate(context.Background()))
	}()

	connstr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	db, err := New(ctx, connstr)
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, db.AutoMigrate(ctx))
	require.NoError(t, db.AutoMigrate(ctx))
}
