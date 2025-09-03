package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"groupie-tracker/api"
	"groupie-tracker/config"
)

// Handler for the /artists route
func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := api.FetchArtists()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des artistes", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", config.GetFrontendURL())
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(artists)
}

// Handler for /artists/{id} route
func ArtistConcertsHandler(w http.ResponseWriter, r *http.Request) {
	// Extraire l'ID de l'artiste de l'URL
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

	// Récupérer les concerts de l'artiste
	concerts, err := api.FetchArtistConcerts(artistID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des concerts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", config.GetFrontendURL())
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(concerts)
}
