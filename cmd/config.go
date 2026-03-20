package main

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type HTTPConfig struct {
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"5s"`
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"`
}

type DBConfig struct {
	Host              string        `env:"HOST"`
	Port              int           `env:"PORT"`
	User              string        `env:"USER"`
	Password          string        `env:"PASSWORD"`
	DBName            string        `env:"NAME"`
	SSLMode           string        `env:"SSLMODE" envDefault:"disable"`
	MaxConns          int32         `env:"MAX_CONNS" envDefault:"10"`
	MinConns          int32         `env:"MIN_CONNS" envDefault:"2"`
	MaxConnLifetime   time.Duration `env:"MAX_CONN_LIFETIME" envDefault:"30m"`
	HealthCheckPeriod time.Duration `env:"HEALTH_CHECK_PERIOD" envDefault:"1m"`
}

type Config struct {
	HTTPAddr string     `env:"HTTP_ADDR" envDefault:":8080"`
	HTTP     HTTPConfig `envPrefix:"HTTP_"`
	DB       DBConfig   `envPrefix:"DB_"`
}

func ReadEnv() (*Config, error) {
	cfg := new(Config)
	opts := env.Options{
		RequiredIfNoDef: true,
	}

	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return nil, err
	}

	return cfg, nil
}
