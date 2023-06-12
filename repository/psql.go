package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	host = "database"
	port = 5432
)

const createActivityTableIfNotQuery = `CREATE TABLE IF NOT EXISTS activity_log(
											id SERIAL PRIMARY KEY,
											key VARCHAR(100) NOT NULL,
											activity VARCHAR(100) NOT NULL,
											created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
										);`

// OpenPgSQL creates a new postgres-sql connection and returns.
func OpenPgSQL(username, password, database string) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, database)

	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "open pgsql")
	}

	err = conn.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "pgsql ping failed")
	}
	log.Println("Database connection established")
	return conn, nil
}

// CreateTableIfNot is responsible for creating the table if not already.
func CreateTableIfNot(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx,
		createActivityTableIfNotQuery)

	if err != nil {
		return errors.Wrap(err, "create table")
	}
	return nil
}
