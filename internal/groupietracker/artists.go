package groupietracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// struct containing all information about an artist
// locations, dates and relations are initially ignored
// because they are filled by different api calls
type ArtistDetails struct {
	ID           int                 `json:"id"`
	Image        string              `json:"image"`
	Name         string              `json:"name"`
	Members      []string            `json:"members"`
	CreationDate int                 `json:"creationDate"`
	FirstAlbum   string              `json:"firstAlbum"`
	Locations    []*Location         `json:"-"`
	Dates        []string            `json:"-"`
	Relations    map[string][]string `json:"-"`
}

type Location struct {
	Country     string
	City        string
	Coordinates *GeocodeResponse
}

func GetAllArtistsDetails() ([]*ArtistDetails, error) {
	resp, err := http.Get(baseURL + endpointArtists)
	if err != nil {
		return nil, fmt.Errorf("could not fetch artists: %v", err)
	}
	// Read the response body and unmarshal into a slice of Artist
	var allArtistsDetails []*ArtistDetails
	if err := json.NewDecoder(resp.Body).Decode(&allArtistsDetails); err != nil {
		return nil, fmt.Errorf("could not parse artists: %v", err)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 3)

	for _, artist := range allArtistsDetails {
		wg.Add(3)

		go func(artist *ArtistDetails) {
			defer wg.Done()
			locations, err := GetLocations(strconv.Itoa(artist.ID))
			if err != nil {
				errCh <- fmt.Errorf("could not fetch locations for artist %d: %v", artist.ID, err)
				return
			}
			artist.Locations = locations
		}(artist)

		go func(artist *ArtistDetails) {
			defer wg.Done()
			dates, err := GetDates(strconv.Itoa(artist.ID))
			if err != nil {
				errCh <- fmt.Errorf("could not fetch dates for artist %d: %v", artist.ID, err)
				return
			}
			artist.Dates = dates
		}(artist)

		go func(artist *ArtistDetails) {
			defer wg.Done()
			relations, err := GetRelation(strconv.Itoa(artist.ID))
			if err != nil {
				errCh <- fmt.Errorf("could not fetch relations for artist %d: %v", artist.ID, err)
				return
			}
			artist.Relations = relations
		}(artist)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Close the error channel and check for any errors
	close(errCh)
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}
	return allArtistsDetails, nil
}

// GetArtistDetails returns all info about an artist of specific ID in a struct
// It gets data from 4 different API calls, and reformats the results before returning
// It returns an error if the process fails at any step, or nil if it succeeds.
func GetArtistDetails(id string) (*ArtistDetails, error) {
	resp, err := http.Get(baseURL + endpointArtists + "/" + id)
	if err != nil {
		return nil, fmt.Errorf("could not fetch artists: %v", err)
	}

	defer resp.Body.Close()

	// Read the response body and unmarshal into a slice of ArtistDetails
	var artistDetails ArtistDetails
	if err := json.NewDecoder(resp.Body).Decode(&artistDetails); err != nil {
		return nil, fmt.Errorf("could not parse artists: %v", err)
	}

	// if called with an id that doesn't exist, the API returns id=0
	if artistDetails.ID == 0 {
		return nil, fmt.Errorf("not found")
	}

	// API call to get locations for the specific artist ID
	artistDetails.Locations, err = GetLocations(id)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	// API call to get dates for the specific artist ID
	artistDetails.Dates, err = GetDates(id)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	// API call to get relation data for the specific artist ID
	artistDetails.Relations, err = GetRelation(id)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return &artistDetails, nil
}
