#!/bin/sh

set -e

echo "Running DB migrations..."

# Check if migrate binary exists
if [ ! -f "./migrate" ]; then
    echo "❌ migrate binary not found"
    exit 1
fi

# Check if migration folder exists
if [ ! -d "./migration" ]; then
    echo "❌ migration folder missing"
    exit 1
fi

echo "✅ migrate binary found"
echo "✅ migration folder found"

# Run migrations
./migrate -path ./migration -database "$DB_SOURCE" -verbose up

echo "Starting API..."

if [ "$PRODUCTION" = "true" ]; then
    echo "Running in PRODUCTION mode"
    exec /app/main
else
    echo "Running in DEVELOPMENT mode with Air"
    exec air -c .air.toml
fi
