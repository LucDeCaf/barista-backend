package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("sqlite3", "file:db.db?_foreign_keys=on")
	if err != nil {
		log.Fatal(err)
	}
}
