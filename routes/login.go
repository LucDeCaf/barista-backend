package routes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/LucDeCaf/go-simple-blog/auth"
	"github.com/LucDeCaf/go-simple-blog/models/users"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		HttpRespond(w, 405, "method not allowed")
		return fmt.Errorf("invalid method %v", r.Method)
	}

	var lr loginRequest

	if err := json.NewDecoder(r.Body).Decode(&lr); err != nil {
		HttpRespond(w, 400, "bad request")
		return err
	}

	lr.Username = strings.TrimSpace(lr.Username)
	lr.Password = strings.TrimSpace(lr.Password)

	if lr.Username == "" || lr.Password == "" {
		HttpRespond(w, 400, "incorrect username or password")
		return fmt.Errorf("missing username or password ('%v', '%v')", lr.Username, lr.Password)
	}

	user, err := users.Get(lr.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			HttpRespond(w, 400, "incorrect username or password")
		} else {
			HttpRespond(w, 500, "internal server error")
		}
		return err
	}

	if !auth.VerifyPassword(lr.Password, user.PasswordHash) {
		HttpRespond(w, 400, "incorrect username or password")
		return fmt.Errorf("incorrect password")
	}

	token, err := auth.NewJWT(lr.Username)
	if err != nil {
		HttpRespond(w, 500, "internal server error")
		return err
	}

	w.Write([]byte(token))

	return nil
}
