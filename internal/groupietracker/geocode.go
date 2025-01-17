package groupietracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type GeocodeResponse struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func Geocode(country, city string) (*GeocodeResponse, error) {

	//construct query parameters
	params := url.Values{}
	if countrycode, ok := countryCodes[country]; ok {
		params.Add("q", city)
		params.Add("countrycodes", countrycode)
	} else {
		params.Add("q", country+"+"+city)
	}

	params.Add("format", "json")

	//make the GET request to Nominatim API
	apiURL := fmt.Sprintf("%s?%s", nominatimAPI, params.Encode())
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//Parse the JSON response
	var results []GeocodeResponse
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, err
	}

	if len(results) > 0 {
		//return the first result
		return &results[0], nil
	}
	return nil, fmt.Errorf("no results for %s, %s", city, country)
}
