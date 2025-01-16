package main

import "net/http"

// creates and returns a new serve mux containing all routes
func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static"))

	mux.Handle("/static/", http.StripPrefix("/static", neuter(fileServer)))

	mux.HandleFunc("/", app.artistsPage)
	mux.HandleFunc("/artist_details/", app.artistDetailsPage)
	mux.HandleFunc("/api/artist-locations", app.getLocationsHandler)

	return mux
}
