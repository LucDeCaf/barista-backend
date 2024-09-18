package routes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/LucDeCaf/go-simple-blog/auth"
	"github.com/LucDeCaf/go-simple-blog/models/users"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}

	var lr loginRequest

	if err := json.NewDecoder(r.Body).Decode(&lr); err != nil {
		log.Println("err decoding body:", err.Error())
		http.Error(w, "bad request", 400)
		return
	}

	lr.Username = strings.TrimSpace(lr.Username)
	lr.Password = strings.TrimSpace(lr.Password)

	if lr.Username == "" || lr.Password == "" {
		http.Error(w, "incorrect username or password", 400)
		return
	}

	user, err := users.Get(lr.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "incorrect username or password", 400)
		} else {
			log.Println("err getting user:", err.Error())
			http.Error(w, "internal server error", 500)
		}

		return
	}

	if !auth.VerifyPassword(lr.Password, user.PasswordHash) {
		http.Error(w, "incorrect username or password", 400)
		return
	}

	token, err := auth.NewJWT(lr.Username)
	if err != nil {
		log.Println("err creating JWT:", err.Error())
		http.Error(w, "internal server error", 500)
		return
	}

	w.Write([]byte(token))

	return
}
