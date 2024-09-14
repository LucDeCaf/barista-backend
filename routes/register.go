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

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		HttpRespond(w, 405, "method not allowed")
		return fmt.Errorf("method '%v' not allowed", r.Method)
	}

	var rr registerRequest
	if err := json.NewDecoder(r.Body).Decode(&rr); err != nil {
		HttpRespond(w, 400, "bad request")
		return err
	}

	// Remove leading or trailing whitespace
	rr.Username = strings.TrimSpace(rr.Username)
	rr.Password = strings.TrimSpace(rr.Password)

	// Prevent empty username / password
	if rr.Username == "" || rr.Password == "" {
		HttpRespond(w, 400, "missing username or password")
		return fmt.Errorf("missing username or password ('%v', '%v')", rr.Username, rr.Password)
	}

	// Password requirements
	if len(rr.Password) < 8 {
		HttpRespond(w, 400, "password must be at least 8 characters")
		return fmt.Errorf("password too short")
	}

	_, err := users.Get(rr.Username)

	// Successful return means username already in use
	if err == nil {
		HttpRespond(w, 400, "username already exists")
		return fmt.Errorf("username already exists")
	}

	/*
		If ErrNoRows is returned then that means username can be used
		and therefore ErrNoRows is actually a success case.

		If any other error occurs, then it is an actual server error.
	*/
	if !errors.Is(err, sql.ErrNoRows) {
		HttpRespond(w, 500, "internal server error")
		return err
	}

	pwHash, err := auth.HashPassword(rr.Password)
	if err != nil {
		HttpRespond(w, 500, "internal server error")
		return err
	}

	user, err := users.Insert(users.User{
		Username:     rr.Username,
		PasswordHash: pwHash,
		Role:         "user",
	})
	if err != nil {
		HttpRespond(w, 500, "internal server error")
		return err
	}

	w.WriteHeader(201)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		return err
	}

	return nil
}
