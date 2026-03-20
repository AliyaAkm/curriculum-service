package config

import (
	"slices"
	"testing"
)

func TestReadEnvSupportsLegacyCORSOrigins(t *testing.T) {
	t.Setenv("HTTP_ADDR", ":8083")
	t.Setenv("HTTP_READ_HEADER_TIMEOUT", "5s")
	t.Setenv("HTTP_SHUTDOWN_TIMEOUT", "5s")
	t.Setenv("CORS_ALLOWED_ORIGINS", "")
	t.Setenv("CORS_ALLOW_ORIGINS", "http://localhost:3000,http://localhost:5173")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "postgres")
	t.Setenv("DB_NAME", "curriculum")
	t.Setenv("DB_SSLMODE", "disable")

	cfg, err := ReadEnv()
	if err != nil {
		t.Fatalf("read env: %v", err)
	}

	expected := []string{"http://localhost:3000", "http://localhost:5173"}
	if !slices.Equal(cfg.CORS.AllowedOrigins, expected) {
		t.Fatalf("unexpected cors origins: %#v", cfg.CORS.AllowedOrigins)
	}
}
