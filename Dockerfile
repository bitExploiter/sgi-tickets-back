FROM golang:1.24.0-bullseye as builder

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

# Copiar directorios necesarios
COPY templates ./templates
COPY public ./public

# Crear directorio para archivos subidos (con permisos)
RUN mkdir -p ./files && chmod 755 ./files

# Copiar binario compilado
COPY --from=builder /app/main .

# NOTA: El .env se debe pasar vía variables de entorno (docker-compose)
# o montar como volumen. Para producción standalone, crear .env-produccion
# y descomentear la siguiente línea antes del build:
# COPY .env-produccion .env

EXPOSE 8080

CMD ["./main"]
