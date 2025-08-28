package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"groupie-tracker/api"
)

type Location struct {
	Name        string    `json:"name"`
	Lat         string    `json:"lat"`
	Lon         string    `json:"lon"`
	BoundingBox [4]string `json:"boundingbox"`
	Type        string    `json:"type"`
	Dates       []string  `json:"dates"`
}

// Handler for the /locations/ route
func LocationsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

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

	response := make(map[string]Location)

	// range on the datesLocations keys
	for location, dates := range relation.DatesLocations {
		coordinates, err := api.GetCoordinates(location)
		if err != nil {
			http.Error(w, "Erreur lors de la récupération des coordonnées", http.StatusInternalServerError)
			return
		}
		response[location] = Location{
			Name:        coordinates[0].Name,
			Lat:         coordinates[0].Lat,
			Lon:         coordinates[0].Lon,
			BoundingBox: coordinates[0].BoundingBox,
			Type:        coordinates[0].Type,
			Dates:       dates,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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
