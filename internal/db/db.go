package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(path string) (*sql.DB, error) {
	// TODO: Open database
	// TODO: Execute schema.sql
	// TODO: Return connection
}
