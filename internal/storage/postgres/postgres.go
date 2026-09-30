package postgres

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type PostgresDB struct {
	pool    *pgxpool.Pool
	connstr string
	close   chan struct{}
}

func New(ctx context.Context, connstr string) (*PostgresDB, error) {
	pool, err := pgxpool.New(ctx, connstr)
	if err != nil {
		slog.Error("error while creating postgres connection pool", slog.Any("error", err))
		return nil, err
	}
	close := make(chan struct{}, 1)
	db := &PostgresDB{pool: pool, connstr: connstr, close: close}
	go db.maintainConnection()
	return db, nil
}

func (db *PostgresDB) AutoMigrate(ctx context.Context) error {
	migrationDB := stdlib.OpenDBFromPool(db.pool)
	defer migrationDB.Close()
	migrations, err := fs.Sub(migrationFiles, "migrations")
	if err != nil {
		slog.Error("error while loading postgres migrations", slog.Any("error", err))
		db.pool.Close()
		return fmt.Errorf("load postgres migrations: %w", err)
	}
	provider, err := goose.NewProvider(goose.DialectPostgres, migrationDB, migrations)
	if err != nil {
		slog.Error("error while configuring postgres migrations", slog.Any("error", err))
		db.pool.Close()
		return fmt.Errorf("configure postgres migrations: %w", err)
	}
	if _, err := provider.Up(ctx); err != nil {
		slog.Error("error while running postgres migrations", slog.Any("error", err))
		db.pool.Close()
		return fmt.Errorf("run postgres migrations: %w", err)
	}
	return nil
}

func (db *PostgresDB) maintainConnection() {
	for {
		select {
		case <-db.close:
			db.pool.Close()
			return
		case <-time.After(5 * time.Second):
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			err := db.pool.Ping(ctx)
			if err != nil {
				slog.Error("postgres connection error, trying to reconnect", slog.Any("error", err))
				db.pool.Close()
				db.pool, err = pgxpool.New(ctx, db.connstr)
				if err != nil {
					slog.Error("postgres failed to reconnect, trying again in 5 sec", slog.Any("error", err))
				}
				slog.Info("postgres successfully reconnected")
			}
			cancel()
		}
	}
}

func (db *PostgresDB) Close() {
	close(db.close)
}

func (db *PostgresDB) GetExchangeRate(ctx context.Context, from, to string) (float32, error) {
	var rate float32
	row := db.pool.QueryRow(ctx, "select rate from Rates WHERE Cur_From = $1 and Cur_To = $2", from, to)
	err := row.Scan(&rate)
	if err != nil {
		slog.ErrorContext(ctx, "postgres error", slog.Any("error", err))
		return 0, err
	}
	return rate, nil
}

func (db *PostgresDB) GetExchangeRates(ctx context.Context) (map[string]float32, error) {
	rates := make(map[string]float32, 0)
	rows, err := db.pool.Query(ctx, "select Cur_From, Cur_To, rate from Rates;")
	if err != nil {
		slog.ErrorContext(ctx, "postgres error", slog.Any("error", err))
		return rates, err
	}
	defer rows.Close()
	for rows.Next() {
		var from string
		var to string
		var rate float32
		err := rows.Scan(&from, &to, &rate)
		if err != nil {
			slog.ErrorContext(ctx, "postgres error", slog.Any("error", err))
		}
		rates[fmt.Sprint(from, "to", to)] = rate
	}
	return rates, nil
}
