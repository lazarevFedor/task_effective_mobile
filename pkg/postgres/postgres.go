// Package postgres provides a small helper for creating a pgx connection
// pool from environment-driven configuration. It also contains the Config
// struct used to populate connection settings via environment variables.
package postgres

import (
	"context"
	"fmt"
	"task_effective_mobile/pkg/logger"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Config describes connection parameters for a Postgres instance.
//
// The fields are tagged for mapping from environment variables (e.g.
// POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASSWORD,
// POSTGRES_DB) by the application's configuration loader.
type Config struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	Username string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Database string `env:"POSTGRES_DB"`

	MinConns int32 `env:"POSTGRES_MIN_CONNS"`
	MaxConns int32 `env:"POSTGRES_MAX_CONNS"`
}

// New creates and returns a pgx connection pool configured according to c.
//
// The provided context is used for pool creation and the service parameter
// is used for logging context only. Returned pool should be closed by the
// caller when no longer needed.
func New(ctx context.Context, c Config, service string) (*pgxpool.Pool, error) {
	log := logger.GetLogger(ctx)
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&pool_min_conns=%d&pool_max_conns=%d",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.MinConns,
		c.MaxConns)
	conn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("new: failed to connect to postgres: %w", err)
	}
	log.Info(fmt.Sprintf("connected to %s", service))
	return conn, nil
}
