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
	Host              string        `env:"HOST"`
	Port              int           `env:"PORT"`
	User              string        `env:"USER"`
	Password          string        `env:"PASSWORD"`
	DBName            string        `env:"NAME"`
	SSLMode           string        `env:"SSLMODE"`
	MaxConns          int32         `env:"MAX_CONNS"`
	MinConns          int32         `env:"MIN_CONNS"`
	MaxConnLifetime   time.Duration `env:"MAX_CONN_LIFETIME"`
	HealthCheckPeriod time.Duration `env:"HEALTH_CHECK_PERIOD"`
}

type JWTConfig struct {
	Secret     string        `env:"SECRET"`
	Issuer     string        `env:"ISSUER"`
	Audience   string        `env:"AUDIENCE"`
	AccessTTL  time.Duration `env:"ACCESS_TTL"`
	RefreshTTL time.Duration `env:"REFRESH_TTL"`
}

type StorageServiceConfig struct {
	URL     string        `env:"URL" envDefault:"http://localhost:8085"`
	Timeout time.Duration `env:"TIMEOUT" envDefault:"30m"`
}

type NotificationServiceConfig struct {
	URL            string        `env:"URL" envDefault:"http://localhost:8087"`
	Timeout        time.Duration `env:"TIMEOUT" envDefault:"5s"`
	InternalAPIKey string        `env:"INTERNAL_API_KEY"`
}

type RedisConfig struct {
	Addr          string        `env:"ADDR" envDefault:""`
	Password      string        `env:"PASSWORD" envDefault:""`
	DB            int           `env:"DB" envDefault:"0"`
	KeyPrefix     string        `env:"KEY_PREFIX" envDefault:"zerde:curriculum"`
	DictionaryTTL time.Duration `env:"DICTIONARY_TTL" envDefault:"24h"`
}

type Config struct {
	HTTPAddr     string                    `env:"HTTP_ADDR" envDefault:":8080"`
	HTTP         HTTPConfig                `envPrefix:"HTTP_"`
	DB           DBConfig                  `envPrefix:"DB_"`
	JWT          JWTConfig                 `envPrefix:"JWT_"`
	Redis        RedisConfig               `envPrefix:"REDIS_"`
	Storage      StorageServiceConfig      `envPrefix:"STORAGE_SERVICE_"`
	Notification NotificationServiceConfig `envPrefix:"NOTIFICATION_SERVICE_"`
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

func (c Config) ListenAddr() string {
	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		return c.HTTPAddr
	}

	return net.JoinHostPort("", port)
}

func (db DBConfig) DatabaseURL() string {
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
