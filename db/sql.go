package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "user:password@/db")
	if err != nil {
		return nil, err
	}

	return db, nil
}
