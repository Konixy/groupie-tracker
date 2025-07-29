package handlers

import (
	"io"
	"net/http"

	"groupie-tracker/api"
)

func ImagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/images/queen.jpeg" {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		http.ServeFile(w, r, "assets/images/queen.jpeg")
		return
	}

	response, err := http.Get(api.BaseURL + r.URL.Path)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération de l'image", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.WriteHeader(response.StatusCode)
	io.Copy(w, response.Body)
}
