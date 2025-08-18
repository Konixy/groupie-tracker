#!/usr/bin/env python3
"""
Script pour générer des données de statistiques depuis l'API Groupie Tracker.
Génère du code Go avec les statistiques des artistes et leurs concerts.
"""

import json
import subprocess


def fetch_stats_data():
    """Récupère et traite les données de statistiques depuis l'API."""
    
    print("Récupération des données d'artistes...")
    # Récupérer les artistes
    artists_cmd = "curl -s 'https://groupietrackers.herokuapp.com/api/artists' | jq '.[] | {id: .id, name: .name}' | jq -s '.'"
    artists_data = json.loads(subprocess.check_output(artists_cmd, shell=True, text=True))

    print("Récupération des données de concerts...")
    # Récupérer les dates
    dates_cmd = "curl -s 'https://groupietrackers.herokuapp.com/api/dates' | jq '.index[] | {id: .id, count: (.dates | length)}' | jq -s '.'"
    dates_data = json.loads(subprocess.check_output(dates_cmd, shell=True, text=True))

    # Créer un map des dates par ID
    dates_map = {item['id']: item['count'] for item in dates_data}

    # Générer les données Go
    print("\n// Données générées automatiquement depuis l'API")
    print("func StatsHandler(w http.ResponseWriter, r *http.Request) {")
    print("	w.Header().Set(\"Content-Type\", \"application/json\")")
    print("	w.Header().Set(\"Access-Control-Allow-Origin\", \"http://localhost:3000\")")
    print("	w.Header().Set(\"Access-Control-Allow-Methods\", \"GET, OPTIONS\")")
    print("	w.Header().Set(\"Access-Control-Allow-Headers\", \"Content-Type\")")
    print("")
    print("	if r.Method == \"OPTIONS\" {")
    print("		w.WriteHeader(http.StatusOK)")
    print("		return")
    print("	}")
    print("")
    print("	// Données statiques basées sur l'API /dates et /artists")
    print("	json.NewEncoder(w).Encode(map[string]interface{}{")
    print("		\"artists\": []map[string]interface{}{")

    for artist in artists_data:
        artist_id = artist['id']
        artist_name = artist['name'].replace('"', '\\"')  # Échapper les guillemets
        concert_count = dates_map.get(artist_id, 0)
        print(f"			{{\"id\": {artist_id}, \"name\": \"{artist_name}\", \"image\": \"http://localhost:8080/images/id{artist_id}.jpg\", \"concertCount\": {concert_count}}},")

    print("		},")
    print("	})")
    print("}")


if __name__ == "__main__":
    fetch_stats_data() 