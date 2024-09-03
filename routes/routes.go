package routes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/LucDeCaf/go-simple-blog/auth"
	e "github.com/LucDeCaf/go-simple-blog/errors"
	"github.com/LucDeCaf/go-simple-blog/models/authors"
	"github.com/LucDeCaf/go-simple-blog/models/blogs"
	"github.com/LucDeCaf/go-simple-blog/models/users"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func httpRespond(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	w.Write([]byte(fmt.Sprintf("%v %v\n", code, message)))
}

func extractUser(r *http.Request) (users.User, *e.HttpError) {
	token, err := auth.ExtractJWT(r)
	if err != nil {
		return users.User{}, e.NewHttpError(401, "unauthenticated")
	}

	claims, err := auth.ExtractClaims(token)
	if err != nil {
		return users.User{}, e.NewHttpError(401, "unauthenticated")
	}

	user, err := users.Get(claims.Username)
	if err != nil {
		return user, e.NewHttpError(500, "internal server error")
	}

	return user, nil
}

func AuthorsHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		allAuthors, err := authors.GetAll()
		if err != nil {
			httpRespond(w, 500, "internal server error")
			return err
		}

		en := json.NewEncoder(w)
		if err := en.Encode(allAuthors); err != nil {
			httpRespond(w, 500, "internal server error")
			return err
		}

	case http.MethodPost:
		user, httpErr := extractUser(r)
		if httpErr != nil {
			httpRespond(w, httpErr.Code, httpErr.Message)
			return httpErr
		}

		if user.Role != users.RoleAdmin {
			httpRespond(w, 403, "unauthorized")
			return e.NewHttpError(403, "unauthorized")
		}

		de := json.NewDecoder(r.Body)
		de.DisallowUnknownFields()

		var author authors.Author

		if err := de.Decode(&author); err != nil {
			httpRespond(w, 400, "bad request")
			return err
		}

		author, err := authors.Insert(author)
		if err != nil {
			httpRespond(w, 500, "internal server error")
			return err
		}

		w.WriteHeader(201)
		if err := json.NewEncoder(w).Encode(author); err != nil {
			return err
		}
	default:
		httpRespond(w, 405, "method not allowed")
		return fmt.Errorf("method '%v' not allowed", r.Method)
	}

	return nil
}

func BlogsHandler(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		allBlogs, err := blogs.GetAll()
		if err != nil {
			httpRespond(w, 500, "internal server error")
			return err
		}

		en := json.NewEncoder(w)
		if err := en.Encode(allBlogs); err != nil {
			httpRespond(w, 500, "internal server error")
			return err
		}

	case http.MethodPost:
		user, httpErr := extractUser(r)
		if httpErr != nil {
			httpRespond(w, httpErr.Code, httpErr.Message)
			return httpErr
		}

		if user.Role != users.RoleAdmin {
			httpRespond(w, 403, "unauthorized")
			return e.NewHttpError(403, fmt.Sprintf("unauthorized role '%v'", user.Role))
		}

		de := json.NewDecoder(r.Body)
		de.DisallowUnknownFields()

		var blog blogs.Blog
		if err := de.Decode(&blog); err != nil {
			httpRespond(w, 400, "bad request")
			return err
		}

		blog, err := blogs.Insert(blog)
		if err != nil {
			httpRespond(w, 500, "internal server error")
			return err
		}

		w.WriteHeader(201)
		if err := json.NewEncoder(w).Encode(blog); err != nil {
			return err
		}
	default:
		httpRespond(w, 405, "method not allowed")
		return fmt.Errorf("method '%v' not allowed", r.Method)
	}

	return nil
}

func AuthorsIdHandler(w http.ResponseWriter, r *http.Request) error {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpRespond(w, 400, "bad request")
		return err
	}

	switch r.Method {
	case http.MethodGet:
		author, err := authors.Get(id)
		if err != nil {
			httpRespond(w, 404, "not found")
			return err
		}

		if err = json.NewEncoder(w).Encode(author); err != nil {
			return err
		}

	case http.MethodDelete:
		user, httpErr := extractUser(r)
		if httpErr != nil {
			httpRespond(w, httpErr.Code, httpErr.Message)
			return httpErr
		}

		if user.Role != users.RoleAdmin {
			httpRespond(w, 403, "unauthorized")
			return e.NewHttpError(403, fmt.Sprintf("unauthorized role '%v'", user.Role))
		}

		author, err := authors.Delete(id)
		if err != nil {
			httpRespond(w, 500, "internal server error")
			return err
		}

		w.WriteHeader(200)
		if err := json.NewEncoder(w).Encode(author); err != nil {
			return err
		}

	default:
		httpRespond(w, 405, "method not allowed")
		return fmt.Errorf("method '%v' not allowed", r.Method)
	}

	return nil
}

func BlogsIdHandler(w http.ResponseWriter, r *http.Request) error {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httpRespond(w, 400, "bad request")
		return err
	}

	switch r.Method {
	case http.MethodGet:
		blog, err := blogs.Get(id)
		if err != nil {
			httpRespond(w, 404, "not found")
			return err
		}

		if err = json.NewEncoder(w).Encode(blog); err != nil {
			return err
		}

	case http.MethodDelete:
		user, httpErr := extractUser(r)
		if httpErr != nil {
			httpRespond(w, httpErr.Code, httpErr.Message)
			return httpErr
		}

		if user.Role != users.RoleAdmin {
			httpRespond(w, 403, "unauthorized")
			return e.NewHttpError(403, fmt.Sprintf("unauthorized role '%v'", user.Role))
		}

		blog, err := blogs.Delete(id)
		if err != nil {
			httpRespond(w, 500, "internal server error")
			return err
		}

		w.WriteHeader(200)
		if err := json.NewEncoder(w).Encode(blog); err != nil {
			return err
		}

	default:
		httpRespond(w, 405, "method not allowed")
		return fmt.Errorf("method '%v' not allowed", r.Method)
	}

	return nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		httpRespond(w, 405, "method not allowed")
		return fmt.Errorf("invalid method %v", r.Method)
	}

	var lr loginRequest

	de := json.NewDecoder(r.Body)
	if err := de.Decode(&lr); err != nil {
		httpRespond(w, 400, "bad request")
		return err
	}

	user, err := users.Get(lr.Username)
	if err != nil {
		httpRespond(w, 500, "internal server error")
		return err
	}

	if !auth.VerifyPassword(lr.Password, user.PasswordHashWithSalt) {
		httpRespond(w, 401, "unauthenticated")
		return fmt.Errorf("wrong password inputted")
	}

	token, err := auth.NewJWT(lr.Username)
	if err != nil {
		httpRespond(w, 500, "internal server error")
		return err
	}

	w.Write([]byte(fmt.Sprintf(`{"token":"%v"}`, token)))

	return nil
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		httpRespond(w, 405, "method not allowed")
		return fmt.Errorf("method '%v' not allowed", r.Method)
	}

	var lr loginRequest
	if err := json.NewDecoder(r.Body).Decode(&lr); err != nil {
		httpRespond(w, 400, "bad request")
		return err
	}

	_, err := users.Get(lr.Username)

	// Successful return means user already exists
	if err == nil {
		httpRespond(w, 400, "user already exists")
		return fmt.Errorf("username already exists")
	}

	// ErrNoRows means record doesn't exist and is required
	// to avoid overwriting users
	if !errors.Is(err, sql.ErrNoRows) {
		httpRespond(w, 500, "internal server error")
		return err
	}

	pwHash, err := auth.HashPassword(lr.Password)
	if err != nil {
		httpRespond(w, 500, "internal server error")
		return err
	}

	user, err := users.Insert(users.User{
		Username:             lr.Username,
		PasswordHashWithSalt: pwHash,
		Role:                 users.RoleUser,
	})
	if err != nil {
		httpRespond(w, 500, "internal server error")
		return err
	}

	w.WriteHeader(201)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		return err
	}

	return nil
}
