#!/bin/bash

# ==============================================
# Script para ejecutar seeders SQL
# ==============================================
# Usa variables de entorno existentes o lee credenciales
# desde .env y ejecuta seeders SQL.
# Permite ejecutar archivos específicos si se pasan
# como argumentos.
# ==============================================

set -e  # Exit on error

# Colores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}======================================"
echo -e "🌱 Ejecutando Seeders de Base de Datos"
echo -e "======================================${NC}\n"

# Cargar variables desde .env solo si no existen en el entorno
if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_USER" ] || [ -z "$DB_NAME" ]; then
    if [ -f ".env" ]; then
        export $(cat .env | grep -v '^#' | xargs)
    fi
fi

# Verificar variables requeridas
if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_USER" ] || [ -z "$DB_NAME" ]; then
    echo -e "${RED}❌ Error: Faltan variables de entorno de BD${NC}"
    echo "Define: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME (o usa .env)"
    exit 1
fi

echo -e "🔗 Conectando a: ${GREEN}$DB_USER@$DB_HOST:$DB_PORT/$DB_NAME${NC}\n"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# Si se pasan argumentos, se ejecutan solo esos archivos.
if [ "$#" -gt 0 ]; then
    FILES_TO_RUN=()
    for arg in "$@"; do
        if [[ "$arg" = /* ]]; then
            FILES_TO_RUN+=("$arg")
        else
            FILES_TO_RUN+=("$SCRIPT_DIR/$arg")
        fi
    done
else
    FILES_TO_RUN=("$SCRIPT_DIR"/*.sql)
fi

# Contar archivos SQL seleccionados
SQL_FILES=0
for file in "${FILES_TO_RUN[@]}"; do
    if [ -f "$file" ]; then
        SQL_FILES=$((SQL_FILES + 1))
    fi
done

if [ "$SQL_FILES" -eq 0 ]; then
    echo -e "${YELLOW}⚠️  No hay archivos .sql para ejecutar${NC}"
    exit 0
fi

# Ejecutar cada archivo SQL en orden
for file in "${FILES_TO_RUN[@]}"; do
    if [ -f "$file" ]; then
        filename=$(basename "$file")
        echo -e "📄 Ejecutando: ${GREEN}$filename${NC}"

        # Ejecutar SQL con password desde variable de entorno
        PGPASSWORD=$DB_PASSWORD psql \
            -h "$DB_HOST" \
            -p "$DB_PORT" \
            -U "$DB_USER" \
            -d "$DB_NAME" \
            -f "$file" \
            -v ON_ERROR_STOP=1 \
            -q  # Quiet mode

        if [ $? -eq 0 ]; then
            echo -e "   ${GREEN}✅ Completado${NC}\n"
        else
            echo -e "   ${RED}❌ Error ejecutando $filename${NC}\n"
            exit 1
        fi
    fi
done

echo -e "${GREEN}======================================"
echo -e "✅ Todos los seeders ejecutados correctamente"
echo -e "======================================${NC}\n"

# Mostrar usuarios creados
echo -e "${YELLOW}👥 Usuarios de prueba creados:${NC}"
PGPASSWORD=$DB_PASSWORD psql \
    -h "$DB_HOST" \
    -p "$DB_PORT" \
    -U "$DB_USER" \
    -d "$DB_NAME" \
    -c "SELECT id, nombres || ' ' || apellidos AS nombre_completo, email, rol, activo FROM tickets_usuarios WHERE origen = 'seeder' ORDER BY id;" \
    2>/dev/null

echo -e "\n${GREEN}🔑 Contraseña para todos: password123${NC}"
echo -e "${YELLOW}💡 Tip: Usa entidad@sgi.gov.co para pruebas sin 2FA${NC}\n"
