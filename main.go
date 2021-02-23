package main

import (
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"goproj/db/documents"
	"goproj/models"
	"html/template"

	"labix.org/v2/mgo"
	"net/http"
)

var postsCollection *mgo.Collection

func indexHandler(rnd render.Render){
	postDocuments := []documents.PostDocument{}
	postsCollection.Find(nil).All(&postDocuments)

	posts := []models.Post{}
	for _, doc := range postDocuments{
		post := models.Post{doc.Id, doc.Title, doc.ContentHtml, doc.ContentMarkDown}
		posts = append(posts, post)
	}
	rnd.HTML(200, "index", posts)
}

func createHandler(rnd render.Render){
	post := models.Post{}

	rnd.HTML(200, "create", post)
}

func editHandler(rnd render.Render, r *http.Request, params martini.Params) {
	id := params["id"]

	postDocument := documents.PostDocument{}
	err := postsCollection.FindId(id).One(&postDocument)
	if err != nil {
		rnd.Redirect("/")
		return
	}
	post := models.Post{postDocument.Id, postDocument.Title, postDocument.ContentHtml, postDocument.ContentMarkDown}

	rnd.HTML(200, "write", post)
}


func savePostHandler(rnd render.Render, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	contentMarkdown := r.FormValue("content")
	contentHtml := ConvertMarkdownToHtml(contentMarkdown)

	postDocument := documents.PostDocument{id, title, contentHtml, contentMarkdown}
	if id != "" {
		postsCollection.UpdateId(id, postDocument)
	} else {
		id = GenerateId()
		postDocument.Id = id
		postsCollection.Insert(postDocument)
	}

	rnd.Redirect("/")
}

func deleteHandler(rnd render.Render, r *http.Request, params martini.Params){
	id := params["id"]
	if id == ""{
		rnd.Redirect("/")
		return
	}

	postsCollection.RemoveId(id)

	rnd.Redirect("/")
}

func getHtmlHandler(rnd render.Render, r *http.Request){
	md := r.FormValue("md")
	html := ConvertMarkdownToHtml(md)

	rnd.JSON(200, map[string]interface{}{"html": html})
}

func unescape(x string) interface{} {
	return template.HTML(x)
}

func main ()  {

	session, err := mgo.Dial("localhost") //27017
	if err != nil{
		panic(err)
	}
	postsCollection = session.DB("blog").C("posts")
	mrtn := martini.Classic()

	unescapeFuncMap := template.FuncMap{"unescape": unescape}

	mrtn.Use(render.Renderer(render.Options{
		Directory:       "templates",                         // Specify what path to load the templates from.
		Layout:          "layout",                            // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions:      []string{".tmpl", ".html"},          // Specify extensions to load for templates.
		Funcs:           []template.FuncMap{unescapeFuncMap}, // Specify helper function maps for templates to access.
		Charset:         "UTF-8",                             // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON:      true,                                // Output human readable JSON

	}))

	staticOptions := martini.StaticOptions{Prefix: "static"}
	mrtn.Use(martini.Static("static", staticOptions))

	mrtn.Get("/", indexHandler)
	mrtn.Get("/create/", createHandler)
	mrtn.Get("/edit/:id", editHandler)
	mrtn.Get("/delete/:id", deleteHandler)
	mrtn.Post("/savePost", savePostHandler)
	mrtn.Post("/gethtml", getHtmlHandler)

	mrtn.Run()

}