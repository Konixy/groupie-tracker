// Configuration pour l'application frontend
// Utilise les variables d'environnement Vite (VITE_*)

// URL de base de l'API backend
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

// Autres configurations peuvent être ajoutées ici
export const config = {
	apiBaseUrl: API_BASE_URL
};
