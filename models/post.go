package models

type Post struct {
	id, title, content string
}

func NewPost(title, content string) *Post {
	return &Post{id, title, content}
}