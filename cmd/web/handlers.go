package main

import (
	"net/http"
	"time"

	"groupie.tracker.filters/internal/groupietracker"
)

// handlers are application methods in order to use custom logging
func (app *application) artistsPage(w http.ResponseWriter, r *http.Request) {
	// return 404 if invalid URL
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	// Lock the mutex for reading the artists data
	app.mu.RLock()
	artists := app.artists
	lastUpdated := app.lastUpdated // Get the last updated time
	app.mu.RUnlock()

	// If there are no artists (for any reason), return an error
	if len(artists) == 0 {
		app.serverError(w, nil)
		return
	}

	data := struct {
		Artists     []*groupietracker.ArtistDetails
		LastUpdated time.Time
		FilterData  *FilterValues
	}{
		Artists:     artists,
		LastUpdated: lastUpdated,
		FilterData:  app.filterData,
	}

	app.executeTemplate(w, "artists.html", data)
}

func (app *application) artistDetailsPage(w http.ResponseWriter, r *http.Request) {
	// if invalid url return 404
	id, section, err := validateURL(r.URL.Path)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Fetch the artist details from the API
	artistDetails, err := groupietracker.GetArtistDetails(id)
	if err != nil {
		if err.Error() == "not found" {
			app.notFound(w)
			return
		} else {
			app.serverError(w, err)
			return
		}
	}

	// the struct sent to the template will now include both ArtistDetails and section
	data := struct {
		Artist      *groupietracker.ArtistDetails
		Section     string
		LastUpdated time.Time
	}{
		Artist:      artistDetails,
		Section:     section,
		LastUpdated: app.lastUpdated,
	}

	app.executeTemplate(w, "artist_details.html", data)
}
