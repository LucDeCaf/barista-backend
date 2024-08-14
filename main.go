package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/LucDeCaf/go-simple-blog/models/author"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqliet3", "./db.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/author", func(w http.ResponseWriter, r *http.Request) {
		// Create new author instance
		a := author.NewAuthor(1, "hello", "world")

		// Create JSON encoder and write response
		wr := json.NewEncoder(w)
		if err := wr.Encode(a); err != nil {
			log.Println(err)
		}
	})

	log.Fatal(http.ListenAndServe(":http", nil))
}
