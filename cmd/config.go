package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

type HTTPConfig struct {
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"5s"`
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"`
}

type DBConfig struct {
	URL               string        `env:"URL" envDefault:""`
	Host              string        `env:"HOST" envDefault:""`
	Port              int           `env:"PORT" envDefault:"0"`
	User              string        `env:"USER" envDefault:""`
	Password          string        `env:"PASSWORD" envDefault:""`
	DBName            string        `env:"NAME" envDefault:""`
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

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c Config) ListenAddr() string {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		return c.HTTPAddr
	}

	return net.JoinHostPort("", port)
}

func (db DBConfig) DatabaseURL() string {
	if strings.TrimSpace(db.URL) != "" {
		return strings.TrimSpace(db.URL)
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.DBName,
		db.SSLMode,
	)
}

func (c Config) validate() error {
	if strings.TrimSpace(c.DB.URL) != "" {
		return nil
	}

	switch {
	case strings.TrimSpace(c.DB.Host) == "":
		return fmt.Errorf("DB_HOST is required when DB_URL is empty")
	case c.DB.Port == 0:
		return fmt.Errorf("DB_PORT is required when DB_URL is empty")
	case strings.TrimSpace(c.DB.User) == "":
		return fmt.Errorf("DB_USER is required when DB_URL is empty")
	case strings.TrimSpace(c.DB.Password) == "":
		return fmt.Errorf("DB_PASSWORD is required when DB_URL is empty")
	case strings.TrimSpace(c.DB.DBName) == "":
		return fmt.Errorf("DB_NAME is required when DB_URL is empty")
	default:
		return nil
	}
}
