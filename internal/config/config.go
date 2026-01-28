package config

import (
	"fmt"
	"task_effective_mobile/pkg/postgres"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres postgres.Config `env:"POSTGRES"`
	Port     string          `env:"SERVER_PORT"`
}

func New() (*Config, error) {
	var config Config
	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, fmt.Errorf("New: reading env error: %w", err)
	}
	return &config, nil
}
