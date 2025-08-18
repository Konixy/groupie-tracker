package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"groupie-tracker/api"
)

// Handler for the /location route
func LocationsHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "ID d'artiste manquant", http.StatusBadRequest)
		return
	}

	artistID, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "ID d'artiste invalide", http.StatusBadRequest)
		return
	}

	relation, err := api.FetchLocations(artistID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des concerts", http.StatusInternalServerError)
		return
	}

	fmt.Println(len(relation.DatesLocations))

	// range on the datesLocations keys
	for location := range relation.DatesLocations {
		fmt.Println(api.GetCoordinates(location))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(relation)
}
