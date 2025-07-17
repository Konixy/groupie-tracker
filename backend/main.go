package main

import (
	"encoding/json"
	"log"
	"net/http"

	"groupie-tracker/handlers"
)

func main() {
	// Route API REST
	http.HandleFunc("/artists", handlers.IndexHandler)

	// Message d'accueil sur /
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Bienvenue sur l'API Groupie Tracker en Go"})
	})

	// Start server
	port := ":8080"
	log.Printf("Server starting on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
