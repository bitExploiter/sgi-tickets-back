#!/bin/bash
set -e

echo "🚀 Iniciando SGI Tickets Backend..."

run_query() {
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "$1" 2>/dev/null | xargs
}

# Ejecutar la aplicación (que corre migraciones automáticamente)
./main &
APP_PID=$!

# Esperar a que existan las tablas necesarias luego de las migraciones
MAX_RETRIES=40
RETRY=0

while [ "$RETRY" -lt "$MAX_RETRIES" ]; do
    if ! kill -0 "$APP_PID" 2>/dev/null; then
        echo "❌ La aplicación se detuvo antes de terminar migraciones"
        wait "$APP_PID"
        exit 1
    fi

    USUARIOS_TABLE=$(run_query "SELECT to_regclass('public.tickets_usuarios') IS NOT NULL;")
    DEPARTAMENTOS_TABLE=$(run_query "SELECT to_regclass('public.departamentos') IS NOT NULL;")
    MUNICIPIOS_TABLE=$(run_query "SELECT to_regclass('public.municipios') IS NOT NULL;")
    TIPOS_DOC_TABLE=$(run_query "SELECT to_regclass('public.tipo_documentos') IS NOT NULL;")
    REGIONALES_TABLE=$(run_query "SELECT to_regclass('public.regionales') IS NOT NULL;")

    if [ "$USUARIOS_TABLE" = "t" ] && [ "$DEPARTAMENTOS_TABLE" = "t" ] && [ "$MUNICIPIOS_TABLE" = "t" ] && [ "$TIPOS_DOC_TABLE" = "t" ] && [ "$REGIONALES_TABLE" = "t" ]; then
        break
    fi

    RETRY=$((RETRY + 1))
    sleep 2
done

if [ "$RETRY" -eq "$MAX_RETRIES" ]; then
    echo "❌ Timeout esperando finalización de migraciones"
    wait "$APP_PID"
    exit 1
fi

USER_COUNT=$(run_query "SELECT COUNT(*) FROM tickets_usuarios;")
DEPARTAMENTOS_COUNT=$(run_query "SELECT COUNT(*) FROM departamentos;")
MUNICIPIOS_COUNT=$(run_query "SELECT COUNT(*) FROM municipios;")
TIPOS_DOC_COUNT=$(run_query "SELECT COUNT(*) FROM tipo_documentos;")
REGIONALES_COUNT=$(run_query "SELECT COUNT(*) FROM regionales;")

if [ "$DEPARTAMENTOS_COUNT" = "0" ] || [ "$MUNICIPIOS_COUNT" = "0" ] || [ "$TIPOS_DOC_COUNT" = "0" ] || [ "$REGIONALES_COUNT" = "0" ]; then
    echo "📊 Catalogos vacíos o incompletos. Ejecutando seeders de catalogos..."

    /app/seeds/run_seeders.sh 02-departamentos.sql 03-municipios.sql 04-regionales.sql 05-tipos_documento.sql

    echo "✅ Seeders de catalogos ejecutados correctamente"
else
    echo "ℹ️  Catalogos ya cargados. Omitiendo seeders de catalogos."
fi

if [ "$USER_COUNT" = "0" ]; then
    echo "📊 Tabla de usuarios vacía. Ejecutando seeder de usuarios..."
    /app/seeds/run_seeders.sh 01-usuarios.sql

    echo "✅ Seeder de usuarios ejecutado correctamente"
else
    echo "ℹ️  Base de datos ya contiene datos ($USER_COUNT usuarios). Omitiendo seeder de usuarios."
fi

# Mantener el proceso principal
wait $APP_PID
