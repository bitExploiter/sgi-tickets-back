# 🔑 Recuperación de Contraseña - Guía de Testing

## 🎯 Flujo Completo

### Método 1: Script Automatizado (Recomendado para Testing)

El script inserta un token conocido en la BD y prueba todo el flujo automáticamente:

```bash
./seeds/test_reset_password_complete.sh
```

Este script hace:
1. ✅ Genera un token de prueba conocido
2. ✅ Calcula el hash SHA-256 del token
3. ✅ Inserta el token en la BD con expiración de 1 hora
4. ✅ Resetea la contraseña usando el token
5. ✅ Verifica el login con la nueva contraseña

**Resultado:**
- Nueva contraseña: `nueva_password_456`
- Para restaurar la original: `./seeds/run_seeders.sh`

---

### Método 2: Flujo Real con Email

Si tienes SMTP configurado en el [.env](../.env), puedes probar el flujo real:

#### Paso 1: Solicitar recuperación

```bash
curl -X POST http://localhost:8080/api/v1/auth/recover \
  -H 'Content-Type: application/json' \
  -d '{
    "email": "entidad@sgi.gov.co"
  }'
```

**Respuesta:**
```json
{
  "success": true,
  "message": "Si el email existe, recibiras instrucciones de recuperacion"
}
```

**Nota:** Por seguridad, siempre devuelve éxito (no revela si el email existe).

#### Paso 2: Revisar email

Busca un email con asunto: **"Recuperar contraseña - SGI Tickets"**

El email contiene un enlace como:
```
http://localhost:5173/reset-password?token=abcd1234567890...
```

#### Paso 3: Resetear contraseña

```bash
curl -X POST http://localhost:8080/api/v1/auth/reset \
  -H 'Content-Type: application/json' \
  -d '{
    "token": "abcd1234567890abcd1234567890abcd1234567890abcd1234567890abcd1234",
    "new_password": "mi_nueva_password_segura"
  }'
```

**Respuesta exitosa:**
```json
{
  "success": true,
  "message": "Contraseña actualizada exitosamente"
}
```

---

## 🔍 Método 3: Manual (Paso a Paso)

### 1. Solicitar recuperación via API

```bash
curl -X POST http://localhost:8080/api/v1/auth/recover \
  -H 'Content-Type: application/json' \
  -d '{"email": "entidad@sgi.gov.co"}'
```

### 2. Obtener token desde la base de datos (solo para testing)

```bash
# Conectar a PostgreSQL
psql -U postgres -d sgi2_db

# Consultar token
SELECT email, reset_token, reset_token_expiry
FROM tickets_usuarios
WHERE email = 'entidad@sgi.gov.co';
```

**IMPORTANTE:** El valor en `reset_token` es el **HASH SHA-256** del token, no el token original.

### 3. Generar token de prueba manualmente

```bash
# Generar token aleatorio de 64 caracteres
TOKEN=$(openssl rand -hex 32)
echo "Token generado: $TOKEN"

# Calcular hash SHA-256
TOKEN_HASH=$(echo -n "$TOKEN" | shasum -a 256 | cut -d' ' -f1)
echo "Hash del token: $TOKEN_HASH"

# Insertar en BD (expira en 1 hora)
psql -U postgres -d sgi2_db -c "
UPDATE tickets_usuarios
SET
    reset_token = '$TOKEN_HASH',
    reset_token_expiry = NOW() + INTERVAL '1 hour'
WHERE email = 'entidad@sgi.gov.co';
"
```

### 4. Usar el token para resetear

```bash
curl -X POST http://localhost:8080/api/v1/auth/reset \
  -H 'Content-Type: application/json' \
  -d "{
    \"token\": \"$TOKEN\",
    \"new_password\": \"password_nueva_123\"
  }"
```

---

## 🛡️ Seguridad del Flujo

### ✅ Características de Seguridad

1. **Token de 32 bytes (64 caracteres hex)** - Alta entropía
2. **Hash SHA-256 en BD** - El token plano nunca se almacena
3. **Expiración de 1 hora** - Ventana de tiempo limitada
4. **Token de un solo uso** - Se elimina después de usarlo
5. **Validación de expiración** - Tokens vencidos son rechazados
6. **Cierre de sesiones activas** - Al cambiar password, todas las cookies se invalidan
7. **Respuesta genérica** - No revela si el email existe o no

### 🔐 Validaciones

```go
// En auth_handler.go:

1. Token debe existir en BD
2. Token no debe estar expirado (< 1 hora desde creación)
3. Hash del token recibido debe coincidir con el hash en BD
4. Nueva contraseña debe tener mínimo 8 caracteres
```

---

## 📧 Configuración de Email (Opcional)

Para probar el flujo completo con emails reales, configura en [.env](../.env):

```env
MAIL_SERVER_HOST=smtp.gmail.com
MAIL_SERVER_PORT=587
MAIL_SERVER_USER=tu_email@gmail.com
MAIL_SERVER_PASSWORD=tu_app_password
MAIL_SERVER_FROM=tu_email@gmail.com
```

**Nota:** Para Gmail, necesitas generar una "Contraseña de aplicación" en:
https://myaccount.google.com/apppasswords

---

## 🧪 Casos de Prueba

### ✅ Caso 1: Flujo exitoso
```bash
./seeds/test_reset_password_complete.sh
```

### ❌ Caso 2: Token inválido
```bash
curl -X POST http://localhost:8080/api/v1/auth/reset \
  -H 'Content-Type: application/json' \
  -d '{"token": "token_invalido", "new_password": "password123"}'
```
**Esperado:** `{"success": false, "error": "Token invalido o expirado"}`

### ❌ Caso 3: Token expirado
```bash
# Insertar token con expiración en el pasado
psql -U postgres -d sgi2_db -c "
UPDATE tickets_usuarios
SET reset_token_expiry = NOW() - INTERVAL '1 hour'
WHERE email = 'entidad@sgi.gov.co';
"

# Intentar usar el token
curl -X POST http://localhost:8080/api/v1/auth/reset \
  -H 'Content-Type: application/json' \
  -d '{"token": "tu_token", "new_password": "password123"}'
```
**Esperado:** `{"success": false, "error": "Token invalido o expirado"}`

### ❌ Caso 4: Contraseña muy corta
```bash
curl -X POST http://localhost:8080/api/v1/auth/reset \
  -H 'Content-Type: application/json' \
  -d '{"token": "token_valido", "new_password": "1234"}'
```
**Esperado:** Error de validación (mínimo 8 caracteres)

---

## 🔄 Restaurar Password Original

Después de probar, restaura la contraseña original (`password123`):

```bash
./seeds/run_seeders.sh
```

O manualmente en PostgreSQL:

```sql
-- Hash bcrypt de "password123"
UPDATE tickets_usuarios
SET password = '$2a$12$Dfgh.Uq8Ig.8Yal3YMRR9eGm5JYEjbSpSyIaXu5xM6hJ8Kk0bBgua',
    reset_token = NULL,
    reset_token_expiry = NULL
WHERE email = 'entidad@sgi.gov.co';
```

---

## 📋 Resumen de Scripts

| Script | Descripción |
|--------|-------------|
| [test_reset_password_complete.sh](test_reset_password_complete.sh) | ✅ Test automatizado completo (recomendado) |
| [test_password_recovery.sh](test_password_recovery.sh) | 📧 Solicita recuperación (requiere email configurado) |
| [run_seeders.sh](run_seeders.sh) | 🔄 Restaura passwords originales |

---

## 📖 Endpoints Relacionados

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| POST | `/api/v1/auth/recover` | Solicitar recuperación (envía email con token) |
| POST | `/api/v1/auth/reset` | Resetear contraseña con token |
| POST | `/api/v1/auth/login` | Login (para verificar nueva contraseña) |

---

## 💡 Tips

- 🔍 **Ver logs del servidor** para debug del envío de emails
- 📧 **Configura SMTP** para probar el flujo completo realista
- 🧪 **Usa el script automatizado** para testing rápido sin configurar email
- 🔄 **Re-ejecuta seeders** para resetear todo a estado inicial
