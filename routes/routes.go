package routes

import (
	"net/http"

	"github.com/LucDeCaf/go-simple-blog/auth"
	e "github.com/LucDeCaf/go-simple-blog/errors"
	"github.com/LucDeCaf/go-simple-blog/models/users"
)

var V1 *http.ServeMux

func ExtractUser(r *http.Request) (users.User, *e.HttpError) {
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
