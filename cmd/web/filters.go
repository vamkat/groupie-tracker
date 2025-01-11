package main

import (
	"sort"
	"strconv"
	"strings"
)

type FilterValues struct {
	MaxMembers  int
	MinAlbum    int
	MaxAlbum    int
	MinCreation int
	MaxCreation int
	Locations   map[string][]string //keys are countries, values are cities
}

func (app *application) getFilterValues() {

	values := &FilterValues{
		MinAlbum:    app.artists[0].CreationDate, //just a random year to initialize the variable
		MinCreation: app.artists[0].CreationDate,
		Locations:   make(map[string][]string),
	}

	for _, artist := range app.artists {
		//get maximum band members
		if len(artist.Members) > values.MaxMembers {
			values.MaxMembers = len(artist.Members)
		}
		//get minimum and maximum creation years
		if artist.CreationDate < values.MinCreation {
			values.MinCreation = artist.CreationDate
		} else if artist.CreationDate > values.MaxCreation {
			values.MaxCreation = artist.CreationDate
		}
		//get minimum and maximum album years
		sepDate := strings.Split(artist.FirstAlbum, "-")
		year, err := strconv.Atoi(sepDate[2])
		if err != nil {
			app.errorLog.Println(err.Error())
		} else {
			if year < values.MinAlbum {
				values.MinAlbum = year
			} else if year > values.MaxAlbum {
				values.MaxAlbum = year
			}
		}
		for _, location := range artist.Locations {
			// map all countries and their cities
			if _, exists := values.Locations[location.Country]; !exists {
				values.Locations[location.Country] = []string{} // Initialize the slice if it doesn't exist
			}
			alreadyExists := false
			for _, city := range values.Locations[location.Country] {
				if city == location.City {
					alreadyExists = true
					break
				}
			}
			if !alreadyExists {
				values.Locations[location.Country] = append(values.Locations[location.Country], location.City)
				sort.Strings(values.Locations[location.Country])
			}
		}
	}
	app.filterData = values
}
