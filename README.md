# Groupie Tracker

Application web moderne pour explorer et suivre des artistes de musique et leurs concerts.

## 🚧 À FAIRE - Priorités de développement

- [ ] **Finir la map** - Compléter l'affichage et l'interactivité de la carte des concerts
- [ ] **Ajouter des filtres à la map** - Permettre de filtrer les concerts par date, lieu, artiste
- [x] **Barre de recherche universelle** - ✅ Implémentée avec recherche par artiste, membre, lieu et date
- [ ] **Finir le fond du site** - Ajouter des détails graphiques et améliorer l'esthétique générale
- [ ] **Responsive** - Adapter le site pour mobile (toutes taille d'écran au finale)

---

## Installation

### Prérequis

- Node.js 18+
- Go 1.21+

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
