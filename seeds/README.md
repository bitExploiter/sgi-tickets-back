# Seeders - Datos de Prueba

## Descripción

Archivos SQL para poblar la base de datos con datos de prueba.

## Usuarios de Prueba

El archivo `01-usuarios.sql` crea 5 usuarios de prueba (uno por cada rol):

| Email | Contraseña | Rol | Requiere 2FA | Descripción |
|-------|-----------|-----|-------------|-------------|
| `entidad@sgi.gov.co` | `password123` | entidad | ❌ No | **Recomendado para pruebas rápidas** |
| `admin@sgi.gov.co` | `password123` | admin | ✅ Sí | Administrador del sistema |
| `supervisor@sgi.gov.co` | `password123` | supervisor | ✅ Sí | Supervisor de tickets |
| `agente@sgi.gov.co` | `password123` | agente | ✅ Sí | Agente de atención |
| `contratista@sgi.gov.co` | `password123` | contratista | ✅ Sí | Contratista externo |

## Cómo Ejecutar

### Opción 1: Usando psql directamente

```bash
psql -U tu_usuario -d tu_base_datos -f seeds/01-usuarios.sql
```

### Opción 2: Desde el archivo .env

```bash
# Leer credenciales del .env
source .env

# Ejecutar seeder
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f seeds/01-usuarios.sql
```

### Opción 3: Script automático

```bash
# Crear script ejecutable
chmod +x seeds/run_seeders.sh

# Ejecutar todos los seeders
./seeds/run_seeders.sh
```

## Verificar Usuarios Creados

```sql
SELECT id, nombres, apellidos, email, rol, activo
FROM tickets_usuarios
WHERE origen = 'seeder';
```

## Resetear Usuarios de Prueba

### Opción 1: Script automático (Recomendado)

Resetea todos los usuarios de prueba a `password123` y limpia tokens:

```bash
# Con .env
export $(cat .env | grep -v '^#' | xargs)
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f seeds/reset_passwords.sql

# O directamente
psql -U postgres -d sgi2_db -f seeds/reset_passwords.sql
```

### Opción 2: Resetear contraseña de un usuario específico

```sql
UPDATE tickets_usuarios
SET password = '$2a$12$Dfgh.Uq8Ig.8Yal3YMRR9eGm5JYEjbSpSyIaXu5xM6hJ8Kk0bBgua', -- password123
    reset_token = NULL,
    reset_token_expiry = NULL,
    totp_token = ''
WHERE email = 'entidad@sgi.gov.co';
```

### Opción 3: Eliminar y recrear todos los usuarios

```sql
-- CUIDADO: Esto eliminará permanentemente los usuarios
DELETE FROM tickets_usuarios WHERE origen = 'seeder';
```

Luego ejecuta `./seeds/run_seeders.sh` para recrearlos.

## Notas Importantes

- ✅ **Solo crea si NO existe** - No sobrescribe usuarios existentes (preserva contraseñas modificadas)
- ⚠️ **Contraseña hasheada con bcrypt (cost 12)** - el hash es diferente cada vez que se genera
- 🔐 **Usuarios con 2FA** requieren configurar Google Authenticator después del primer login
- ✅ **Usuario "entidad"** no requiere 2FA - ideal para pruebas de API
- 📧 **Todos los emails** terminan en `@sgi.gov.co`
- 🔑 **Contraseña inicial:** `password123` (solo para desarrollo)

## Generar Hash de Nueva Contraseña

Si necesitas crear usuarios con diferentes contraseñas:

```bash
go run scripts/hash_password.go "tu_contraseña_aqui"
```

Esto generará un hash bcrypt con un ejemplo de SQL INSERT completo.

## Testing de Funcionalidades

### 🔐 Autenticación de Dos Factores (2FA)

Ver guía completa: [USUARIOS_PRUEBA.md](USUARIOS_PRUEBA.md#-flujo-de-login-con-2fa)

```bash
# Script automatizado para configurar 2FA
./seeds/setup_2fa_example.sh
```

### 🔑 Recuperación de Contraseña

Ver guía completa: [RECUPERAR_PASSWORD.md](RECUPERAR_PASSWORD.md)

```bash
# Test automatizado del flujo completo
./seeds/test_reset_password_complete.sh

# O flujo manual con email real
./seeds/test_password_recovery.sh
```
