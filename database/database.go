package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/scys12/go-attendance/config"
)

const (
	dbCredentialsFormat = "%v:%v@tcp(%v:%v)/%v"
)

func GetDatabaseConnection(cfg config.Config) *sql.DB {
	address := fmt.Sprintf(dbCredentialsFormat,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	address = address + "?parseTime=True&loc=Asia%2FJakarta"

	db, err := sql.Open("mysql", address)
	if err != nil {
		log.Fatal("[Database] failed connecting to DB: " + address + ", err: " + err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Fatal("[Database] db is unreachable: " + address + ", err: " + err.Error())
	}

	return db
}
