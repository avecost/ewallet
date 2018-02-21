package models

import (
	"database/sql"

	_ "github.com/lib/pq" // we are using PostgreSQL DB
)

// DB contains Database Object
type DB struct {
	*sql.DB
}

// NewDB return a DB Object using the dataSourceName
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
