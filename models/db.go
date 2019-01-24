package models

import (
	"database/sql"

	// imporing the mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// DataAccessor defines the methods to retrieve data
type DataAccessor interface {
	UserDataReader
	PostDataReader
}

// DB wraps the database object and will implement the interface DataAccessor
type DB struct {
	*sql.DB
}

// NewDB returns and instnace of the DB
func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
