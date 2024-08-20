package tables_test

import (
	"github.com/LucDeCaf/go-simple-blog/models"
	"github.com/LucDeCaf/go-simple-blog/models/author"
	"github.com/LucDeCaf/go-simple-blog/models/blog"
	"github.com/LucDeCaf/go-simple-blog/models/user"
	"github.com/LucDeCaf/go-simple-blog/tables"
)

func useTable[M models.Model, PK any](t tables.Table[M, PK]) {}

// Assert tables implement Table
func _() {
	useTable(author.AuthorTable{})
	useTable(blog.BlogTable{})
	useTable(user.UserTable{})
}
