package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	mw "github.com/LucDeCaf/go-simple-blog/middleware"
	"github.com/LucDeCaf/go-simple-blog/models/author"
	"github.com/LucDeCaf/go-simple-blog/models/blog"
	"github.com/LucDeCaf/go-simple-blog/models/user"
)

var db *sql.DB

type LoginRequest struct {
	Username string `json:"username"`
	Uassword string `json:"password"`
}

func init() {
	var err error

	db, err = sql.Open("sqlite3", "./db.db")
	if err != nil {
		log.Fatal("failed to open db:", err)
	}
}

func main() {
	portPtr := flag.String(
		"port",
		"8080",
		"the port the application will run on",
	)

	flag.Parse()

	port := *portPtr

	http.Handle("/author", middlewares(authorHandler))
	http.Handle("/blog", middlewares(blogHandler))
	http.Handle("/login", middlewares(loginHandler))

	log.Println("api listening on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Print(err)
	}
}

func middlewares(h mw.Handler) http.Handler {
	return mw.Build(mw.GenericErrorWriter(mw.ErrorLogger(h)))
}

func authorHandler(w http.ResponseWriter, r *http.Request) error {
	authorTable := author.NewAuthorTable(db)

	switch r.Method {
	case http.MethodGet:
		authors, err := authorTable.GetAll()
		if err != nil {
			return err
		}

		en := json.NewEncoder(w)
		if err := en.Encode(authors); err != nil {
			return err
		}
	case http.MethodPost:
		de := json.NewDecoder(r.Body)
		de.DisallowUnknownFields()

		var a author.Author

		if err := de.Decode(&a); err != nil {
			return err
		}

		// Echo back request body
		en := json.NewEncoder(w)
		if err := en.Encode(a); err != nil {
			return err
		}
	}

	return nil
}

func blogHandler(w http.ResponseWriter, r *http.Request) error {
	blogTable := blog.NewBlogTable(db)

	switch r.Method {
	case http.MethodGet:
		blogs, err := blogTable.GetAll()
		if err != nil {
			return err
		}

		en := json.NewEncoder(w)
		if err := en.Encode(blogs); err != nil {
			return err
		}
	case http.MethodPost:
		de := json.NewDecoder(r.Body)
		de.DisallowUnknownFields()

		var a author.Author

		if err := de.Decode(&a); err != nil {
			return err
		}

		// Echo back request body
		en := json.NewEncoder(w)
		if err := en.Encode(a); err != nil {
			return err
		}
	}

	return nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("invalid method %v", r.Method)
	}

	var lr LoginRequest

	de := json.NewDecoder(r.Body)
	if err := de.Decode(&lr); err != nil {
		return err
	}

	userTable := user.NewUserTable(db)

	user, err := userTable.Get(lr.Username)
	if err != nil {
		return err
	}

	en := json.NewEncoder(w)
	if err := en.Encode(user); err != nil {
		return err
	}

	return nil
}
