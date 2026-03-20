package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

type HTTPConfig struct {
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" envDefault:"5s"`
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"`
}

type CORSConfig struct {
	AllowedOrigins []string `env:"ALLOWED_ORIGINS" envSeparator:"," envDefault:"*"`
	AllowedMethods []string `env:"ALLOWED_METHODS" envSeparator:"," envDefault:"GET,OPTIONS"`
	AllowedHeaders []string `env:"ALLOWED_HEADERS" envSeparator:"," envDefault:"Authorization,Content-Type,Accept,Origin,X-Request-ID"`
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
	CORS     CORSConfig `envPrefix:"CORS_"`
	DB       DBConfig   `envPrefix:"DB_"`

	LegacyCORSAllowOrigins []string `env:"CORS_ALLOW_ORIGINS" envSeparator:","`
}

func ReadEnv() (*Config, error) {
	cfg := new(Config)
	opts := env.Options{
		RequiredIfNoDef: true,
	}
	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return nil, err
	}
	cfg.applyLegacyCompatibility()
	return cfg, nil
}

func (c Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.DBName, c.DB.SSLMode,
	)
}

func (c *Config) applyLegacyCompatibility() {
	if strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS")) == "" && len(c.LegacyCORSAllowOrigins) > 0 {
		c.CORS.AllowedOrigins = c.LegacyCORSAllowOrigins
	}
}
