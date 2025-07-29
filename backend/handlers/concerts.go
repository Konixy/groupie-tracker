package handlers

import (
	"encoding/json"
	"net/http"

	"groupie-tracker/api"
)

// Handler for the /concerts route
func ConcertsHandler(w http.ResponseWriter, r *http.Request) {
	concerts, err := api.FetchAllConcerts()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des concerts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(concerts)
}
