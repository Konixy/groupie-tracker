# 🎵 Groupie Tracker

Application web moderne pour explorer et suivre des artistes de musique et leurs concerts.

## 🏗️ Architecture

- **Frontend** : Svelte 5 avec TypeScript
- **Backend** : Go avec API REST
- **Données** : API publique Groupie Tracker

## 🚀 Installation

### Prérequis
- Node.js 18+
- Go 1.21+
- Python 3.8+ (pour les scripts de génération de données)

### Backend
```bash
cd backend
go mod tidy
go run main.go
```
Le serveur démarre sur `http://localhost:8080`

### Frontend
```bash
cd frontend
npm install
npm run dev
```
L'application démarre sur `http://localhost:3000`

## 📁 Structure du projet

```
groupie-tracker/
├── backend/           # API Go
│   ├── api/          # Client API
│   ├── handlers/     # Gestionnaires HTTP
│   ├── assets/       # Images des artistes
│   └── main.go       # Point d'entrée
├── frontend/         # Interface Svelte
│   ├── src/
│   │   ├── lib/      # Composants
│   │   ├── App.svelte
│   │   └── types.ts
│   └── package.json
└── scripts/          # Scripts Python
    ├── generate_real_data.py
    └── generate_stats_data.py
```

## 🎯 Fonctionnalités

- **Carousel interactif** des artistes avec navigation 3D
- **Recherche** et filtrage d'artistes
- **Détails des artistes** avec informations complètes
- **Carte des concerts** avec localisation
- **Statistiques** et visualisations
- **Design responsive** et animations fluides

## 🛠️ Scripts utiles

### Génération de données
```bash
# Générer les données d'artistes
python generate_real_data.py

# Générer les statistiques
python generate_stats_data.py
```

## 🎨 Technologies utilisées

- **Frontend** : Svelte 5, TypeScript, Vite
- **Backend** : Go, HTTP/JSON
- **Styling** : CSS moderne avec variables
- **Images** : Extraction de couleurs avec Vibrant.js
- **Données** : API REST avec CORS

## 📝 Licence

Projet éducatif - Libre d'utilisation 