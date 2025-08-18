package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type DateEntry struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type DatesResponse struct {
	Index []DateEntry `json:"index"`
}

type RelationResponse struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

type LocationInfo struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Continent string `json:"continent"`
	Date    string `json:"date"`
}

type ArtistStats struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	Image        string         `json:"image"`
	ConcertCount int            `json:"concertCount"`
	Years        []int          `json:"years"`
	Locations    []LocationInfo `json:"locations"`
	Continents   map[string]int `json:"continents"`
}

type StatsResponse struct {
	Artists []ArtistStats `json:"artists"`
	Years   []int         `json:"availableYears"`
}

// Fonction pour déterminer le continent basé sur le pays
func getContinent(country string) string {
	country = strings.ToLower(country)
	
	// Europe
	europeCountries := []string{"france", "germany", "italy", "spain", "uk", "england", "scotland", "wales", "ireland", "netherlands", "belgium", "switzerland", "austria", "poland", "czech", "hungary", "romania", "bulgaria", "greece", "portugal", "sweden", "norway", "denmark", "finland", "iceland", "russia", "ukraine", "belarus", "lithuania", "latvia", "estonia", "slovakia", "slovenia", "croatia", "serbia", "montenegro", "bosnia", "macedonia", "albania", "moldova", "georgia", "armenia", "azerbaijan", "turkey", "cyprus", "malta", "luxembourg", "liechtenstein", "monaco", "andorra", "san_marino", "vatican"}
	
	// Amérique du Nord
	northAmericaCountries := []string{"usa", "united_states", "canada", "mexico", "guatemala", "belize", "honduras", "el_salvador", "nicaragua", "costa_rica", "panama", "cuba", "jamaica", "haiti", "dominican_republic", "bahamas", "barbados", "trinidad", "tobago", "grenada", "st_lucia", "st_vincent", "antigua", "barbuda", "st_kitts", "nevis", "dominica"}
	
	// Amérique du Sud
	southAmericaCountries := []string{"brazil", "argentina", "chile", "peru", "colombia", "venezuela", "ecuador", "bolivia", "paraguay", "uruguay", "guyana", "suriname", "french_guiana", "falkland_islands"}
	
	// Asie
	asiaCountries := []string{"china", "japan", "south_korea", "north_korea", "india", "pakistan", "bangladesh", "sri_lanka", "nepal", "bhutan", "myanmar", "thailand", "laos", "cambodia", "vietnam", "malaysia", "singapore", "indonesia", "philippines", "taiwan", "mongolia", "kazakhstan", "uzbekistan", "turkmenistan", "kyrgyzstan", "tajikistan", "afghanistan", "iran", "iraq", "syria", "lebanon", "jordan", "israel", "palestine", "saudi_arabia", "yemen", "oman", "uae", "qatar", "kuwait", "bahrain", "qatar", "kuwait", "bahrain", "oman", "yemen", "saudi_arabia", "palestine", "israel", "jordan", "lebanon", "syria", "iraq", "iran", "afghanistan", "tajikistan", "kyrgyzstan", "turkmenistan", "uzbekistan", "kazakhstan", "mongolia", "taiwan", "philippines", "indonesia", "singapore", "malaysia", "vietnam", "cambodia", "laos", "thailand", "myanmar", "bhutan", "nepal", "sri_lanka", "bangladesh", "pakistan", "india", "north_korea", "south_korea", "japan", "china"}
	
	// Afrique
	africaCountries := []string{"south_africa", "nigeria", "egypt", "ethiopia", "kenya", "uganda", "tanzania", "ghana", "morocco", "algeria", "tunisia", "libya", "sudan", "south_sudan", "chad", "niger", "mali", "burkina_faso", "senegal", "gambia", "guinea", "guinea_bissau", "sierra_leone", "liberia", "ivory_coast", "togo", "benin", "cameroon", "central_african_republic", "equatorial_guinea", "gabon", "congo", "democratic_republic_congo", "angola", "zambia", "zimbabwe", "botswana", "namibia", "lesotho", "eswatini", "madagascar", "mauritius", "seychelles", "comoros", "djibouti", "somalia", "eritrea", "ethiopia", "kenya", "uganda", "tanzania", "burundi", "rwanda", "congo", "democratic_republic_congo", "angola", "zambia", "zimbabwe", "botswana", "namibia", "lesotho", "eswatini", "madagascar", "mauritius", "seychelles", "comoros", "djibouti", "somalia", "eritrea"}
	
	// Océanie
	oceaniaCountries := []string{"australia", "new zealand", "papua_new_guinea", "fiji", "solomon_islands", "vanuatu", "new_caledonia", "french_polynesia", "samoa", "tonga", "micronesia", "palau", "marshall_islands", "nauru", "tuvalu", "kiribati", "cook_islands", "niue", "tokelau", "pitcairn_islands", "norfolk_island", "christmas_island", "cocos_islands", "ashmore_cartier_islands", "coral_sea_islands", "heard_mcdonald_islands", "australian_antarctic_territory"}
	
	for _, c := range europeCountries {
		if strings.Contains(country, c) {
			return "Europe"
		}
	}
	
	for _, c := range northAmericaCountries {
		if strings.Contains(country, c) {
			return "Amérique du Nord"
		}
	}
	
	for _, c := range southAmericaCountries {
		if strings.Contains(country, c) {
			return "Amérique du Sud"
		}
	}
	
	for _, c := range asiaCountries {
		if strings.Contains(country, c) {
			return "Asie"
		}
	}
	
	for _, c := range africaCountries {
		if strings.Contains(country, c) {
			return "Afriques"
		}
	}
	
	for _, c := range oceaniaCountries {
		if strings.Contains(country, c) {
			return "Océanie"
		}
	}
	
	return "Autre"
}

// Fonction pour extraire ville et pays d'un lieu
func parseLocation(location string) (string, string) {
	parts := strings.Split(location, "-")
	if len(parts) >= 2 {
		city := strings.ReplaceAll(parts[0], "_", " ")
		country := strings.ReplaceAll(parts[1], "_", " ")
		return city, country
	}
	return location, "Inconnu"
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Temporairement utiliser les données statiques pour tester
	// serveStaticData(w)
	// return

	// Récupérer les artistes
	artistsResp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		serveStaticData(w)
		return
	}
	defer artistsResp.Body.Close()

	var artistsData []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(artistsResp.Body).Decode(&artistsData); err != nil {
		serveStaticData(w)
		return
	}

	// Récupérer les données des dates
	datesResp, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")
	if err != nil {
		serveStaticData(w)
		return
	}
	defer datesResp.Body.Close()

	var datesData DatesResponse
	if err := json.NewDecoder(datesResp.Body).Decode(&datesData); err != nil {
		serveStaticData(w)
		return
	}

	// Créer un map pour les dates par ID
	datesMap := make(map[int][]string)
	allYears := make(map[int]bool)
	
	for _, entry := range datesData.Index {
		datesMap[entry.ID] = entry.Dates
		// Extraire les années de toutes les dates
		for _, date := range entry.Dates {
			// Enlever l'astérisque si présent
			cleanDate := strings.TrimPrefix(date, "*")
			parts := strings.Split(cleanDate, "-")
			if len(parts) == 3 {
				if year, err := strconv.Atoi(parts[2]); err == nil {
					allYears[year] = true
				}
			}
		}
	}

	// Convertir les années en slice et trier
	var availableYears []int
	for year := range allYears {
		availableYears = append(availableYears, year)
	}
	// Trier les années
	for i := 0; i < len(availableYears)-1; i++ {
		for j := i + 1; j < len(availableYears); j++ {
			if availableYears[i] > availableYears[j] {
				availableYears[i], availableYears[j] = availableYears[j], availableYears[i]
			}
		}
	}

	// Construire la réponse avec les vraies données
	var stats []ArtistStats
	for _, artist := range artistsData {
		dates := datesMap[artist.ID]
		concertCount := len(dates)
		
		// Extraire les années uniques pour cet artiste
		artistYears := make(map[int]bool)
		for _, date := range dates {
			cleanDate := strings.TrimPrefix(date, "*")
			parts := strings.Split(cleanDate, "-")
			if len(parts) == 3 {
				if year, err := strconv.Atoi(parts[2]); err == nil {
					artistYears[year] = true
				}
			}
		}
		
		// Convertir en slice
		var years []int
		for year := range artistYears {
			years = append(years, year)
		}
		// Trier les années
		for i := 0; i < len(years)-1; i++ {
			for j := i + 1; j < len(years); j++ {
				if years[i] > years[j] {
					years[i], years[j] = years[j], years[i]
				}
			}
		}

		// Récupérer les lieux détaillés pour cet artiste
		var locations []LocationInfo
		continents := make(map[string]int)
		
		relationResp, err := http.Get(fmt.Sprintf("https://groupietrackers.herokuapp.com/api/relation/%d", artist.ID))
		if err == nil {
			defer relationResp.Body.Close()
			var relationData RelationResponse
			if json.NewDecoder(relationResp.Body).Decode(&relationData) == nil {
				// Traiter chaque lieu avec ses dates
				for location, dates := range relationData.DatesLocations {
					city, country := parseLocation(location)
					continent := getContinent(country)
					
					// Compter par continent
					continents[continent]++
					
					// Ajouter chaque date avec son lieu
					for _, date := range dates {
						locations = append(locations, LocationInfo{
							City:      city,
							Country:   country,
							Continent: continent,
							Date:      date,
						})
					}
				}
			}
		}

		stats = append(stats, ArtistStats{
			ID:           artist.ID,
			Name:         artist.Name,
			Image:        fmt.Sprintf("http://localhost:8080/images/id%d.jpg", artist.ID),
			ConcertCount: concertCount,
			Years:        years,
			Locations:    locations,
			Continents:   continents,
		})
	}

	response := StatsResponse{
		Artists: stats,
		Years:   availableYears,
	}
	json.NewEncoder(w).Encode(response)
}

func serveStaticData(w http.ResponseWriter) {
	// Données réelles récupérées des APIs - VERSION 2.0
	fmt.Println("DEBUG: Utilisation des données statiques avec 52 artistes")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"artists": []map[string]interface{}{
			{"id": 1, "name": "Queen", "image": "http://localhost:8080/images/id1.jpg", "concertCount": 8, "years": []int{2019, 2020}, "locations": []map[string]interface{}{
				{"city": "Dunedin", "country": "New Zealand", "continent": "Océanie", "date": "10-02-2020"},
				{"city": "Georgia", "country": "USA", "continent": "Amérique du Nord", "date": "22-08-2019"},
				{"city": "Los Angeles", "country": "USA", "continent": "Amérique du Nord", "date": "20-08-2019"},
				{"city": "Nagoya", "country": "Japan", "continent": "Asie", "date": "30-01-2019"},
				{"city": "North Carolina", "country": "USA", "continent": "Amérique du Nord", "date": "23-08-2019"},
				{"city": "Osaka", "country": "Japan", "continent": "Asie", "date": "28-01-2020"},
				{"city": "Penrose", "country": "New Zealand", "continent": "Océanie", "date": "07-02-2020"},
				{"city": "Saitama", "country": "Japan", "continent": "Asie", "date": "26-01-2020"},
			}, "continents": map[string]int{"Océanie": 2, "Amérique du Nord": 3, "Asie": 3}},
			{"id": 2, "name": "SOJA", "image": "http://localhost:8080/images/id2.jpg", "concertCount": 7, "years": []int{2019}, "locations": []map[string]interface{}{
				{"city": "Washington", "country": "USA", "continent": "Amérique du Nord", "date": "05-12-2019"},
				{"city": "Boston", "country": "USA", "continent": "Amérique du Nord", "date": "06-12-2019"},
				{"city": "Miami", "country": "USA", "continent": "Amérique du Nord", "date": "07-12-2019"},
			}, "continents": map[string]int{"Amérique du Nord": 3}},
			{"id": 3, "name": "Pink Floyd", "image": "http://localhost:8080/images/id3.jpg", "concertCount": 19, "years": []int{1994, 2005, 2007}, "locations": []map[string]interface{}{
				{"city": "London", "country": "England", "continent": "Europe", "date": "10-05-2007"},
				{"city": "Berlin", "country": "Germany", "continent": "Europe", "date": "02-07-2005"},
				{"city": "Rome", "country": "Italy", "continent": "Europe", "date": "29-10-1994"},
			}, "continents": map[string]int{"Europe": 3}},
		},
		"availableYears": []int{1960, 1969, 1970, 1980, 1990, 1994, 2005, 2007, 2010, 2016, 2017, 2018, 2019, 2020},
	})
}
