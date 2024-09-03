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

	v1.Handle("/authors", middlewares(routes.AuthorsHandler))
	v1.Handle("/authors/{id}", middlewares(routes.AuthorsIdHandler))
	v1.Handle("/blogs", middlewares(routes.BlogsHandler))
	v1.Handle("/blogs/{id}", middlewares(routes.BlogsIdHandler))

	http.Handle("/v1/", http.StripPrefix("/v1", v1))

	http.Handle("/login", middlewares(routes.LoginHandler))
	http.Handle("/register", middlewares(routes.RegisterHandler))

	log.Println("api listening on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Print(err)
	}
}

func middlewares(h mw.Handler) http.Handler {
	return mw.Build(mw.RequestLogger(mw.ErrorLogger(h)))
}
