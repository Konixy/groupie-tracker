# Construction de l'application Svelte
FROM node:current-alpine AS builder

# Arguments de build pour les variables d'environnement
ARG VITE_API_BASE_URL=http://localhost:8080

# Définition du répertoire de travail
WORKDIR /app

# Copie des fichiers de configuration npm
COPY package*.json ./

# Installation des dépendances
RUN npm i

# Copie du code source
COPY . .

# Configuration des variables d'environnement pour le build
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

# Construction de l'application pour la production
RUN npm run build

# Image finale avec serveur Nginx pour servir l'application
FROM nginx:alpine

# Copie de l'application construite
COPY --from=builder /app/dist /usr/share/nginx/html

# Copie de la configuration Nginx personnalisée
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Exposition du port HTTP
EXPOSE 80

# Démarrage de Nginx
CMD ["nginx", "-g", "daemon off;"]

