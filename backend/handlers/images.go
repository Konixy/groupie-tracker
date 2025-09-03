package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"groupie-tracker/config"
)

func ImagesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", config.GetFrontendURL())

	// Extraire le nom du fichier de l'URL
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "URL d'image invalide", http.StatusBadRequest)
		return
	}

	// Le nom du fichier est la dernière partie de l'URL
	filename := parts[len(parts)-1]

	// Construire le chemin vers l'image locale
	localImagePath := filepath.Join("assets", "images", "artistspict", filename)

	// Vérifier si l'image locale existe
	if strings.HasSuffix(filename, ".jpg") {
		// Vérifier si le fichier existe
		if _, err := os.Stat(localImagePath); err == nil {
			// Définir les headers CORS
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			w.Header().Set("Content-Type", "image/jpeg")

			// Lire et servir le fichier manuellement
			file, err := os.Open(localImagePath)
			if err != nil {
				log.Printf("Erreur lors de l'ouverture de l'image: %v", err)
				http.Error(w, "Erreur lors de l'ouverture de l'image", http.StatusInternalServerError)
				return
			}
			defer file.Close()

			io.Copy(w, file)
			return
		}
	}

	http.Error(w, "Image non trouvée", http.StatusNotFound)
}
