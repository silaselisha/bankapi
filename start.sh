#!/bin/sh
set -e 

echo "run database migration"
/usr/bin/migrate -path /app/db/migrations -database "$DB_SOURCE" -verbose up

echo "start the app"
exec "$@"