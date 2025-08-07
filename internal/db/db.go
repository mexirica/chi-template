// Package db handles database connection establishment and session management.
//
// Handles Database connection establishment
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

var dbConn = &DB{}

const maxDBLifetime = 5 * time.Minute

func Connect(connString string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	conn, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if err := otelpgx.RecordStats(conn); err != nil {
		return nil, fmt.Errorf("unable to record database stats: %w", err)
	}
	return conn, nil
}
