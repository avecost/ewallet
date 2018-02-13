package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// DBApp struct will hold the pointer
// of the Open DB
type DBApp struct {
	*sql.DB
}

// Connect will Open DB connection
// using the provided credentials
func Connect(user, passwd, dbname string) (*DBApp, error) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		user, passwd, dbname)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &DBApp{db}, nil
}
