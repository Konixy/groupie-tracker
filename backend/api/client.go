package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// BaseURL est l'adresse racine de l'API groupie tracker
const (
	BaseURL = "https://groupietrackers.herokuapp.com/api"
)

// Representation de l'api en "struct" Go
// Si tu vas sur l'url de l'api, tu verras que les données sont sous forme de json
// la c'est la meme chose mais traduit en struct Go
// Les tags `json:"..."` servent a faire le raccord entre le json de l'api et le struct Go
// Sans ça, go ne sait pas quoi faire avec les données de l'api
type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

// FetchArtists récupère tous les artistes depuis l'API Groupie Tracker
// Retourne une slice de structs Artist et une erreur si il y en a une (sinon "nil")
func FetchArtists() ([]Artist, error) {
	// Fait une requête HTTP GET vers l'adresse de base de l'api + /artists
	response, err := http.Get(BaseURL + "/artists")
	if err != nil {
		// Si la requête réseau échoue, on log l'erreur et on la retourne
		log.Printf("Error fetching artists: %v", err)
		return nil, err
	}
	// "defer" permet d'executer une fonction apres la fonction dans laquelle il est appelé
	// En gros, on ferme la connection HTTP apres que la fonction FetchArtists se termine
	// Si tu comprends pas trop cette ligne, c'est pas hyper grave
	defer response.Body.Close()

	// Lit et récupère tout le contenu de la réponse HTTP (les données JSON brutes)
	// ENTRÉE: response.Body = un "stream" de données (comme un tuyau d'eau qui coule)
	// SORTIE: body = toutes les données JSON sous forme de []byte (tableau d'octets)
	// POURQUOI: json.Unmarshal a besoin de TOUTES les données d'un coup, pas un stream
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	// Crée une slice vide pour contenir nos artistes après la conversion JSON
	var artists []Artist
	// json.Unmarshal convertit les bytes JSON en structs Go
	// Utilise les tags json:"..." pour mapper les champs JSON aux champs de struct
	err = json.Unmarshal(body, &artists)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		return nil, err
	}

	return artists, nil
}

// Concert représente un concert avec sa date et ses lieux
type Concert struct {
	Date      string   `json:"date"`
	Locations []string `json:"locations"`
}

// ArtistConcerts représente les concerts d'un artiste
type ArtistConcerts struct {
	ID      int                `json:"id"`
	Concerts map[string][]string `json:"concerts"`
}

// FetchArtistConcerts récupère les concerts d'un artiste spécifique depuis l'API Groupie Tracker
func FetchArtistConcerts(artistID int) (*ArtistConcerts, error) {
	// Récupère les dates de concerts
	datesResponse, err := http.Get(BaseURL + "/dates/" + fmt.Sprintf("%d", artistID))
	if err != nil {
		log.Printf("Error fetching dates: %v", err)
		return nil, err
	}
	defer datesResponse.Body.Close()

	// Récupère les lieux de concerts
	locationsResponse, err := http.Get(BaseURL + "/locations/" + fmt.Sprintf("%d", artistID))
	if err != nil {
		log.Printf("Error fetching locations: %v", err)
		return nil, err
	}
	defer locationsResponse.Body.Close()

	// Parse les dates
	var datesData struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	}
	
	datesBody, err := io.ReadAll(datesResponse.Body)
	if err != nil {
		log.Printf("Error reading dates response: %v", err)
		return nil, err
	}
	
	err = json.Unmarshal(datesBody, &datesData)
	if err != nil {
		log.Printf("Error unmarshaling dates JSON: %v", err)
		return nil, err
	}

	// Parse les lieux
	var locationsData struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
	}
	
	locationsBody, err := io.ReadAll(locationsResponse.Body)
	if err != nil {
		log.Printf("Error reading locations response: %v", err)
		return nil, err
	}
	
	err = json.Unmarshal(locationsBody, &locationsData)
	if err != nil {
		log.Printf("Error unmarshaling locations JSON: %v", err)
		return nil, err
	}

	// Combine les dates et lieux
	concerts := make(map[string][]string)
	
	// Assure que nous avons le même nombre de dates et de lieux
	minLen := len(datesData.Dates)
	if len(locationsData.Locations) < minLen {
		minLen = len(locationsData.Locations)
	}
	
	for i := 0; i < minLen; i++ {
		date := datesData.Dates[i]
		location := locationsData.Locations[i]
		
		// Nettoyer la date (enlever l'astérisque si présente)
		cleanDate := strings.TrimPrefix(date, "*")
		
		// Formater la location pour la rendre plus lisible
		formattedLocation := formatLocation(location)
		
		concerts[cleanDate] = append(concerts[cleanDate], formattedLocation)
	}

	return &ArtistConcerts{
		ID:       artistID,
		Concerts: concerts,
	}, nil
}

// formatLocation formate une location pour la rendre plus lisible
func formatLocation(location string) string {
	// Remplacer les underscores par des espaces et capitaliser
	parts := strings.Split(location, "-")
	if len(parts) >= 2 {
		city := strings.ReplaceAll(parts[0], "_", " ")
		country := strings.ToUpper(parts[1])
		return strings.Title(city) + ", " + country
	}
	
	// Si pas de format attendu, juste nettoyer les underscores
	return strings.Title(strings.ReplaceAll(location, "_", " "))
}
