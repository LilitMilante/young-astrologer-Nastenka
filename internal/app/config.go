package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort    string `env:"HTTP_PORT"`
	PostgresDSN string `env:"POSTGRES_DSN"`
	NASAAPI     string `env:"NASA_API"`
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load env: %w", err)
	}

	cfg := &Config{}

	err = env.ParseWithOptions(cfg, env.Options{RequiredIfNoDef: true})
	if err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	return cfg, nil
}
