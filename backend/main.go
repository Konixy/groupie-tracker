package main

import (
	"encoding/json"
	"log"
	"net/http"

	"groupie-tracker/handlers"
)

func main() {
	// Route API REST
	http.HandleFunc("/artists", handlers.ArtistsHandler)
	http.HandleFunc("/artists/", handlers.ArtistConcertsHandler)
	http.HandleFunc("/images/", handlers.ImagesHandler)
	http.HandleFunc("/locations/", handlers.LocationsHandler)
	http.HandleFunc("/all-locations", handlers.AllLocationsHandler)

	// Message d'accueil sur /
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		json.NewEncoder(w).Encode(map[string]string{"message": "Bienvenue sur l'API Groupie Tracker en Go"})
	})

	// Start server
	port := ":8080"
	log.Printf("Server starting on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
