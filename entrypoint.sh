#!/bin/sh
set -e

echo "🔄 Запуск ожидания PostgreSQL..."
/search/wait-for-db.sh

echo "🔄 Запуск ожидания Elasticsearch..."
/search/wait-for-es.sh

echo "🚀 Запуск приложения..."
exec /search/search_service "$@"
