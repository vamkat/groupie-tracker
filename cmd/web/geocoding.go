package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"groupie.tracker.filters/internal/groupietracker"
)

func (app *application) getCoordinates(index int) {

}

func (app *application) getLocationsHandler(w http.ResponseWriter, r *http.Request) {
	// Get artist ID from query parameter
	artistID := r.URL.Query().Get("id")
	if artistID == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(artistID)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var artistIndex int
	for i, artist := range app.artists {
		if artist.ID == id {
			artistIndex = i
			break
		}
	}

	for i, location := range app.artists[artistIndex].Locations {

		var err error
		if app.artists[artistIndex].Locations[i].Coordinates == nil { //do not fetch again if already fetched
			app.artists[artistIndex].Locations[i].Coordinates, err = groupietracker.Geocode(location.Country, location.City)
			if err != nil {
				//just log the error and carry on
				app.errorLog.Print(err)
			}
		}
	}

	// Return locations as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app.artists[artistIndex].Locations)
}
