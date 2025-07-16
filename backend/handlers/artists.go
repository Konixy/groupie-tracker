package handlers

import (
	"net/http"

	"groupie-tracker/api"
)

// Handler for the / route
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := api.FetchArtists()

	// A toi de faire le reste
}
