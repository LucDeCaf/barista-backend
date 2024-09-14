package middleware

import (
	"fmt"
	"log"
	"net/http"
)

// A type identical to the signature of http.HandlerFunc,
// but with the benefit of being able to return errors to
// be handled by middleware functions.
type Handler func(w http.ResponseWriter, r *http.Request) error

// Converts middleware.Handler into http.Handler
func (h Handler) Build() http.Handler {
	// Create wrapper around h
	f := func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
	}
	return http.HandlerFunc(f)
}

func Logging(h Handler) Handler {
	return RequestLogger(ErrorLogger(h))
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

func ValidateRequest(w http.ResponseWriter, r *http.Request) error {
	if !(r.Method == http.MethodPost) {
		w.WriteHeader(405)
		w.Write([]byte("method not allowed"))
		return fmt.Errorf("method %v not allowed", r.Method)
	}

	header := w.Header().Get("Server-Action")

	if header == "" {
		w.WriteHeader(400)
		w.Write([]byte("bad request"))
		return fmt.Errorf("missing Server-Action header")
	}

	return nil
}
