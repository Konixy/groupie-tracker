package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"groupie-tracker/api"
)

// Handler for the /locations route
func LocationsHandler(w http.ResponseWriter, r *http.Request) {
	// Extraire l'ID de l'artiste de l'URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "ID d'artiste manquant", http.StatusBadRequest)
		return
	}

	artistID := pathParts[2]
	if artistID == "" {
		http.Error(w, "ID d'artiste manquant", http.StatusBadRequest)
		return
	}

	// Convertir l'ID en int
	artistIDInt, err := strconv.Atoi(artistID)
	if err != nil {
		http.Error(w, "ID d'artiste invalide", http.StatusBadRequest)
		return
	}

	// Récupérer les concerts de l'artiste
	concerts, err := api.FetchArtistConcerts(artistIDInt)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des concerts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(concerts)
}

// Handler pour récupérer tous les lieux disponibles
func AllLocationsHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer tous les artistes pour collecter leurs lieux
	artists, err := api.FetchArtists()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des artistes", http.StatusInternalServerError)
		return
	}

	// Collecter tous les lieux uniques
	allLocations := make(map[string][]string) // lieu -> [artistes]

	for _, artist := range artists {
		// Récupérer les concerts de cet artiste
		concerts, err := api.FetchArtistConcerts(artist.ID)
		if err != nil {
			continue // Ignorer les erreurs pour cet artiste
		}

		// Parcourir tous les concerts
		for _, locations := range concerts.Concerts {
			for _, location := range locations {
				// Nettoyer le nom du lieu
				cleanLocation := strings.TrimSpace(location)
				if cleanLocation != "" {
					// Ajouter l'artiste à ce lieu
					if _, exists := allLocations[cleanLocation]; !exists {
						allLocations[cleanLocation] = []string{}
					}
					
					// Éviter les doublons d'artistes pour le même lieu
					artistExists := false
					for _, existingArtist := range allLocations[cleanLocation] {
						if existingArtist == artist.Name {
							artistExists = true
							break
						}
					}
					
					if !artistExists {
						allLocations[cleanLocation] = append(allLocations[cleanLocation], artist.Name)
					}
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allLocations)
}
