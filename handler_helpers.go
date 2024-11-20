package main

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// executeTemplate executes the given template to the buffer first,
// then copies the buffer to ResponseWriter. This way if an error happens
// during execution it returns immediately instead of half loading the page
func executeTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	// Write to buffer first
	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, templateName, data)
	if err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// If template executed successfully, go ahead and write to response
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Printf("Error writing response: %v", err)
		// At this point we can't really recover or send an error to the client
	}

}

// validateURL checks that the URL for artist_details.html is valid
// invalid formats include no artist id, no valid artist id, invalid section name
func validateURL(url string) (id, section string, err error) {

	// Remove the prefix to get everything after /artist_details/
	path := strings.TrimPrefix(url, "/artist_details/")

	// If nothing after /artist_details/ or ends with /
	if path == "" || path == "/" {
		err = errors.New("incomplete url path")
		return
	}

	// Split the remaining path into parts
	parts := strings.Split(path, "/")

	// Validate URL structure
	if len(parts) > 2 {
		// Too many parts in URL
		err = errors.New("url path too long")
		return
	}

	// Get the ID and validate it's a number
	id = parts[0]
	if _, err = strconv.Atoi(id); err != nil {
		// ID is not a valid number
		err = errors.New("artist ID not a valid number")
		return
	}

	// Get section if it exists and validate
	section = ""
	if len(parts) == 2 {
		section = parts[1]
		// Validate section is one of the allowed values
		if section != "locations" && section != "dates" && section != "concerts" {
			err = errors.New("invalid section")
			return
		}
	}

	return
}
