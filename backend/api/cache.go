package api

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

// Cache stocke les coordonnées géographiques en mémoire et sur disque
// Structure thread-safe pour éviter les problèmes de concurrence
type Cache struct {
	data map[string]GeocodeResponse // map[location]coordonnées
	mu   sync.RWMutex               // mutex pour la sécurité thread
}

var (
	// Instance globale du cache
	coordinatesCache *Cache
	// Nom du fichier où sauvegarder le cache
	cacheFileName = "coordinates_cache.json"
)

// init initialise le cache au démarrage du programme
func init() {
	coordinatesCache = &Cache{
		data: make(map[string]GeocodeResponse),
	}
	loadCacheFromDisk()
}

// GetFromCache récupère une coordonnée depuis le cache
// Retourne les coordonnées et un booléen indiquant si elles ont été trouvées
func GetFromCache(location string) (GeocodeResponse, bool) {
	coordinatesCache.mu.RLock() // verrou en lecture
	defer coordinatesCache.mu.RUnlock()

	response, found := coordinatesCache.data[location]
	return response, found
}

// SaveToCache sauvegarde une coordonnée dans le cache en mémoire et sur disque
func SaveToCache(location string, response GeocodeResponse) {
	coordinatesCache.mu.Lock() // verrou en écriture
	defer coordinatesCache.mu.Unlock()

	// Sauvegarder en mémoire
	coordinatesCache.data[location] = response

	// Sauvegarder sur disque de manière asynchrone pour ne pas bloquer
	go saveCacheToDisk()
}

// loadCacheFromDisk charge le cache depuis le fichier JSON
func loadCacheFromDisk() {
	file, err := os.Open(cacheFileName)
	if err != nil {
		// Si le fichier n'existe pas, ce n'est pas une erreur grave
		if os.IsNotExist(err) {
			log.Printf("Fichier de cache non trouvé, création d'un nouveau cache")
			return
		}
		log.Printf("Erreur lors de l'ouverture du fichier de cache: %v", err)
		return
	}
	defer file.Close()

	var cacheData map[string]GeocodeResponse
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cacheData)
	if err != nil {
		log.Printf("Erreur lors du décodage du cache: %v", err)
		return
	}

	coordinatesCache.mu.Lock()
	coordinatesCache.data = cacheData
	coordinatesCache.mu.Unlock()

	log.Printf("Cache chargé avec %d entrées", len(cacheData))
}

// saveCacheToDisk sauvegarde le cache actuel dans le fichier JSON
func saveCacheToDisk() {
	coordinatesCache.mu.RLock()
	// Créer une copie des données pour éviter de garder le verrou trop longtemps
	dataCopy := make(map[string]GeocodeResponse)
	for k, v := range coordinatesCache.data {
		dataCopy[k] = v
	}
	coordinatesCache.mu.RUnlock()

	file, err := os.Create(cacheFileName)
	if err != nil {
		log.Printf("Erreur lors de la création du fichier de cache: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Format JSON lisible
	err = encoder.Encode(dataCopy)
	if err != nil {
		log.Printf("Erreur lors de l'encodage du cache: %v", err)
		return
	}
}

// GetCacheStats retourne des statistiques sur le cache
func GetCacheStats() map[string]int {
	coordinatesCache.mu.RLock()
	defer coordinatesCache.mu.RUnlock()

	return map[string]int{
		"entries": len(coordinatesCache.data),
	}
}
