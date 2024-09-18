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

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", 405)
		return
	}

	var rr registerRequest
	if err := json.NewDecoder(r.Body).Decode(&rr); err != nil {
		log.Println("err decoding request:", err.Error())
		http.Error(w, "bad request", 400)
		return
	}

	// Remove leading or trailing whitespace
	rr.Username = strings.TrimSpace(rr.Username)
	rr.Password = strings.TrimSpace(rr.Password)

	// Prevent empty username / password
	if rr.Username == "" || rr.Password == "" {
		log.Println("missing username or password")
		http.Error(w, "missing username or password", 400)
		return
	}

	// Password requirements
	if len(rr.Password) < 8 {
		log.Println("password too short")
		http.Error(w, "password must be at least 8 characters", 400)
		return
	}

	_, err := users.Get(rr.Username)

	// Successful return means username already in use
	if err == nil {
		http.Error(w, "username already exists", 400)
		return
	}

	/*
		If ErrNoRows is returned then that means username can be used
		and therefore ErrNoRows is actually a success case.

		If any other error occurs, then it is an actual server error.
	*/
	if !errors.Is(err, sql.ErrNoRows) {
		log.Println("err checking for user:", err.Error())
		http.Error(w, "internal server error", 500)
		return
	}

	pwHash, err := auth.HashPassword(rr.Password)
	if err != nil {
		log.Println("err hashing password:", err.Error())
		http.Error(w, "internal server error", 500)
		return
	}

	user, err := users.Insert(users.User{
		Username:     rr.Username,
		PasswordHash: pwHash,
		Role:         "user",
	})
	if err != nil {
		log.Println("err creating user:", err.Error())
		http.Error(w, "internal server error", 500)
		return
	}

	w.WriteHeader(201)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Println("err encoding user to json:", err.Error())
		return
	}

	return
}
