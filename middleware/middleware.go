package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/LucDeCaf/go-simple-blog/auth"
)

// A type identical to the signature of http.HandlerFunc,
// but with the benefit of being able to return errors to
// be handled by middleware functions.
type Handler func(w http.ResponseWriter, r *http.Request) error

// Converts middleware.Handler into http.Handler
func Build(next Handler) http.Handler {
	// Create wrapper around h
	f := func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
	return http.HandlerFunc(f)
}

func RequestLogger(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// Log request
		log.Printf("%v %v\n", r.Method, r.URL.Path)
		return next(w, r)
	}
}

func ErrorLogger(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// Log errors
		err := next(w, r)
		if err != nil {
			log.Println("err:", err)
		}
		return err
	}
}

func Protected(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			return fmt.Errorf("missing authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return fmt.Errorf("invalid authorization header format")
		}

		_, err := auth.ExtractClaims(parts[1])
		if err != nil {
			return err
		}

		// TODO: Use claims to verify user identity and permissions

		return next(w, r)
	}
}

func GenericErrorWriter(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := next(w, r)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("500 internal server error\n"))
		}
		return err
	}
}
