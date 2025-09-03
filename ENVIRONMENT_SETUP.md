# Configuration des Variables d'Environnement - Groupie Tracker

Ce document explique comment configurer les variables d'environnement pour différents environnements de déploiement.

## Variables d'environnement Backend (Go)

Le backend utilise les variables suivantes :

- `PORT` : Port sur lequel le serveur écoute (défaut: 8080)
- `FRONTEND_URL` : URL du frontend pour CORS (défaut: http://localhost:3000)
- `API_BASE_URL` : URL de base de l'API pour les images (défaut: http://localhost:8080)

### Exemple backend en développement local :

```bash
export PORT=8080
export FRONTEND_URL=http://localhost:3000
export API_BASE_URL=http://localhost:8080
```

## Variables d'environnement Frontend (Svelte/Vite)

Le frontend utilise les variables suivantes (préfixées par `VITE_`) :

- `VITE_API_BASE_URL` : URL de l'API backend (défaut: http://localhost:8080)

### Fichier .env pour le développement local :

```env
VITE_API_BASE_URL=http://localhost:8080
```

### Fichier .env.production pour la production :

```env
VITE_API_BASE_URL=http://backend:8080
```

## Configuration Docker

### Docker Compose

Le fichier `docker-compose.yml` configure automatiquement les variables d'environnement :

```yaml
services:
  backend:
    environment:
      - PORT=8080
      - FRONTEND_URL=http://frontend
      - API_BASE_URL=http://backend:8080

  frontend:
    build:
      args:
        - VITE_API_BASE_URL=http://localhost:89
```

### Variables spécifiques à Docker

- Backend : Utilise les noms de services Docker (`http://frontend`, `http://backend:8080`)
- Frontend : Utilise l'URL publique accessible depuis l'hôte (`http://localhost:89`)

## Déploiement en production

Pour un déploiement en production, modifiez les variables dans `docker-compose.yml` :

```yaml
services:
  backend:
    environment:
      - FRONTEND_URL=https://votre-domaine.com
      - API_BASE_URL=https://api.votre-domaine.com

  frontend:
    build:
      args:
        - VITE_API_BASE_URL=https://api.votre-domaine.com
```

## Développement local

### Backend seul

```bash
cd backend
export FRONTEND_URL=http://localhost:3000
go run main.go
```

### Frontend seul

```bash
cd frontend
echo "VITE_API_BASE_URL=http://localhost:8080" > .env
npm run dev
```

## Vérification

Pour vérifier que les variables sont correctement configurées :

### Backend

- Vérifiez les logs de démarrage pour le port
- Testez l'endpoint `/` pour voir si CORS fonctionne

### Frontend

- Ouvrez les outils de développement du navigateur
- Vérifiez que les appels API utilisent la bonne URL dans l'onglet Network
