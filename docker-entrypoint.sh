#!/bin/bash
set -e

echo "🚀 Iniciando SGI Tickets Backend..."

# Ejecutar la aplicación (que corre migraciones automáticamente)
./main &
APP_PID=$!

# Esperar a que la app arranque y las migraciones se ejecuten
sleep 3

# Verificar si hay usuarios en la BD
USER_COUNT=$(PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT COUNT(*) FROM tickets_usuarios;" 2>/dev/null | xargs)

if [ "$USER_COUNT" = "0" ]; then
    echo "📊 Base de datos vacía. Ejecutando seeders..."

    # Ejecutar seeders SQL si existen
    for seed_file in /app/seeds/*.sql; do
        if [ -f "$seed_file" ]; then
            echo "   🌱 Ejecutando $(basename $seed_file)..."
            PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$seed_file" -q
        fi
    done

    echo "✅ Seeders ejecutados correctamente"
else
    echo "ℹ️  Base de datos ya contiene datos ($USER_COUNT usuarios). Omitiendo seeders."
fi

# Mantener el proceso principal
wait $APP_PID
