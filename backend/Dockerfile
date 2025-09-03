# Construction de l'application Go
FROM golang:1.24.3-alpine AS builder

# Installation des dépendances système nécessaires
RUN apk add --no-cache git

# Définition du répertoire de travail
WORKDIR /app

# Copie des fichiers de configuration Go
COPY go.mod go.sum* ./

# Téléchargement des dépendances
RUN go mod download

# Copie du code source
COPY . .

# Construction de l'application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Image finale optimisée
FROM alpine:latest

# Installation des certificats SSL pour les appels API externes
RUN apk --no-cache add ca-certificates

# Création d'un utilisateur non-root pour la sécurité
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

# Définition du répertoire de travail accessible à l'utilisateur
WORKDIR /app

# Copie de l'exécutable depuis l'étape de construction
COPY --from=builder /app/main .

# Copie des assets statiques (images des artistes)
COPY --from=builder /app/assets ./assets
COPY --from=builder /app/coordinates_cache.json ./coordinates_cache.json

# Changement des permissions pour l'utilisateur appuser
RUN chown -R appuser:appgroup /app

# Changement vers l'utilisateur non-root
USER appuser

# Exposition du port de l'API
EXPOSE 8080

# Commande de démarrage
CMD ["./main"]
