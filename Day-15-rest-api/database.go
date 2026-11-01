package main

import (
	"context"
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "go_practice.db")
	if err != nil {
		return nil, err
	}

	// Keep the pool simple for this exercise.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
