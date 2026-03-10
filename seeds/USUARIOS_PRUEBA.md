# 👥 Usuarios de Prueba - Credenciales

## 🚀 Inicio Rápido

Para probar el login **sin 2FA** (recomendado):

```
Email: entidad@sgi.gov.co
Password: password123
```

## 📋 Lista Completa de Usuarios

| Rol | Email | Contraseña | 2FA |
|-----|-------|-----------|-----|
| 🔴 Admin | `admin@sgi.gov.co` | `password123` | ✅ Requerido |
| 🟡 Supervisor | `supervisor@sgi.gov.co` | `password123` | ✅ Requerido |
| 🔵 Agente | `agente@sgi.gov.co` | `password123` | ✅ Requerido |
| 🟣 Contratista | `contratista@sgi.gov.co` | `password123` | ✅ Requerido |
| 🟢 Entidad | `entidad@sgi.gov.co` | `password123` | ❌ No |

## 🔐 Sobre la Autenticación de Dos Factores (2FA)

### Roles SIN 2FA
- **entidad**: Login directo, ideal para pruebas de API

### Roles CON 2FA
- **admin, supervisor, agente, contratista**: Requieren configurar Google Authenticator

## 🔑 Flujo de Login con 2FA

### Primera Vez (Configurar 2FA)

1. **Hacer login** con email y contraseña
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H 'Content-Type: application/json' \
     -d '{"email":"admin@sgi.gov.co","password":"password123"}' \
     -c cookies.txt
   ```

2. El sistema responderá:
   ```json
   {
     "success": true,
     "require_2fa": true,
     "totp_enabled": false,
     "message": "Debes configurar autenticacion de dos factores"
   }
   ```

3. **Obtener código QR:**
   ```bash
   curl -X GET http://localhost:8080/api/v1/auth/2fa/setup \
     -b cookies.txt
   ```

   **Respuesta:**
   ```json
   {
     "success": true,
     "data": {
       "secret": "JBSWY3DPEHPK3PXP...",
       "qr_code": "data:image/png;base64,...",
       "issuer": "SGI Tickets",
       "account": "admin@sgi.gov.co"
     }
   }
   ```

4. **Escanear QR o usar el script automatizado:**
   ```bash
   ./seeds/setup_2fa_example.sh
   ```

   O manualmente:
   - Abre **Google Authenticator** en tu teléfono
   - Escanea el QR o ingresa el `secret` manualmente
   - Ingresa el código de 6 dígitos generado

5. **Confirmar código:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/2fa/setup \
     -H 'Content-Type: application/json' \
     -b cookies.txt \
     -d '{"code":"123456"}'
   ```

### Próximos Logins

1. **Login con email y password**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H 'Content-Type: application/json' \
     -d '{"email":"admin@sgi.gov.co","password":"password123"}' \
     -c cookies.txt
   ```

2. **Verificar código de Google Authenticator:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/2fa/verify \
     -H 'Content-Type: application/json' \
     -b cookies.txt \
     -d '{"code":"123456"}'
   ```

## 📡 Pruebas con Postman/cURL

### Login sin 2FA (Entidad)

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "entidad@sgi.gov.co",
    "password": "password123"
  }' \
  -c cookies.txt
```

### Login con 2FA (Admin)

**Paso 1: Login inicial**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@sgi.gov.co",
    "password": "password123"
  }' \
  -c cookies.txt
```

**Paso 2: Setup 2FA (primera vez)**
```bash
# Obtener QR
curl -X GET http://localhost:8080/api/v1/auth/2fa/setup \
  -b cookies.txt

# Confirmar código (reemplazar CODE con código de Google Authenticator)
curl -X POST http://localhost:8080/api/v1/auth/2fa/setup \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"code": "123456"}'
```

**Paso 3: Verificar 2FA (logins posteriores)**
```bash
curl -X POST http://localhost:8080/api/v1/auth/2fa/verify \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{"code": "123456"}'
```

## 🗑️ Limpiar Usuarios de Prueba

```sql
DELETE FROM tickets_usuarios WHERE origen = 'seeder';
```

## 🔄 Re-ejecutar Seeder

```bash
# Opción 1: Script automático
./seeds/run_seeders.sh

# Opción 2: Manual
psql -U postgres -d sgi2_db -f seeds/01-usuarios.sql
```

## ⚠️ Notas Importantes

- 🔒 **Contraseñas hasheadas** con bcrypt (cost factor 12)
- 🍪 **Cookies de sesión** válidas por 7 días
- 📧 **Solo desarrollo**: No usar estos usuarios en producción
- 🔑 **Password universal** para pruebas: `password123`
