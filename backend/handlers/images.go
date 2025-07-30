package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"groupie-tracker/api"
)

func ImagesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("ImagesHandler appelé pour: %s", r.URL.Path)
	
	// Gérer les requêtes OPTIONS pour CORS
	if r.Method == "OPTIONS" {
		log.Printf("Requête OPTIONS détectée")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Extraire le nom du fichier de l'URL
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "URL d'image invalide", http.StatusBadRequest)
		return
	}

	// Le nom du fichier est la dernière partie de l'URL
	filename := parts[len(parts)-1]
	log.Printf("Nom du fichier: %s", filename)
	
	// Construire le chemin vers l'image locale
	localImagePath := filepath.Join("assets", "images", "artistspict", filename)
	log.Printf("Chemin local: %s", localImagePath)
	
	// Vérifier si l'image locale existe
	if strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") || strings.HasSuffix(filename, ".png") {
		// Vérifier si le fichier existe
		if _, err := os.Stat(localImagePath); err == nil {
			log.Printf("Image locale trouvée: %s", localImagePath)
			
			// Définir les headers CORS
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			
			// Déterminer le Content-Type basé sur l'extension
			contentType := "image/jpeg"
			if strings.HasSuffix(filename, ".png") {
				contentType = "image/png"
			}
			w.Header().Set("Content-Type", contentType)
			
			log.Printf("Headers CORS définis pour l'image locale")
			
			// Lire et servir le fichier manuellement
			file, err := os.Open(localImagePath)
			if err != nil {
				log.Printf("Erreur lors de l'ouverture de l'image: %v", err)
				http.Error(w, "Erreur lors de l'ouverture de l'image", http.StatusInternalServerError)
				return
			}
			defer file.Close()
			
			io.Copy(w, file)
			log.Printf("Image locale servie avec succès")
			return
		} else {
			log.Printf("Image locale non trouvée: %s", localImagePath)
		}
	}

	log.Printf("Utilisation de l'API distante pour: %s", r.URL.Path)
	// Si l'image locale n'existe pas, essayer de la récupérer depuis l'API distante
	response, err := http.Get(api.BaseURL + r.URL.Path)
	if err != nil {
		log.Printf("Erreur lors de la récupération depuis l'API distante: %v", err)
		http.Error(w, "Erreur lors de la récupération de l'image", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.WriteHeader(response.StatusCode)
	io.Copy(w, response.Body)
	log.Printf("Image distante servie avec succès")
}
