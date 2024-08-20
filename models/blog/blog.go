package blog

import (
	"database/sql"
	"fmt"
	"html/template"
	"time"

	"github.com/LucDeCaf/go-simple-blog/models"
)

type Blog struct {
	Id        int           `json:"id"`
	AuthorId  int           `json:"author_id"`
	Title     string        `json:"title"`
	Content   template.HTML `json:"content"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type BlogTable struct {
	db *sql.DB
}

func (*Blog) FieldNames() string {
	return "id,author_id,title,content,created_at,updated_at"
}

func FieldNames() string {
	return (*Blog).FieldNames(nil)
}

func (b *Blog) Fields() []any {
	return []any{
		&b.Id,
		&b.AuthorId,
		&b.Title,
		&b.Content,
		&b.CreatedAt,
		&b.UpdatedAt,
	}
}

func (b *Blog) Values() []any {
	return models.ValuesFromFields(b)
}

func NewBlogTable(db *sql.DB) BlogTable {
	return BlogTable{db: db}
}

func (t BlogTable) Get(id int) (*Blog, error) {
	query := fmt.Sprintf("SELECT %v FROM blogs WHERE id=?;", FieldNames())

	row := t.db.QueryRow(query, id)

	var b Blog
	if err := row.Scan(b.Fields()...); err != nil {
		return nil, err
	}

	return &b, nil
}

func (t BlogTable) GetAll() ([]*Blog, error) {
	query := fmt.Sprintf("SELECT %v FROM blogs;", FieldNames())

	rows, err := t.db.Query(query)
	if err != nil {
		return nil, err
	}

	blogs := make([]*Blog, 0)
	for {
		if !rows.Next() {
			break
		}

		var b Blog

		rows.Scan(
			&b.Id,
			&b.AuthorId,
			&b.Title,
			&b.Content,
			&b.CreatedAt,
			&b.UpdatedAt,
		)

		blogs = append(blogs, &b)
	}

	return blogs, rows.Err()
}

func (t BlogTable) Insert(blog *Blog) (*Blog, error) {
	query := fmt.Sprintf("INSERT INTO blogs (author_id,title,content) VALUES (?,?,?) RETURNING %v;", FieldNames())

	row := t.db.QueryRow(query, blog.AuthorId, blog.Title, blog.Content)

	var b Blog
	if err := row.Scan(b.Fields()...); err != nil {
		return nil, err
	}

	return &b, nil
}

func (t BlogTable) Delete(id int) (*Blog, error) {
	query := fmt.Sprintf("DELETE FROM blogs WHERE id=? RETURNING %v;", FieldNames())

	row := t.db.QueryRow(query, id)

	var b Blog
	if err := row.Scan(b.Fields()...); err != nil {
		return nil, err
	}

	return &b, nil
}
