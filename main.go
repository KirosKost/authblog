package main

import (
	"fmt"
	"html/template"
	"net/http"
	"github.com/KirosKost/authblog/models"
)

func indexHandler(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil{
		fmt.Fprintf(w, err.Error())
	}

	tmpl.ExecuteTemplate(w, "index", nil)
}

func createHandler(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil{
		fmt.Fprintf(w, err.Error())
	}

	tmpl.ExecuteTemplate(w, "create", nil)
}


func savePostHandler(w http.ResponseWriter, r *http.Request){
	title := r.FormValue("title")
	content := r.FormValue("content")

	post := models.NewPost(title, content)
}

func main ()  {

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create/", createHandler)
	http.HandleFunc("/savePost", savePostHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8080", nil)

}