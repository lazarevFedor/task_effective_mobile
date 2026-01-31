package postgres

import (
	"context"
	"fmt"
	"task_effective_mobile/pkg/logger"

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
