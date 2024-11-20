package groupietracker

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// struct for the homepage, fetching only the data needed
type Artist struct {
	ID    int    `json:"id"`
	Image string `json:"image"`
	Name  string `json:"name"`
}

// struct for the details page, containing all information about an artist
// locations, dates and relations are ignored because they are filled by different api calls
type ArtistDetails struct {
	ID           int                 `json:"id"`
	Image        string              `json:"image"`
	Name         string              `json:"name"`
	Members      []string            `json:"members"`
	CreationDate int                 `json:"creationDate"`
	FirstAlbum   string              `json:"firstAlbum"`
	Locations    []string            `json:"-"`
	Dates        []string            `json:"-"`
	Relations    map[string][]string `json:"-"`
}

// GetAllArtists returns a slice of all artists read from the API. The Artist struct includes ID, name and image.
// It also returns an error if the process of getting the data from the API fails at any step, or nil if it succeeds.
func GetAllArtists() ([]Artist, error) {
	resp, err := http.Get(baseURL + endpointArtists)
	if err != nil {
		return nil, fmt.Errorf("could not fetch artists: %v", err)
	}

	defer resp.Body.Close()

	// Read the response body and unmarshal into a slice of Artist
	var artists []Artist
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		return nil, fmt.Errorf("could not parse artists: %v", err)
	}

	return artists, nil
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
