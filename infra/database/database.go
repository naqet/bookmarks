package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
)

func Init() *sql.DB {
	dbUser, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		slog.Error("POSTGRES_USER env is not set")
		os.Exit(1)
	}

	dbPass, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		slog.Error("POSTGRES_PASSWORD env is not set")
		os.Exit(1)
	}

	dbName, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		slog.Error("POSTGRES_DB env is not set")
		os.Exit(1)
	}

	conn, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", dbUser, dbPass, dbName))
	if err != nil {
		slog.Error("couldn't open a database", slog.Any("error", err))
	}

	conn.SetMaxOpenConns(1)

	return conn
}

func Migrate(db *sql.DB) error {
	content, err := os.ReadFile("infra/database/init.sql")

	if err != nil {
		slog.Error("migration failed", slog.Any("error", err))
		os.Exit(1)
	}

	_, err = db.Exec(string(content))
	return err
}
