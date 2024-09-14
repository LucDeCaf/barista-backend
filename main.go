package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	mw "github.com/LucDeCaf/go-simple-blog/middleware"
	"github.com/LucDeCaf/go-simple-blog/routes"
)

func main() {
	portPtr := flag.String(
		"port",
		"8080",
		"the port the application will run on",
	)

	flag.Parse()

	port := *portPtr

	v1 := http.NewServeMux()

	v1.Handle("/blogs", mw.Logging(routes.BlogsHandler).Build())
	v1.Handle("/blogs/{id}", mw.Logging(routes.BlogsIdHandler).Build())
	v1.Handle("/users", mw.Logging(routes.UsersHandler).Build())

	http.Handle("/v1/", http.StripPrefix("/v1", v1))

	http.Handle("/login", mw.Logging(routes.LoginHandler).Build())
	http.Handle("/register", mw.Logging(routes.RegisterHandler).Build())

	log.Println("api listening on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Print(err)
	}
}
