package models

type Post struct {
	title, content string
}

func NewPost(title, content string) *Post {
	return &Post{title, content}
}