package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func NewPool(
	ctx context.Context,
	dbConfig DBConfig,
) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dbConfig.DatabaseURL())
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = dbConfig.MaxConns
	cfg.MinConns = dbConfig.MinConns
	cfg.MaxConnLifetime = dbConfig.MaxConnLifetime
	cfg.HealthCheckPeriod = dbConfig.HealthCheckPeriod

	return pgxpool.NewWithConfig(ctx, cfg)
}

func NewDB(
	ctx context.Context,
	dbConfig DBConfig,
) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbConfig.DatabaseURL()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(int(dbConfig.MaxConns))
	sqlDB.SetMaxIdleConns(int(dbConfig.MinConns))
	sqlDB.SetConnMaxLifetime(dbConfig.MaxConnLifetime)
	sqlDB.SetConnMaxIdleTime(dbConfig.HealthCheckPeriod)

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = sqlDB.PingContext(ctxTimeout); err != nil {
		return nil, err
	}

	return db, nil
}
