package main

import "groupie.tracker.filters/internal/groupietracker"

func (app *application) getCoordinates(index int) {

	for i, location := range app.artists[index].Locations {

		var err error
		app.artists[index].Locations[i].Coordinates, err = groupietracker.Geocode(location.Country, location.City)
		if err != nil {
			//just log the error and carry on
			app.errorLog.Print(err)
		}

	}

}
