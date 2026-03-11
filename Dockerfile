FROM golang:1.24.0-bullseye AS builder

WORKDIR /app

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar binario
RUN go build -o /app/main .

# Imagen final
FROM debian:bullseye-slim

WORKDIR /app

# Instalar postgresql-client para ejecutar seeders
RUN apt-get update && apt-get install -y postgresql-client && rm -rf /var/lib/apt/lists/*

# Copiar directorios necesarios
COPY templates ./templates
COPY public ./public
COPY seeds ./seeds
COPY docs ./docs

# Crear directorio para archivos subidos (con permisos)
RUN mkdir -p ./files && chmod 755 ./files

# Copiar binario compilado
COPY --from=builder /app/main .

# Copiar script de entrypoint
COPY docker-entrypoint.sh /app/
RUN chmod +x /app/docker-entrypoint.sh

# NOTA: El .env se debe pasar vía variables de entorno (docker-compose)
# o montar como volumen. Para producción standalone, crear .env-produccion
# y descomentear la siguiente línea antes del build:
# COPY .env-produccion .env

EXPOSE 8080

ENTRYPOINT ["/app/docker-entrypoint.sh"]
