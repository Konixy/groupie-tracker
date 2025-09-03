package config

import "os"

// GetFrontendURL retourne l'URL du frontend depuis les variables d'environnement
func GetFrontendURL() string {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}
	return frontendURL
}

// GetAPIBaseURL retourne l'URL de base de l'API depuis les variables d'environnement
func GetAPIBaseURL() string {
	apiURL := os.Getenv("API_BASE_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}
	return apiURL
}

// GetPort retourne le port du serveur depuis les variables d'environnement
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
