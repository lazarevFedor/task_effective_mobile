package postgres

import (
	"context"
	"errors"
	"fmt"
	"task_effective_mobile/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	Username string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Database string `env:"POSTGRES_DB"`

	MinConns int32 `env:"POSTGRES_MIN_CONNS"`
	MaxConns int32 `env:"POSTGRES_MAX_CONNS"`
}

// New returns pool of connections to postgres DB
func New(ctx context.Context, c Config, service string) (*pgxpool.Pool, error) {
	log := logger.GetLoggerFromCtx(ctx)
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
	log.Info(ctx, fmt.Sprintf("connected to %s_postgres", service))

	migration, err := migrate.New(
		fmt.Sprintf("file://db/migrations/%s", service),
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			c.Username,
			c.Password,
			c.Host,
			c.Port,
			c.Database,
		),
	)

	if err != nil {
		return nil, fmt.Errorf("new: failed to create migration instance: %w", err)
	}

	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("new: failed to Up migration: %w", err)
	}
	log.Info(ctx, "Successfully Applied Migration")
	return conn, nil
}
