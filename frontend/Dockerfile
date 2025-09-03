# Construction de l'application Svelte
FROM node:18-alpine AS builder

# Définition du répertoire de travail
WORKDIR /app

# Copie des fichiers de configuration npm
COPY package*.json ./

# Installation des dépendances
RUN npm i

# Copie du code source
COPY . .

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

