package blogs

import (
	"html/template"
	"time"

	"github.com/LucDeCaf/go-simple-blog/db"
)

type Blog struct {
	Id            int           `json:"id"`
	OwnerUsername string        `json:"owner_username"`
	Title         string        `json:"title"`
	Content       template.HTML `json:"content"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

func Get(id int) (Blog, error) {
	const query string = "SELECT id,owner_username,title,content,created_at,updated_at FROM blogs WHERE id=?;"

	var b Blog
	err := db.DB.QueryRow(query, id).Scan(
		&b.Id,
		&b.OwnerUsername,
		&b.Title,
		&b.Content,
		&b.CreatedAt,
		&b.UpdatedAt,
	)

	return b, err
}

func GetAll() ([]Blog, error) {
	const query string = "SELECT id,owner_username,title,content,created_at,updated_at FROM blogs;"

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}

	blogs := make([]Blog, 0)
	for {
		if !rows.Next() {
			break
		}

		var b Blog

		if err := rows.Scan(
			&b.Id,
			&b.OwnerUsername,
			&b.Title,
			&b.Content,
			&b.CreatedAt,
			&b.UpdatedAt,
		); err != nil {
			return blogs, err
		}

		blogs = append(blogs, b)
	}

	return blogs, rows.Err()
}

func Insert(blog Blog) (Blog, error) {
	const query string = "INSERT INTO blogs (owner_username,title,content) VALUES (?,?,?) RETURNING id,owner_username,title,content,created_at,updated_at;"

	var b Blog
	err := db.DB.QueryRow(query, blog.OwnerUsername, blog.Title, blog.Content).Scan(
		&b.Id,
		&b.OwnerUsername,
		&b.Title,
		&b.Content,
		&b.CreatedAt,
		&b.UpdatedAt,
	)

	return b, err
}

func Update(blog Blog) (Blog, error) {
	const query string = "INSERT INTO blogs (owner_username,title,content,updated_at) VALUES (?,?,?,?) WHERE id=? RETURNING id,owner_username,title,content,created_at,updated_at;"

	var b Blog
	err := db.DB.QueryRow(query, blog.OwnerUsername, blog.Title, blog.Content, time.Now(), blog.Id).Scan(
		&b.Id,
		&b.OwnerUsername,
		&b.Title,
		&b.Content,
		&b.CreatedAt,
		&b.UpdatedAt,
	)

	return b, err
}

func Delete(id int) (Blog, error) {
	const query string = "DELETE FROM blogs WHERE id=? RETURNING id,owner_username,title,content,created_at,updated_at;"

	var b Blog
	err := db.DB.QueryRow(query, id).Scan(
		&b.Id,
		&b.OwnerUsername,
		&b.Title,
		&b.Content,
		&b.CreatedAt,
		&b.UpdatedAt,
	)

	return b, err
}
