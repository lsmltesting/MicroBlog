package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/lsmltesting/MicroBlog/internal/logger"
)

const defaultMaxConns = int32(4)
const defaultMinConns = int32(0)
const defaultMaxConnLifetime = time.Hour
const defaultMaxConnIdleTime = time.Minute * 30
const defaultHealthCheckPeriod = time.Minute
const defaultConnectTimeout = time.Second * 5

func NewPostgresPool(lg logger.Logger) (*pgxpool.Pool, error) {
	wd, _ := os.Getwd()
	envPath := filepath.Join(wd, "..", ".env")
	_ = godotenv.Load(envPath)

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	lg.AddLog(
		logger.LevelDebug,
		logger.SourceDB,
		map[string]string{},
		fmt.Sprintf("dbUrl is equal -> %v", dbURL),
	)

	dbConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
		return nil, err
	}
	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		log.Fatal("Do not create pgxpool: %w", err)
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("No ping: %w", err)
		return nil, err
	}

	return pool, nil
}
