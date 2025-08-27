package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// GeocodeResponse représente la réponse de l'API Nominatim d'OpenStreetMap
type GeocodeResponse [1]struct {
	PlaceID     int       `json:"place_id"`
	Licence     string    `json:"licence"`
	OsmType     string    `json:"osm_type"`
	OsmID       int       `json:"osm_id"`
	BoundingBox [4]string `json:"boundingbox"`
	Lat         string    `json:"lat"`
	Lon         string    `json:"lon"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Category    string    `json:"category"`
	Type        string    `json:"addresstype"`
	Importance  float64   `json:"importance"`
}

// GetCoordinates récupère les coordonnées géographiques d'un lieu donné
// Utilise l'API Nominatim d'OpenStreetMap pour convertir "ville-pays" en latitude/longitude
// Vérifie d'abord le cache avant de faire un appel API
func GetCoordinates(location string) (GeocodeResponse, error) {
	// Vérifier le cache d'abord
	if cachedResponse, found := GetFromCache(location); found {
		return cachedResponse, nil
	}

	// Divise la location en ville et pays (format: "ville-pays")
	str := strings.Split(location, "-")
	if len(str) < 2 {
		return GeocodeResponse{}, fmt.Errorf("format de location invalide: %s", location)
	}

	city := str[0]
	country := str[1]

	// Appel à l'API Nominatim
	response, err := http.Get(fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s,%s&format=jsonv2&accept-language=fr", city, country))
	if err != nil {
		log.Printf("Erreur lors de la récupération des coordonnées: %v", err)
		return GeocodeResponse{}, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Erreur lors de la lecture de la réponse: %v", err)
		return GeocodeResponse{}, err
	}

	var geocodeResponse GeocodeResponse
	err = json.Unmarshal(body, &geocodeResponse)
	if err != nil {
		log.Printf("Erreur lors du parsing JSON: %v", err)
		return GeocodeResponse{}, err
	}

	// Sauvegarder dans le cache pour les prochaines fois
	SaveToCache(location, geocodeResponse)

	return geocodeResponse, nil
}
