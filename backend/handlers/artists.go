package handlers

import (
	"encoding/json"
	"net/http"

	"groupie-tracker/api"
)

// Handler for the / route
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := api.FetchArtists()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des artistes", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(artists)
}
