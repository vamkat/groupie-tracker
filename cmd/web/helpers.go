package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

// serverError writes an error message and stack trace to the errorlog,
// then sends a generic 500 response to the user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Print(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError spends a specific status code and description to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// wrapper for clientError that sends a 404 notFound response
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// executeTemplate executes the given template to the buffer first,
// then copies the buffer to ResponseWriter. This way if an error happens
// during execution it returns immediately instead of half loading the page
func (app *application) executeTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	// Write to buffer first
	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, templateName, data)
	if err != nil {
		app.serverError(w, err)
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
func validateURL(url string) (id int, section string, err error) {

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
	if id, err = strconv.Atoi(parts[0]); err != nil {
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

// this blocks directory listing for the static folder, but doesn't block access
// to individual files if the user knows the exact path
func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
