package main

import (
	"html/template"
	"log"
	"net/http"
)

// package-level variable that will hold all parsed templates
var tmpl *template.Template

func main() {
	//serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/artist_details/", artistDetailsPage)
	http.HandleFunc("/", artistsPage)

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func init() {
	//Parse all templates when program starts
	tmpl = template.Must(template.New("").ParseGlob("templates/*.html"))
}
