package middleware

import (
	"log"
	"net/http"
)

// A type identical to the signature of http.HandlerFunc,
// but with the benefit of being able to return errors to
// be handled by middleware functions.
//
// The 
type Handler func(w http.ResponseWriter, r *http.Request) error

// Converts `middleware.Handler` into `http.Handler`
func Build(h Handler) http.Handler {
	// Create wrapper around h
	f := func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 internal server error"))
		}
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
