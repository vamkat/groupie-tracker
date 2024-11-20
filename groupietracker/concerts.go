package groupietracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// GetLocations gets the locations for a specific artist ID from the API.
// After reformatting them, it returns them as a slice of string.
// It returns an error if the process fails at any point, or nil if it succeeds.
func GetLocations(id string) ([]string, error) {

	//local struct to unmarshal into
	type Locations struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
		Dates     string   `json:"dates"`
	}

	resp, err := http.Get(baseURL + endpointLocations + "/" + id)
	if err != nil {
		return nil, fmt.Errorf("could not fetch locations: %v", err)
	}

	defer resp.Body.Close()

	// unmarshal the JSON
	var locations Locations
	if err := json.NewDecoder(resp.Body).Decode(&locations); err != nil {
		return nil, fmt.Errorf("could not parse locations: %v", err)
	}

	// reformat the received locations
	for i := range locations.Locations {
		locations.Locations[i] = formatLocationString(locations.Locations[i])
	}

	return locations.Locations, nil
}

// GetDates gets the dates for a specific artist ID from the API and returns them.
// It also returns an error if the process fails at any step, or nil if it succeeds.
func GetDates(id string) ([]string, error) {
	//local struct to unmarshal into
	type Dates struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	}

	resp, err := http.Get(baseURL + endpointDates + "/" + id)
	if err != nil {
		return nil, fmt.Errorf("could not fetch dates: %v", err)
	}

	defer resp.Body.Close()

	// unmarshal the JSON
	var dates Dates
	if err := json.NewDecoder(resp.Body).Decode(&dates); err != nil {
		return nil, fmt.Errorf("could not parse dates: %v", err)
	}

	return dates.Dates, nil
}

// GetLocations gets the relation data for a specific artist ID from the API.
// After reformatting the locations, it returns a map with the locations as key and the dates as value.
// It returns an error if the process fails at any point, or nil if it succeeds.
func GetRelation(id string) (map[string][]string, error) {

	type Relation struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	}

	resp, err := http.Get(baseURL + endpointRelation + "/" + id)
	if err != nil {
		return nil, fmt.Errorf("could not fetch relation: %v", err)
	}

	defer resp.Body.Close()

	// unmarshal the JSON
	var relation Relation
	if err := json.NewDecoder(resp.Body).Decode(&relation); err != nil {
		return nil, fmt.Errorf("could not parse relation data: %v", err)
	}

	// change the locations into a more readable format
	relation.DatesLocations = formatRelations(relation.DatesLocations)

	return relation.DatesLocations, nil
}

// formatRelations copies the map sent to it with the keys(locations) reformatted
func formatRelations(relations map[string][]string) map[string][]string {
	formattedMap := make(map[string][]string)
	for key, dates := range relations {
		// Transform the key
		newKey := formatLocationString(key)
		// Insert the transformed key-value pair into the new map
		formattedMap[newKey] = dates
	}
	return formattedMap
}

// formatLocationString receives a string(location) and changes it to a more readable format.
// It replaces underscores with spaces and dashes with commas.
// It capitalizes all words as they are place names.
// In the specific cases of USA and UK it turns the whole word to uppercase.
func formatLocationString(s string) string {
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", ", ")

	// Handle edge cases for "Uk" and "Usa"
	words := strings.Fields(s)
	for j, word := range words {
		if word == "uk" {
			words[j] = "UK"
		} else if word == "usa" {
			words[j] = "USA"
		} else {
			words[j] = strings.Title(words[j])
		}
	}

	// Join the words back into a single string
	s = strings.Join(words, " ")

	return s
}
