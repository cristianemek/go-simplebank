#!/bin/sh

set -e

echo "Running database migrations..."
. /app/app.env.example
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "Starting the application..."
exec "$@"