package main

import (
	"errors"
	"net/http"
	"time"

	"groupie.tracker.filters/internal/groupietracker"
)

// handlers are application methods in order to use custom logging
func (app *application) artistsPage(w http.ResponseWriter, r *http.Request) {
	defer func() {
		// Catch any panic that might occur
		if err := recover(); err != nil {
			// Log the panic for debugging
			app.infoLog.Printf("Recovered from panic in artistsPage: %v", err)

			// Respond with a 500 Internal Server Error
			app.serverError(w, errors.New("internal server error"))
		}
	}()

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

	// If there are no artists return a server error
	if len(artists) == 0 {
		err := errors.New("no artists found")
		app.serverError(w, err)
		return
	}

	//read queries
	queryParams := r.URL.Query()
	if len(queryParams) > 0 {
		membersList, minAlbum, maxAlbum, minCreation, maxCreation, country, city, err := app.handleQuery(queryParams)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}
		artists, err = app.executeFilters(membersList, minAlbum, maxAlbum, minCreation, maxCreation, country, city)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	//send data to template
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
	defer func() {
		// Catch any panic that might occur
		if err := recover(); err != nil {
			// Log the panic for debugging
			app.infoLog.Printf("Recovered from panic in artistDetailsPage: %v", err)

			// Respond with a 500 Internal Server Error
			app.serverError(w, errors.New("internal server error"))
		}
	}()

	// if invalid url return 404
	id, section, err := validateURL(r.URL.Path)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var artistDetails *groupietracker.ArtistDetails
	for i, artist := range app.artists {
		if id == artist.ID {
			if section == "locations" && app.artists[i].Locations[0].Coordinates == nil { //do not fetch again if already fetched
				//get geocoding information for all locations

				app.getCoordinates(i)

			}
			artistDetails = artist
		}
	}
	if artistDetails == nil {
		app.notFound(w)
		return
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
