package main

import (
	"log"
	"time"

	"groupie.tracker.filters/internal/groupietracker"
)

// Fetch data from the API
func (app *application) fetchArtistsData() error {
	artists, err := groupietracker.GetAllArtistsDetails()
	if err != nil {
		return err
	}

	// Lock the mutex to safely update the artists slice
	app.mu.Lock()
	app.artists = artists
	app.lastUpdated = time.Now()
	app.mu.Unlock()

	return nil
}

// Periodically refresh the artists data every hour
func (app *application) refreshArtistsData() {

	// Use a ticker to refresh the data every hour
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		<-ticker.C
		err := app.fetchArtistsData()
		if err != nil {
			log.Printf("Error fetching artists data: %v", err)
		} else {
			log.Println("Successfully refreshed artists data")
		}
	}
}
