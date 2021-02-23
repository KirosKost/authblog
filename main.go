package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"goproj/routes"
	"goproj/session"
	"html/template"

	"labix.org/v2/mgo"
)

func unescape(x string) interface{} {
return template.HTML(x)
}

func main() {
	fmt.Println("Listening on port :3000")

	mongoSession, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	db := mongoSession.DB("blog")

	mrtn := martini.Classic()

	unescapeFuncMap := template.FuncMap{"unescape": unescape}

	mrtn.Map(db)

	mrtn.Use(session.Middleware)

	mrtn.Use(render.Renderer(render.Options{
		Directory:  "templates",                         // Specify what path to load the templates from.
		Layout:     "layout",                            // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"},          // Specify extensions to load for templates.
		Funcs:      []template.FuncMap{unescapeFuncMap}, // Specify helper function maps for templates to access.
		Charset:    "UTF-8",                             // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,                                // Output human readable JSON
	}))

	staticOptions := martini.StaticOptions{Prefix: "static"}
	mrtn.Use(martini.Static("static", staticOptions))

	mrtn.Get("/", routes.IndexHandler)
	mrtn.Get("/login", routes.GetLoginHandler)
	mrtn.Post("/login", routes.PostLoginHandler)
	mrtn.Get("/create/", routes.CreateHandler)
	mrtn.Get("/edit/:id", routes.EditHandler)
	mrtn.Get("/delete/:id", routes.DeleteHandler)
	mrtn.Post("/savePost", routes.SavePostHandler)
	mrtn.Post("/gethtml", routes.GetHtmlHandler)

	mrtn.Run()
}