package main

import (
	"errors"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"groupie.tracker.filters/internal/groupietracker"
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

func (app *application) handleQuery(queryParams url.Values) (memberList []int, minAlbum, maxAlbum, minCreation, maxCreation int, country, city string, err error) {
	//members
	members := queryParams["members"]
	if len(members) > 0 {
		for _, member := range members {
			// Convert string to int and add to the list
			var memberInt int
			memberInt, err = strconv.Atoi(member)
			//validate members
			if err != nil || memberInt < 1 || memberInt > app.filterData.MaxMembers {

				err = errors.New("invalid member number")
				return
			}
			memberList = append(memberList, memberInt)
		}
	}
	//minimum album date
	minAlbum, err = strconv.Atoi(queryParams.Get("minAlbum"))
	//validate minAlbum
	if err != nil || minAlbum < app.filterData.MinAlbum {
		err = errors.New("invalid first album date")
		return
	}
	//maximum album date
	maxAlbum, err = strconv.Atoi(queryParams.Get("maxAlbum"))
	//validate maxAlbum
	if err != nil || maxAlbum > app.filterData.MaxAlbum {
		err = errors.New("invalid first album date")
		return
	}
	//min creation date
	minCreation, err = strconv.Atoi(queryParams.Get("minCreation"))
	//validate minCreation
	if err != nil || minCreation < app.filterData.MinCreation {
		err = errors.New("invalid creation date")
		return
	}
	//max creation date
	maxCreation, err = strconv.Atoi(queryParams.Get("maxCreation"))
	//validate maxCreation
	if err != nil || maxCreation > app.filterData.MaxCreation {
		err = errors.New("invalid creation date")
		return
	}
	//country
	country = queryParams.Get("country")
	//validate country
	if country != "" {
		_, ok := app.filterData.Locations[country]
		if !ok {
			err = errors.New("invalid country")
			return
		}
	}

	//city
	city = queryParams.Get("city")
	if country != "" && city != "" {
		var cityFound bool
		for _, mapCity := range app.filterData.Locations[country] {
			if mapCity == city {
				cityFound = true
				break
			}
		}
		if !cityFound {
			err = errors.New("invalid city")
			return
		}
	}
	err = nil
	return
}

func (app *application) executeFilters(membersList []int, minAlbum, maxAlbum, minCreation, maxCreation int, country, city string) ([]*groupietracker.ArtistDetails, error) {
	// find all ids that much all the criteria
	var filteredArtists []*groupietracker.ArtistDetails
	for _, artist := range app.artists {
		//check for matchig members (automatic true if no members filter was used)
		var membersMatch bool
		if len(membersList) > 0 {
			for _, membersNum := range membersList {
				if membersNum == len(artist.Members) {
					membersMatch = true
					break
				}
			}
		} else {
			membersMatch = true
		}
		//check for first album date match
		var albumMatch bool
		firstAlbumYear, err := strconv.Atoi(strings.Split(artist.FirstAlbum, "-")[2])
		if err != nil {
			return nil, errors.New("error reading first album date")
		}
		if minAlbum <= firstAlbumYear && maxAlbum >= firstAlbumYear {
			albumMatch = true
		}
		// check for creation date match
		var creationMatch bool
		if minCreation <= artist.CreationDate && maxCreation >= artist.CreationDate {
			creationMatch = true
		}
		// check for location match (automatic match if no location set)
		var locationMatch bool
		if country == "" {
			locationMatch = true
		} else if city == "" {
			//match the country
			for _, location := range artist.Locations {
				if location.Country == country {
					locationMatch = true
					break
				}
			}
		} else {
			for _, location := range artist.Locations {
				if location.City == city {
					locationMatch = true
					break
				}
			}
		}
		if membersMatch && albumMatch && creationMatch && locationMatch {
			filteredArtists = append(filteredArtists, artist)
		}
	}
	return filteredArtists, nil
}
