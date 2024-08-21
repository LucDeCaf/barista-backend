package db

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("sqlite3", "../db.db")
	if err != nil {
		log.Fatal(err)
	}
}
