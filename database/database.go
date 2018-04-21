package database

import (
	"database/sql"

	// _ подключение пакета для работы с PostgreSQL
	_ "github.com/lib/pq"

	"github.com/Slava12/Computer_Market/config"
)

var db *sql.DB

// Connect осуществляет подключение к базе данных
func Connect(dbConfig config.Config) (err error) {
	connStr := "user=" + dbConfig.Database.User + " password=" + dbConfig.Database.Password + " dbname=" + dbConfig.Database.DBname + " host=" + dbConfig.Database.Host + " port=" + dbConfig.Database.Port + " sslmode=" + dbConfig.Database.SSLmode
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
