package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func NewPool(
	ctx context.Context,
	dbConfig DBConfig,
) (*pgxpool.Pool, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName, dbConfig.SSLMode,
	)

	cfg, err := pgxpool.ParseConfig(url)
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

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dbConfig.Host,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.Port,
		dbConfig.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
