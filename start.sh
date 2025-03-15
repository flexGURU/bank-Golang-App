#!/bin/sh

# script exits immediatley if the script returns a non zero status
set -e

echo "run db migration"
/app/migrate -path /app/migration -database "$DSN" -verbose up

echo "start the app"
exec "$@"