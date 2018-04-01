package database

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/Slava12/Computer_Market/config"
)

var db *sql.DB

func Connect(dbConfig config.Config) (err error) {
	connStr := "user=postgres password=Vbhfrk_)! dbname=market host=localhost port=5432 sslmode=disable"
	if db == nil {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			return err
		}
	}
	if err = db.Ping(); err != nil {
		return err
	}
	return nil
}
