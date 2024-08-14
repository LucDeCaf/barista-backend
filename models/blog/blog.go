package blog

import (
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

func NewBlog(id, authorId int, title string, content template.HTML) Blog {
	return Blog{
		Id:        id,
		AuthorId:  authorId,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
