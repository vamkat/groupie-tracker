package main

import (
	"log"
	"net/http"

	"groupie-tracker/groupietracker"
)

func artistsPage(w http.ResponseWriter, r *http.Request) {
	// return 400 if invalid URL
	if r.URL.Path != "/" && r.URL.Path != "/artists" {
		http.NotFound(w, r)
		return
	}

	// Get all artists from the API
	artists, err := groupietracker.GetAllArtists()
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	executeTemplate(w, "artists.html", artists)
}

func artistDetailsPage(w http.ResponseWriter, r *http.Request) {
	// if invalid url return 404
	id, section, err := validateURL(r.URL.Path)
	if err != nil {
		log.Print(err)

		http.Error(w, "Bad Request", http.StatusBadRequest)

		return
	}

	// Fetch the artist details from the API
	artistDetails, err := groupietracker.GetArtistDetails(id)
	if err != nil {
		if err.Error() == "not found" {
			log.Print(err)
			http.NotFound(w, r)
			return
		} else {
			log.Print(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// the struct sent to the template will now include both ArtistDetails and section
	data := struct {
		Artist  *groupietracker.ArtistDetails
		Section string
	}{
		Artist:  artistDetails,
		Section: section,
	}

	executeTemplate(w, "artist_details.html", data)
}
