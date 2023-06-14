package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	host              = "db"
	port              = 5432
	tableActivityLogs = "activity_logs"
)

const createActivityTableIfNotQuery = `CREATE TABLE IF NOT EXISTS activity_logs(
											id SERIAL PRIMARY KEY,
											key VARCHAR(100) NOT NULL,
											activity VARCHAR(100) NOT NULL,
											created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
										);`

// OpenPgSQL creates a new postgres-sql connection and returns.
func OpenPgSQL() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "open pgsql")
	}

	err = conn.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "pgsql ping failed")
	}
	CreateTableIfNot(context.Background(), conn)
	return conn, nil
}

// CreateTableIfNot is responsible for creating the table if not already.
func CreateTableIfNot(ctx context.Context, db *sql.DB) {
	_, err := db.ExecContext(ctx,
		createActivityTableIfNotQuery)

	if err != nil {
		log.Fatal().Msg("table creation failed")
	}
}
