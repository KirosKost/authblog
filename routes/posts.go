package routes

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"goproj/db/documents"
	"goproj/models"
	"goproj/utils"
	"labix.org/v2/mgo"
	"net/http"
)

func CreateHandler(rnd render.Render) {
post := models.Post{}
rnd.HTML(200, "create", post)
}

func EditHandler(rnd render.Render, r *http.Request, params martini.Params, db *mgo.Database) {
	postsCollection := db.C("posts")

	id := params["id"]

	postDocument := documents.PostDocument{}
	err := postsCollection.FindId(id).One(&postDocument)
	if err != nil {
		rnd.Redirect("/")
		return
	}
	post := models.Post{postDocument.Id, postDocument.Title, postDocument.ContentHtml, postDocument.ContentMarkDown}

	rnd.HTML(200, "create", post)
}

func SavePostHandler(rnd render.Render, r *http.Request, db *mgo.Database) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	contentMarkdown := r.FormValue("content")
	contentHtml := utils.ConvertMarkdownToHtml(contentMarkdown)

	postDocument := documents.PostDocument{id, title, contentHtml, contentMarkdown}

	postsCollection := db.C("posts")
	if id != "" {
		postsCollection.UpdateId(id, postDocument)
	} else {
		id = utils.GenerateId()
		postDocument.Id = id
		postsCollection.Insert(postDocument)
	}

	rnd.Redirect("/")
}

func DeleteHandler(rnd render.Render, r *http.Request, params martini.Params, db *mgo.Database) {
	id := params["id"]
	if id == "" {
		rnd.Redirect("/")
		return
	}

	postsCollection := db.C("posts")
	postsCollection.RemoveId(id)

	rnd.Redirect("/")
}

func GetHtmlHandler(rnd render.Render, r *http.Request) {
	md := r.FormValue("md")
	html := utils.ConvertMarkdownToHtml(md)

	rnd.JSON(200, map[string]interface{}{"html": html})
}