package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/LucDeCaf/go-simple-blog/routes"
	v1 "github.com/LucDeCaf/go-simple-blog/routes/v1"
)

func main() {
	portPtr := flag.String(
		"port",
		"8080",
		"the port the application will run on",
	)

	flag.Parse()

	port := *portPtr

	v1Router := http.NewServeMux()

	v1Router.HandleFunc("/blogs", v1.BlogsHandler)
	v1Router.HandleFunc("/blogs/{id}", v1.BlogsIdHandler)
	v1Router.HandleFunc("/users", v1.UsersHandler)

	http.Handle("/v1/", http.StripPrefix("/v1", v1Router))

	http.HandleFunc("/login", routes.LoginHandler)
	http.HandleFunc("/register", routes.RegisterHandler)

	log.Println("api listening on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Print(err)
	}
}
