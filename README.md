# SGI Tickets Backend

Backend REST API para Sistema de Gestión de Interventorías.

**Stack**: Go 1.24.0 + Fiber v2 + GORM + PostgreSQL

---

## 🚀 Quick Start (Desarrollo Local)

### Prerrequisitos
- [Docker](https://docs.docker.com/get-docker/) y Docker Compose instalados
- Git

### Iniciar el proyecto

```bash
# 1. Clonar repositorio
git clone <url-del-repo>
cd sgi-tickets-back

# 2. Levantar servicios (PostgreSQL + Backend)
docker-compose up -d

# 3. Ver logs para verificar que todo está corriendo
docker-compose logs -f app
```

¡Listo! El backend estará disponible en `http://localhost:8080`

### Endpoints disponibles

- **Health check**: `GET http://localhost:8080/api/version`
- **API Base**: `http://localhost:8080/api/v1/`
- **Autenticación**: `http://localhost:8080/api/v1/auth/`

### Comandos útiles

```bash
# Ver logs del backend
docker-compose logs -f app

# Ver logs de la base de datos
docker-compose logs -f db

# Detener servicios
docker-compose down

# Reiniciar servicios
docker-compose restart

# Limpiar todo (incluyendo la BD)
docker-compose down -v

# Reconstruir las imágenes
docker-compose up -d --build
```

### Conectar a la base de datos

Si necesitas conectarte directamente a PostgreSQL:

```
Host: localhost
Port: 5432
Database: sgi2_db
User: postgres
Password: abcd.1234
```

---

## 🔧 Desarrollo Local (sin Docker)

Si prefieres ejecutar Go directamente:

### Prerrequisitos
- Go 1.24.0
- PostgreSQL 16 instalado y corriendo

### Configuración

```bash
# 1. Instalar dependencias
go mod download && go mod tidy

# 2. Crear archivo .env (copiar desde .env.example)
cp .env.example .env

# 3. Editar .env con tus credenciales de BD local

# 4. Instalar Air para hot reload (opcional)
go install github.com/air-verse/air@latest

# 5. Ejecutar con hot reload
air

# O ejecutar directamente
go run main.go
```

---

## 📁 Estructura del proyecto

```
handlers/          # HTTP handlers + middleware
models/            # Modelos GORM
toolbox/           # Lógica de negocio y utilidades
storage/           # Conexión a BD
migrations/        # Migraciones de base de datos
templates/         # Templates HTML (emails, PDFs)
public/            # Frontend React compilado
files/             # Archivos subidos por usuarios
```

---

## 🐛 Troubleshooting

### El backend no inicia
```bash
# Ver logs completos
docker-compose logs app

# Verificar que PostgreSQL esté ready
docker-compose logs db
```

### Error de conexión a BD
```bash
# Reiniciar servicios
docker-compose restart

# O reconstruir desde cero
docker-compose down -v
docker-compose up -d
```

### Cambios en el código no se reflejan
```bash
# Reconstruir la imagen
docker-compose up -d --build
```

---

## 📚 Documentación adicional

- [README-DOCKER.md](README-DOCKER.md) - Guía completa de Docker
- [CLAUDE.md](CLAUDE.md) - Patrones de código y arquitectura del proyecto

---

## 🤝 Para el equipo frontend

El backend expone una API REST en `http://localhost:8080/api/v1/`

Para probar endpoints protegidos, primero debes:
1. Login: `POST /api/v1/auth/login`
2. Setup 2FA: `GET /api/v1/auth/2fa/setup`
3. Verificar 2FA: `POST /api/v1/auth/2fa/verify`

Las cookies de sesión se manejan automáticamente.
