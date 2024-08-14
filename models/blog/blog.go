package blog

import (
	"database/sql"
	"html/template"
	"time"
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

func NewBlogTable(db *sql.DB) BlogTable {
	return BlogTable{db: db}
}

func (bt BlogTable) Get(id int) (Blog, error) {
	var b Blog
	if err := bt.db.QueryRow("SELECT id, author_id, title, content, created_at, updated_at FROM blogs WHERE id = ?;", id).Scan(
		&b.Id,
		&b.AuthorId,
		&b.Title,
		&b.Content,
		&b.CreatedAt,
		&b.UpdatedAt,
	); err != nil {
		return Blog{}, err
	}
	return b, nil
}

func (at BlogTable) GetAll() ([]Blog, error) {
	// Select rows
	rows, err := at.db.Query("SELECT id, author_id, title, content, created_at, updated_at FROM blogs;")
	if err != nil {
		return nil, err
	}

	// Read blogs from rows into slice
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
