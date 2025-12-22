#!/bin/sh

set -e

if [ -z "$DSN" ]; then
  echo "Error: DSN environment variable is not set."
  exit 1
fi

echo "Waiting for PostgreSQL (via DSN)..."

until pg_isready -d "$DSN"; do
  echo "DB unavailable, waiting..."
  sleep 2
done

echo "DB ACCEPTING QUERIES, starting app!"

exec "$@"