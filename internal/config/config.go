// Package config provides loading of application configuration from environment
// variables. It uses the cleanenv package to populate a Config struct that
// contains Postgres connection settings and the server port.
package config

import (
	"fmt"
	"task_effective_mobile/pkg/postgres"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config holds application configuration read from environment variables.
//
// Fields are tagged for cleanenv so that environment variables like
// POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, etc. are automatically
// mapped into the nested Postgres.Config. The Port field is populated from
// SERVER_PORT.
type Config struct {
	Postgres postgres.Config `env:"POSTGRES"`
	Port     string          `env:"SERVER_PORT"`
}

// New reads configuration from environment variables and returns a populated
// Config instance. If reading environment variables fails the function
// returns an error describing the problem.
func New() (*Config, error) {
	var config Config
	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, fmt.Errorf("New: reading env error: %w", err)
	}
	return &config, nil
}
