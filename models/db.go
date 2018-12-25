package models

import (
	"database/sql"

	// imporing the mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// IDataStore defines the methods to retrieve data from the database
type IDataStore interface {
	GetUserByUsername(username string) (user User, err error)
}

// DB wraps the database object and will implement the interface IDataStore
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
