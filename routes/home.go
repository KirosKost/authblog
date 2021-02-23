package routes

import (
	"fmt"
	"github.com/martini-contrib/render"
	"goproj/db/documents"
	"goproj/models"
	"goproj/session"
	"labix.org/v2/mgo"
)


func IndexHandler(rnd render.Render, s *session.Session, db *mgo.Database) {
	fmt.Println(s.Username)

	postDocuments := []documents.PostDocument{}
	postsCollection := db.C("posts")
	postsCollection.Find(nil).All(&postDocuments)

	posts := []models.Post{}
	for _, doc := range postDocuments {
		post := models.Post{doc.Id, doc.Title, doc.ContentHtml, doc.ContentMarkDown}
		posts = append(posts, post)
	}

	rnd.HTML(200, "index", posts)
}