package v1

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/LucDeCaf/go-simple-blog/models/users"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		w.Write([]byte("method not allowed"))
		return fmt.Errorf("method not allowed: %v", r.Method)
	}

	switch r.Header.Get("Server-Action") {
	case "GetByUsername":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return err
		}

		username := string(body)

		user, err := users.Get(username)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				w.WriteHeader(404)
				w.Write([]byte("not found"))
			} else {
				w.WriteHeader(500)
				w.Write([]byte("internal server error"))
			}

			return err
		}

		resp, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return err
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)

	case "GetAll":
		users, err := users.GetAll()
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return err
		}

		resp, err := json.Marshal(users)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("internal server error"))
			return err
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)

	default:
		w.WriteHeader(400)
		w.Write([]byte("missing Server-Action header"))
		return fmt.Errorf("missing Server-Action header")
	}

	return nil
}
