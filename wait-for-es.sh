#!/bin/sh

set -e

ES_HOST="${1:-elasticsearch}"
ES_PORT="${2:-9200}"
shift 2

echo "⏳ Waiting for Elasticsearch at $ES_HOST:$ES_PORT..."

until curl -s "http://$ES_HOST:$ES_PORT" > /dev/null; do
  echo "Waiting for Elasticsearch..."
  sleep 2
done

echo "✅ Elasticsearch is up - executing command: $*"
exec "$@"
