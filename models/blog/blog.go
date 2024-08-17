package blog

import (
	"database/sql"
	"html/template"
	"time"
)

// TODO: This is error-prone to Scan from, possibly create func `ScanDefault(*sql.Row, *Blog) error`
const PubliclyReturned string = "id,author_id,title,content,created_at,updated_at"

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

func NewBlogTable(db *sql.DB) BlogTable {
	return BlogTable{db: db}
}

func (t *BlogTable) Get(id int) (*Blog, error) {
	var b Blog
	if err := t.db.QueryRow("SELECT "+PubliclyReturned+" FROM blogs WHERE id = ?;", id).Scan(
		&b.Id,
		&b.AuthorId,
		&b.Title,
		&b.Content,
		&b.CreatedAt,
		&b.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &b, nil
}

func (t *BlogTable) GetAll() ([]Blog, error) {
	rows, err := t.db.Query("SELECT " + PubliclyReturned + " FROM blogs;")
	if err != nil {
		return nil, err
	}

	blogs := make([]Blog, 0)
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

		blogs = append(blogs, b)
	}

	return blogs, rows.Err()
}

func (t *BlogTable) Insert(blog *Blog) (*Blog, error) {
	row := t.db.QueryRow("INSERT INTO blogs (author_id,title,content) VALUES (?,?,?) RETURNING "+PubliclyReturned+";", blog.AuthorId, blog.Title, blog.Content)

	var b Blog
	if err := row.Scan(&b.Id, &b.AuthorId, &b.Title, &b.Content, &b.CreatedAt, &b.UpdatedAt); err != nil {
		return nil, err
	}

	return &b, nil
}

func (t *BlogTable) Delete(blogId int) (*Blog, error) {
	row := t.db.QueryRow("DELETE FROM blogs WHERE id = ? RETURNING " + PubliclyReturned + ";", blogId)

	var b Blog
	if err := row.Scan(&b.Id, &b.AuthorId, &b.Title, &b.Content, &b.CreatedAt, &b.UpdatedAt); err != nil {
		return nil, err
	}
	return &b, nil
}
