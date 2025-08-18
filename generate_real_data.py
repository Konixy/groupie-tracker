#!/usr/bin/env python3
"""
Script pour générer des données d'artistes depuis l'API Groupie Tracker.
Génère du code Go avec les données réelles des artistes, concerts et lieux.
"""

import requests
import json
from collections import defaultdict


def fetch_data():
    """Récupère et traite les données depuis l'API Groupie Tracker."""
    
    # Récupérer les artistes
    print("Récupération des artistes...")
    artists_url = "https://groupietrackers.herokuapp.com/api/artists"
    artists_response = requests.get(artists_url)
    artists_data = artists_response.json()
    
    # Récupérer les dates
    print("Récupération des dates de concerts...")
    dates_url = "https://groupietrackers.herokuapp.com/api/dates"
    dates_response = requests.get(dates_url)
    dates_data = dates_response.json()
    
    # Créer un map des dates par ID
    dates_map = {}
    all_years = set()
    
    for entry in dates_data['index']:
        artist_id = entry['id']
        dates = entry['dates']
        dates_map[artist_id] = dates
        
        # Extraire les années
        for date in dates:
            clean_date = date.replace('*', '')
            parts = clean_date.split('-')
            if len(parts) == 3:
                try:
                    year = int(parts[2])
                    all_years.add(year)
                except ValueError:
                    pass
    
    # Récupérer les relations pour les lieux
    print("Récupération des lieux de concerts...")
    locations_map = {}
    for artist in artists_data:
        artist_id = artist['id']
        relation_url = f"https://groupietrackers.herokuapp.com/api/relation/{artist_id}"
        try:
            relation_response = requests.get(relation_url)
            relation_data = relation_response.json()
            
            # Extraire les lieux uniques
            locations = set()
            for location_list in relation_data.get('datesLocations', {}).values():
                for location in location_list:
                    locations.add(location)
            
            locations_map[artist_id] = list(locations)
        except Exception as e:
            print(f"Erreur pour l'artiste {artist_id}: {e}")
            locations_map[artist_id] = []
    
    # Générer le code Go
    print("\n// Code généré automatiquement avec les vraies données de l'API")
    print("func serveStaticData(w http.ResponseWriter) {")
    print("    json.NewEncoder(w).Encode(map[string]interface{}{")
    print("        \"artists\": []map[string]interface{}{")
    
    for artist in artists_data:
        artist_id = artist['id']
        name = artist['name']
        dates = dates_map.get(artist_id, [])
        concert_count = len(dates)
        
        # Extraire les années pour cet artiste
        artist_years = set()
        for date in dates:
            clean_date = date.replace('*', '')
            parts = clean_date.split('-')
            if len(parts) == 3:
                try:
                    year = int(parts[2])
                    artist_years.add(year)
                except ValueError:
                    pass
        
        years_list = sorted(list(artist_years))
        locations = locations_map.get(artist_id, [])
        
        print(f"            {{\"id\": {artist_id}, \"name\": \"{name}\", \"image\": \"http://localhost:8080/images/id{artist_id}.jpg\", \"concertCount\": {concert_count}, \"years\": {years_list}, \"locations\": {json.dumps(locations)}}},")
    
    print("        },")
    print(f"        \"availableYears\": {sorted(list(all_years))},")
    print("    })")
    print("}")


if __name__ == "__main__":
    fetch_data() 