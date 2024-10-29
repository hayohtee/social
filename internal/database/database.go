package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// New create an return a sql.DB connection to a PostgreSQL database.
func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, fmt.Errorf("error parsing maxIdleTime: %w", err)
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxIdleConns)
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
