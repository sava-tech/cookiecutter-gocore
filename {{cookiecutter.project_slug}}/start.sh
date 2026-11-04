#!/bin/sh

set -e

# Live reload needs the Go toolchain, which only the dev image has
if [ "$PRODUCTION" != "true" ] && command -v go >/dev/null 2>&1; then
    DEV_MODE=true
    MIGRATE_BIN=migrate
    # Gather each module's migrations from the mounted source, so new ones
    # are picked up on the next container start
    MIGRATION_DIR=/tmp/migration
    rm -rf "$MIGRATION_DIR"
    mkdir -p "$MIGRATION_DIR"
    cp internal/*/migration/* "$MIGRATION_DIR"/
else
    DEV_MODE=false
    MIGRATE_BIN=./migrate
    MIGRATION_DIR=./migration
fi

echo "Running DB migrations..."

# Check if migrate binary exists
if ! command -v "$MIGRATE_BIN" >/dev/null 2>&1; then
    echo "❌ migrate binary not found"
    exit 1
fi

# Check if migration folder exists
if [ ! -d "$MIGRATION_DIR" ]; then
    echo "❌ migration folder missing"
    exit 1
fi

echo "✅ migrate binary found"
echo "✅ migration folder found"

# Run migrations
"$MIGRATE_BIN" -path "$MIGRATION_DIR" -database "$DB_SOURCE" -verbose up

echo "Starting API..."

if [ "$DEV_MODE" = "true" ]; then
    echo "Running in DEVELOPMENT mode with Air"
    echo "✅ server: http://localhost:8080"
    echo "✅ API doc: http://localhost:8080/swagger/index.html"
    echo "✅ mailpit: http://localhost:8025"
    # The source mount hides the docs generated at image build time
    swag init -g cmd/api/main.go --parseDependency --parseInternal --output ./docs
    # Hand the generated files back to whoever owns the mounted source
    chown -R "$(stat -c '%u:%g' .)" docs
    # Build outside the mount so no root-owned files land in the source folder
    exec air -c .air.toml --tmp_dir /tmp/air \
        --build.cmd "go build -o /tmp/air/main ./cmd/api" --build.bin /tmp/air/main
else
    echo "Running in PRODUCTION mode"
    exec /app/main
fi
