package v1

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/LucDeCaf/go-simple-blog/models/blogs"
	"github.com/LucDeCaf/go-simple-blog/routes"
	"github.com/LucDeCaf/go-simple-blog/sanitizer"
)

func BlogsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		allBlogs, err := blogs.GetAll()
		if err != nil {
			log.Println("err getting all blogs:", err.Error())
			http.Error(w, "internal server error", 500)
			return
		}

		en := json.NewEncoder(w)
		if err := en.Encode(allBlogs); err != nil {
			log.Println("err encoding body to blog:", err.Error())
			http.Error(w, "internal server error", 500)
			return
		}

	case http.MethodPost:
		user, httpErr := routes.ExtractUser(r)
		if httpErr != nil {
			log.Println("err extracting user from headers:", httpErr.Error())
			http.Error(w, httpErr.Message, httpErr.Code)
			return
		}

		if user.Role != "admin" {
			log.Println("user is not admin")
			http.Error(w, "unauthorized", 403)
			return
		}

		de := json.NewDecoder(r.Body)
		de.DisallowUnknownFields()

		var blog blogs.Blog
		if err := de.Decode(&blog); err != nil {
			log.Println("err decoding blog:", err.Error())
			http.Error(w, "bad request", 400)
			return
		}

		// Sanitize HTML content using UGC (User Generated Content) sanitizer
		blog.Content = template.HTML(sanitizer.Sanitize(string(blog.Content)))

		blog, err := blogs.Insert(blog)
		if err != nil {
			log.Println("err inserting blog:", err.Error())
			http.Error(w, "internal server error", 500)
			return
		}

		w.WriteHeader(201)
		if err := json.NewEncoder(w).Encode(blog); err != nil {
			log.Println("err encoding blog:", err.Error())
		}
	default:
		http.Error(w, "method not allowed", 405)
		return
	}
}

func BlogsIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("err parsing id:", err.Error())
		http.Error(w, "bad request", 400)
		return
	}

	switch r.Method {
	case http.MethodGet:
		blog, err := blogs.Get(id)
		if err != nil {
			log.Println("err getting blog:", err.Error())
			http.Error(w, "not found", 404)
			return
		}

		if err = json.NewEncoder(w).Encode(blog); err != nil {
			log.Println("err encoding blog:", err.Error())
			http.Error(w, "internal server error", 500)
			return
		}

	case http.MethodDelete:
		user, httpErr := routes.ExtractUser(r)
		if httpErr != nil {
			log.Println("err extracting user:", httpErr.Error())
			http.Error(w, httpErr.Message, httpErr.Code)
			return
		}

		if user.Role != "admin" {
			log.Println("user role unauthorized:", user.Role)
			http.Error(w, "unauthorized", 403)
			return
		}

		blog, err := blogs.Delete(id)
		if err != nil {
			log.Println("err deleting blog:", err.Error())
			http.Error(w, "internal server error", 500)
			return
		}

		w.WriteHeader(200)
		if err := json.NewEncoder(w).Encode(blog); err != nil {
			log.Println("err encoding blog:", err.Error())
			return
		}

	default:
		http.Error(w, "method not allowed", 405)
		return
	}
}
