package blog

import (
	"html/template"
	"time"
)

type Blog struct {
	Id        int           `json:"id"`
	AuthorId  int           `json:"authorId"`
	Title     string        `json:"title"`
	Content   template.HTML `json:"content"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}
