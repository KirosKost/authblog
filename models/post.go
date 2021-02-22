package models

type Post struct {
	Id, Title, Content string
}

func NewPost(id, title, content string) *Post {
	return &Post{id, title, content}
}