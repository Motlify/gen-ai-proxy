package database

import (
	"context"
	"fmt"
	"gen-ai-proxy/src/config"
	"time"

	"github.com/jackc/pgx/v5"
)

func Connect(cfg *config.Config) (*pgx.Conn, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	pgxConfig, err := pgx.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database config: %w", err)
	}

	// Explicitly disable TLS
	pgxConfig.TLSConfig = nil

	maxAttempts := 5
	initialDelay := 1 * time.Second

	for i := 0; i < maxAttempts; i++ {
		conn, err := pgx.ConnectConfig(context.Background(), pgxConfig)
		if err == nil {
			return conn, nil // Successfully connected
		}

		// Log the error for debugging
		fmt.Printf("Attempt %d to connect to database failed: %v\n", i+1, err)

		if i < maxAttempts-1 {
			// Exponential backoff
			time.Sleep(initialDelay * time.Duration(1<<i))
		} else {
			return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxAttempts, err)
		}
	}
	return nil, fmt.Errorf("unexpected error: should not reach here")
}
