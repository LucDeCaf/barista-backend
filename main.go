package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/LucDeCaf/go-simple-blog/auth"
	mw "github.com/LucDeCaf/go-simple-blog/middleware"
	"github.com/LucDeCaf/go-simple-blog/models/author"
	"github.com/LucDeCaf/go-simple-blog/models/blog"
	"github.com/LucDeCaf/go-simple-blog/models/user"
)

var db *sql.DB

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
	return mw.Build(mw.RequestLogger(mw.ErrorLogger(h)))
}

func httpRespond(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf("%v %v\n", code, message)))
}

func authorHandler(w http.ResponseWriter, r *http.Request) error {
	authorTable := author.NewAuthorTable(db)

	switch r.Method {
	case http.MethodGet:
		authors, err := authorTable.GetAll()
		if err != nil {
			httpRespond(w, 500, "internal server error")
			return err
		}

		en := json.NewEncoder(w)
		if err := en.Encode(authors); err != nil {
			httpRespond(w, 500, "internal server error")
			return err
		}
	case http.MethodPost:
		de := json.NewDecoder(r.Body)
		de.DisallowUnknownFields()

		var a author.Author

		if err := de.Decode(&a); err != nil {
			httpRespond(w, 400, "bad request")
			return err
		}

		en := json.NewEncoder(w)
		if err := en.Encode(a); err != nil {
			httpRespond(w, 500, "internal server error")
			return err
		}
	default:
		httpRespond(w, 405, "method not allowed")
	}

	return nil
}

func blogHandler(w http.ResponseWriter, r *http.Request) error {
	blogTable := blog.NewBlogTable(db)

	switch r.Method {
	case http.MethodGet:
		blogs, err := blogTable.GetAll()
		if err != nil {
			httpRespond(w, 500, "internal server error")
			return err
		}

		en := json.NewEncoder(w)
		if err := en.Encode(blogs); err != nil {
			httpRespond(w, 500, "internal server error")
			return err
		}
	case http.MethodPost:
		de := json.NewDecoder(r.Body)
		de.DisallowUnknownFields()

		var a author.Author

		if err := de.Decode(&a); err != nil {
			httpRespond(w, 400, "bad request")
			return err
		}

		en := json.NewEncoder(w)
		if err := en.Encode(a); err != nil {
			httpRespond(w, 500, "internal server error")
			return err
		}
	default:
		httpRespond(w, 405, "method not allowed")
	}

	return nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		httpRespond(w, 405, "method not allowed")
		return fmt.Errorf("invalid method %v", r.Method)
	}

	var lr LoginRequest

	de := json.NewDecoder(r.Body)
	if err := de.Decode(&lr); err != nil {
		httpRespond(w, 400, "bad request")
		return err
	}

	userTable := user.NewUserTable(db)

	user, err := userTable.Get(lr.Username)
	if err != nil {
		httpRespond(w, 500, "internal server error")
		return err
	}

	if !auth.VerifyPassword(lr.Password, user.PasswordHashWithSalt) {
		httpRespond(w, 401, "unauthenticated")
		return fmt.Errorf("wrong password inputted")
	}

	token, err := auth.NewJWT(lr.Username)
	if err != nil {
		httpRespond(w, 500, "internal server error")
		return err
	}

	w.Write([]byte(fmt.Sprintf(`{"token":"%v"}`, token)))

	return nil
}
