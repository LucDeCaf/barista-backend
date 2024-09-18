package v1

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/LucDeCaf/go-simple-blog/models/users"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}

	switch r.Header.Get("Server-Action") {
	case "GetByUsername":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("err reading body:", err.Error())
			http.Error(w, "internal server error", 500)
			return
		}

		username := string(body)

		user, err := users.Get(username)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "not found", 404)
			} else {
				http.Error(w, "internal server error", 500)
			}

			log.Println("err getting user:", err.Error())
			return
		}

		resp, err := json.Marshal(user)
		if err != nil {
			log.Println("err converting user to json:", err.Error())
			http.Error(w, "internal server error", 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)

	case "GetAll":
		users, err := users.GetAll()
		if err != nil {
			log.Println("err getting all users:", err.Error())
			http.Error(w, "internal server error", 500)
			return
		}

		resp, err := json.Marshal(users)
		if err != nil {
			log.Println("err converting user to json:", err.Error())
			http.Error(w, "internal server error", 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)

	default:
		http.Error(w, "missing Server-Action header", 400)
		return
	}

	return
}
