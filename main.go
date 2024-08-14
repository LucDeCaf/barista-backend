package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/LucDeCaf/go-simple-blog/author"
)

func main() {
	http.HandleFunc("/author", func(w http.ResponseWriter, r *http.Request) {
		// Create new author instance
		a := author.Author{
			Id:        1,
			Firstname: "hello",
			Lastname:  "world",
		}

		// Create JSON encoder and write response
		wr := json.NewEncoder(w)
		if err := wr.Encode(a); err != nil {
			log.Println(err)
		}
	})

	log.Fatal(http.ListenAndServe(":http", nil))
}
