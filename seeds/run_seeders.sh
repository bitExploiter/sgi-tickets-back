#!/bin/bash

# ==============================================
# Script para ejecutar seeders SQL
# ==============================================
# Lee credenciales del archivo .env y ejecuta
# todos los archivos .sql en orden numérico
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

# Verificar que existe .env
if [ ! -f ".env" ]; then
    echo -e "${RED}❌ Error: Archivo .env no encontrado${NC}"
    echo "Por favor crea un archivo .env con las credenciales de la BD"
    exit 1
fi

# Cargar variables del .env
export $(cat .env | grep -v '^#' | xargs)

# Verificar variables requeridas
if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_USER" ] || [ -z "$DB_NAME" ]; then
    echo -e "${RED}❌ Error: Faltan variables de entorno en .env${NC}"
    echo "Asegúrate de definir: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME"
    exit 1
fi

echo -e "🔗 Conectando a: ${GREEN}$DB_USER@$DB_HOST:$DB_PORT/$DB_NAME${NC}\n"

# Contar archivos SQL
SQL_FILES=$(ls seeds/*.sql 2>/dev/null | wc -l)

if [ "$SQL_FILES" -eq 0 ]; then
    echo -e "${YELLOW}⚠️  No hay archivos .sql en seeds/${NC}"
    exit 0
fi

# Ejecutar cada archivo SQL en orden
for file in seeds/*.sql; do
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
