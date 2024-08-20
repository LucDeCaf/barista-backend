package models_test

import (
	"github.com/LucDeCaf/go-simple-blog/models"
	"github.com/LucDeCaf/go-simple-blog/models/author"
	"github.com/LucDeCaf/go-simple-blog/models/blog"
	"github.com/LucDeCaf/go-simple-blog/models/user"
)

func useModel[M models.Model]() {}

// Assert models implement Model
func _() {
	useModel[*author.Author]()
	useModel[*blog.Blog]()
	useModel[*user.User]()
}
