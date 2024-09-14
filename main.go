package main

import (
	"flag"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	mw "github.com/LucDeCaf/go-simple-blog/middleware"
	"github.com/LucDeCaf/go-simple-blog/routes"
	"github.com/LucDeCaf/go-simple-blog/routes/v1"
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

	v1Router.Handle("/blogs", mw.Logging(v1.BlogsHandler).Build())
	v1Router.Handle("/blogs/{id}", mw.Logging(v1.BlogsIdHandler).Build())
	v1Router.Handle("/users", mw.Logging(v1.UsersHandler).Build())

	http.Handle("/v1/", http.StripPrefix("/v1", v1Router))

	http.Handle("/login", mw.Logging(routes.LoginHandler).Build())
	http.Handle("/register", mw.Logging(routes.RegisterHandler).Build())

	log.Println("api listening on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Print(err)
	}
}
