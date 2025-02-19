 
# Groupie Tracker 

## Description

Groupie Tracker is a website that receives an API with music artist information, manipulates the data received and displays it for the user. It includes filtering functionality and uses a geocoding API to display locations on a map. 

The server is written entirely in Go.

This is a project from Zone01.

 ## How to run
 Start the server with this command:
  `go run ./cmd/web`

  If you want to specify a different port to the default 8080, use the port flag:
  `go run ./cmd/web -port=:3000`

 This command redirects logging from standard output and standard error to files in the tmp folder and appends as needed:
  `go run ./cmd/web >>/tmp/info.log 2>>tmp/error.log`
