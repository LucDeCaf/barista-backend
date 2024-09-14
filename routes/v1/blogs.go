package v1

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	e "github.com/LucDeCaf/go-simple-blog/errors"
	"github.com/LucDeCaf/go-simple-blog/models/blogs"
	"github.com/LucDeCaf/go-simple-blog/routes"
	"github.com/LucDeCaf/go-simple-blog/sanitizer"
)

func BlogsHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		allBlogs, err := blogs.GetAll()
		if err != nil {
			routes.HttpRespond(w, 500, "internal server error")
			return err
		}

		en := json.NewEncoder(w)
		if err := en.Encode(allBlogs); err != nil {
			routes.HttpRespond(w, 500, "internal server error")
			return err
		}

	case http.MethodPost:
		user, httpErr := routes.ExtractUser(r)
		if httpErr != nil {
			routes.HttpRespond(w, httpErr.Code, httpErr.Message)
			return httpErr
		}

		if user.Role != "admin" {
			routes.HttpRespond(w, 403, "unauthorized")
			return e.NewHttpError(403, fmt.Sprintf("unauthorized role '%v'", user.Role))
		}

		de := json.NewDecoder(r.Body)
		de.DisallowUnknownFields()

		var blog blogs.Blog
		if err := de.Decode(&blog); err != nil {
			routes.HttpRespond(w, 400, "bad request")
			return err
		}

		// Sanitize HTML content using UGC (User Generated Content) sanitizer
		blog.Content = template.HTML(sanitizer.Sanitize(string(blog.Content)))

		blog, err := blogs.Insert(blog)
		if err != nil {
			routes.HttpRespond(w, 500, "internal server error")
			return err
		}

		w.WriteHeader(201)
		if err := json.NewEncoder(w).Encode(blog); err != nil {
			return err
		}
	default:
		routes.HttpRespond(w, 405, "method not allowed")
		return fmt.Errorf("method '%v' not allowed", r.Method)
	}

	return nil
}

func BlogsIdHandler(w http.ResponseWriter, r *http.Request) error {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		routes.HttpRespond(w, 400, "bad request")
		return err
	}

	switch r.Method {
	case http.MethodGet:
		blog, err := blogs.Get(id)
		if err != nil {
			routes.HttpRespond(w, 404, "not found")
			return err
		}

		if err = json.NewEncoder(w).Encode(blog); err != nil {
			return err
		}

	case http.MethodDelete:
		user, httpErr := routes.ExtractUser(r)
		if httpErr != nil {
			routes.HttpRespond(w, httpErr.Code, httpErr.Message)
			return httpErr
		}

		if user.Role != "admin" {
			routes.HttpRespond(w, 403, "unauthorized")
			return e.NewHttpError(403, fmt.Sprintf("unauthorized role '%v'", user.Role))
		}

		blog, err := blogs.Delete(id)
		if err != nil {
			routes.HttpRespond(w, 500, "internal server error")
			return err
		}

		w.WriteHeader(200)
		if err := json.NewEncoder(w).Encode(blog); err != nil {
			return err
		}

	default:
		routes.HttpRespond(w, 405, "method not allowed")
		return fmt.Errorf("method '%v' not allowed", r.Method)
	}

	return nil
}
