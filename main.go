package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"

	"github.com/LucDeCaf/go-simple-blog/auth"
	"github.com/LucDeCaf/go-simple-blog/errors"
	mw "github.com/LucDeCaf/go-simple-blog/middleware"
	"github.com/LucDeCaf/go-simple-blog/models/author"
	"github.com/LucDeCaf/go-simple-blog/models/blog"
	"github.com/LucDeCaf/go-simple-blog/models/user"
)

var db *sql.DB

type loginRequest struct {
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

	v1 := http.NewServeMux()
	v1.Handle("/authors", middlewares(authorsHandler))
	v1.Handle("/authors/{id}", middlewares(authorsIdHandler))
	v1.Handle("/blogs", middlewares(blogsHandler))
	v1.Handle("/blogs/{id}", middlewares(blogsIdHandler))
	v1.Handle("/login", middlewares(loginHandler))

	http.Handle("/v1/", http.StripPrefix("/v1", v1))

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

func extractUser(r *http.Request) (*user.User, *errors.HttpError) {
	token, err := auth.ExtractJWT(r)
	if err != nil {
		return nil, errors.NewHttpError(401, "unauthenticated")
	}

	claims, err := auth.ExtractClaims(token)
	if err != nil {
		return nil, errors.NewHttpError(401, "unauthenticated")
	}

	userTable := user.NewUserTable(db)
	user, err := userTable.Get(claims.Username)
	if err != nil {
		return nil, errors.NewHttpError(500, "internal server error")
	}

	return &user, nil
}

func authorsHandler(w http.ResponseWriter, r *http.Request) error {
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
		if user, err := extractUser(r); err != nil {
			httpRespond(w, err.Code, err.Message)
			return err
		} else if user.Role != "admin" {
			httpRespond(w, 403, "unauthorized")
			return err
		}

		de := json.NewDecoder(r.Body)
		de.DisallowUnknownFields()

		var a *author.Author

		if err := de.Decode(&a); err != nil {
			httpRespond(w, 400, "bad request")
			return err
		}

		a, err := authorTable.Insert(a)
		if err != nil {
			httpRespond(w, 500, "internal server error")
			return err
		}

		w.WriteHeader(201)
		if err := json.NewEncoder(w).Encode(a); err != nil {
			return err
		}
	default:
		httpRespond(w, 405, "method not allowed")
		return fmt.Errorf("method '%v' not allowed", r.Method)
	}

	return nil
}

func blogsHandler(w http.ResponseWriter, r *http.Request) error {
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
		if user, err := extractUser(r); err != nil {
			httpRespond(w, err.Code, err.Message)
			return err
		} else if user.Role != "admin" {
			httpRespond(w, 403, "unauthorized")
			return err
		}

		de := json.NewDecoder(r.Body)
		de.DisallowUnknownFields()

		var b *blog.Blog

		if err := de.Decode(&b); err != nil {
			httpRespond(w, 400, "bad request")
			return err
		}

		blogTable := blog.NewBlogTable(db)
		b, err := blogTable.Insert(b)
		if err != nil {
			httpRespond(w, 500, "internal server error")
			return err
		}

		w.WriteHeader(201)
		if err := json.NewEncoder(w).Encode(b); err != nil {
			return err
		}
	default:
		httpRespond(w, 405, "method not allowed")
		return fmt.Errorf("method '%v' not allowed", r.Method)
	}

	return nil
}

func authorsIdHandler(w http.ResponseWriter, r *http.Request) error {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpRespond(w, 400, "bad request")
		return err
	}

	switch r.Method {
	case http.MethodGet:
		authorTable := author.NewAuthorTable(db)
		a, err := authorTable.Get(id)
		if err != nil {
			httpRespond(w, 404, "not found")
			return err
		}

		if err = json.NewEncoder(w).Encode(a); err != nil {
			return err
		}

	default:
		httpRespond(w, 405, "method not allowed")
		return fmt.Errorf("method '%v' not allowed", r.Method)
	}

	return nil
}

func blogsIdHandler(w http.ResponseWriter, r *http.Request) error {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpRespond(w, 400, "bad request")
		return err
	}

	switch r.Method {
	case http.MethodGet:
		blogTable := blog.NewBlogTable(db)
		b, err := blogTable.Get(id)
		if err != nil {
			httpRespond(w, 404, "not found")
			return err
		}

		if err = json.NewEncoder(w).Encode(b); err != nil {
			return err
		}

	default:
		httpRespond(w, 405, "method not allowed")
		return fmt.Errorf("method '%v' not allowed", r.Method)
	}

	return nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		httpRespond(w, 405, "method not allowed")
		return fmt.Errorf("invalid method %v", r.Method)
	}

	var lr loginRequest

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
