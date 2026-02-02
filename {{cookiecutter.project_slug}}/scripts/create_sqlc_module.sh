#!/bin/bash
MODULE_NAME="$1"
PROJECT_NAME="$2"

BASE_DIR="internal/$MODULE_NAME"

# Create directories
mkdir -p "$BASE_DIR/migration" "$BASE_DIR/query" "$BASE_DIR/models" "$BASE_DIR/domain"

# Create Go files
for f in router.go repository.go service.go handler.go pg_repo.go dto.go; do
  echo "package $MODULE_NAME" > "$BASE_DIR/$f"
done

# Create initial migration (assuming migrate CLI installed)
VERSION=$(printf "%06d" $(ls -1 "$BASE_DIR/migration" 2>/dev/null | wc -l))
MIGRATION_NAME="${VERSION}_init_schema"
migrate create -ext sql -dir "$BASE_DIR/migration" "$MIGRATION_NAME"

# Ensure sqlc.yaml exists
if [ ! -f sqlc.yaml ] || [ ! -s sqlc.yaml ]; then
cat > sqlc.yaml <<EOL
version: "2"
cloud:
  project: "$PROJECT_NAME"
sql:
EOL
  echo "✅ Created base sqlc.yaml"
fi

# Append module to sqlc.yaml if not already present
if ! grep -q "schema: ./internal/$MODULE_NAME/migration" sqlc.yaml; then
cat >> sqlc.yaml <<EOL
  - engine: "postgresql"
    schema: "./internal/$MODULE_NAME/migration"
    queries: "./internal/$MODULE_NAME/query"

    gen:
      go:
        package: "$MODULE_NAME"
        out: "./internal/$MODULE_NAME/models"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: true
        emit_exact_table_names: false
        emit_empty_slices: true
        emit_pointers_for_null_types: true
        overrides:
          - db_type: timestamptz
            go_type: time.Time
EOL
  echo "✅ Added SQLC config for module $MODULE_NAME"
else
  echo "⚠️ SQLC config for module $MODULE_NAME already exists — skipping"
fi
