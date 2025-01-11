package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"groupie.tracker.filters/internal/groupietracker"
)

type application struct {
	errorLog    *log.Logger
	infoLog     *log.Logger
	artists     []*groupietracker.ArtistDetails
	mu          sync.RWMutex
	lastUpdated time.Time
	filterData  *FilterValues
}

// package-level variable that will hold all parsed templates
var tmpl *template.Template

func main() {
	//get the port from an optional command line flag when starting the server
	port := flag.String("port", ":8080", "HTTP network address")
	flag.Parse()

	//create two logs for information and errors
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	//initialize a new application
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		artists:  []*groupietracker.ArtistDetails{},
	}

	// Call fetchArtistsData initially to load the data
	err := app.fetchArtistsData()
	if err != nil {
		errorLog.Printf("Error fetching artists data: %v", err)
	}

	// Start the background task to refresh the artists data
	go app.refreshArtistsData()

	// Get the values needed for the filters
	app.getFilterValues()

	//preParse()

	//initialize new http.Server struct so that it uses our custom errorLog
	srv := &http.Server{
		Addr:     *port,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	//start the server
	infoLog.Printf("Server started at http://localhost%s\n", *port)
	errorLog.Fatal(srv.ListenAndServe())
}

func init() {
	// Create the function map first
	funcMap := template.FuncMap{
		"iterate": func(start, end int) []int {
			var result []int
			for i := start; i <= end; i++ {
				result = append(result, i)
			}
			return result
		},
	}
	// Create new template with functions, then parse all files
	tmpl = template.Must(template.New("").Funcs(funcMap).ParseGlob("./ui/html/*.html"))
}
