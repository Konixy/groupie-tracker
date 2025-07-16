package main

import (
	"log"
	"net/http"

	"groupie-tracker/handlers"
)

func main() {
	// Serve static files (CSS, images, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Route /
	http.HandleFunc("/", handlers.IndexHandler)

	// Start server
	port := ":3000"
	log.Printf("Server starting on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
