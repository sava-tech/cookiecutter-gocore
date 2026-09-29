#!/bin/bash

# Usage: ./scripts/delete_sqlc_module.sh <module_name>
#        (or: make delete-module name=<module_name>)
# Set FORCE=1 to skip the confirmation prompt.
MODULE_NAME="$1"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if module name is provided
if [ -z "$MODULE_NAME" ]; then
    echo -e "${RED}ERROR: Module name is required${NC}"
    echo "Usage: $0 <module_name>"
    echo "Example: $0 posts"
    exit 1
fi

# Only plain module names: this path is passed to rm -rf
if [[ ! "$MODULE_NAME" =~ ^[a-z][a-z0-9_]*$ ]]; then
    echo -e "${RED}ERROR: Invalid module name '$MODULE_NAME' (use lowercase letters, digits and underscores)${NC}"
    exit 1
fi

# Core packages under internal/ that aren't removable modules
case "$MODULE_NAME" in
    server|shared|database|doc)
        echo -e "${RED}ERROR: '$MODULE_NAME' is a core package, not a module${NC}"
        exit 1
        ;;
esac

BASE_DIR="internal/$MODULE_NAME"
SCHEMA_PATH="./internal/$MODULE_NAME/migration"

# Check if module exists
if [ ! -d "$BASE_DIR" ]; then
    echo -e "${RED}ERROR: Module '$MODULE_NAME' does not exist in internal/${NC}"
    exit 1
fi

echo -e "${YELLOW}About to delete module: $MODULE_NAME${NC}"
echo "This will delete:"
echo "  - Directory: $BASE_DIR (including its migrations)"
echo "  - SQLC configuration for $MODULE_NAME in sqlc.yaml"
echo ""
if [ "$FORCE" != "1" ]; then
    read -p "Are you sure you want to continue? (y/N): " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}Deletion cancelled${NC}"
        exit 0
    fi
fi

# Step 1: Remove the module directory
echo -e "${GREEN}Removing module directory...${NC}"
if rm -rf "$BASE_DIR"; then
    echo -e "${GREEN}Module directory deleted: $BASE_DIR${NC}"
else
    echo -e "${RED}Failed to delete module directory${NC}"
    exit 1
fi

# Step 2: Remove the module's entry from sqlc.yaml
#
# Each module is one "  - engine:" list item under "sql:", optionally
# preceded by "  # comment" lines. Buffer each item and drop the one whose
# schema points at this module's migration folder; everything else,
# including blank lines and comments, is kept as-is.
if [ -f "sqlc.yaml" ]; then
    echo -e "${GREEN}Removing SQLC configuration for module: $MODULE_NAME${NC}"

    TMP_FILE=$(mktemp)
    awk -v schema="$SCHEMA_PATH" '
        function flush() {
            if (!drop) printf "%s", block
            block = ""; drop = 0
        }
        /^  #/ { pending = pending $0 "\n"; next }
        /^  - engine:/ {
            flush()
            block = pending $0 "\n"; pending = ""; in_block = 1
            next
        }
        /^[[:space:]]*$/ {
            if (pending != "") pending = pending $0 "\n"
            else if (in_block) block = block $0 "\n"
            else print
            next
        }
        {
            if (pending != "") {
                if (in_block) block = block pending; else printf "%s", pending
                pending = ""
            }
            if (!in_block) { print; next }
            block = block $0 "\n"
            if ($0 ~ /^    schema:/ && (index($0, "\"" schema "\"") || index($0, " " schema))) drop = 1
        }
        END { flush(); printf "%s", pending }
    ' sqlc.yaml > "$TMP_FILE"

    if cmp -s sqlc.yaml "$TMP_FILE"; then
        echo -e "${YELLOW}No SQLC entry found for $MODULE_NAME, sqlc.yaml unchanged${NC}"
        rm -f "$TMP_FILE"
    elif ! grep -q "^  - engine:" "$TMP_FILE"; then
        # sqlc rejects a config with no sql entries; create-module recreates it
        rm -f sqlc.yaml "$TMP_FILE"
        echo -e "${GREEN}Removed sqlc.yaml (no modules left)${NC}"
    else
        mv "$TMP_FILE" sqlc.yaml
        echo -e "${GREEN}SQLC configuration removed for module: $MODULE_NAME${NC}"
    fi
else
    echo -e "${YELLOW}sqlc.yaml not found, skipping SQLC cleanup${NC}"
fi

# Step 3: Warn about code that still references the module
echo -e "${GREEN}Checking for references in code...${NC}"
REFS=$(grep -rln --include='*.go' "internal/$MODULE_NAME\"\|internal/$MODULE_NAME/" cmd internal 2>/dev/null)
if [ -n "$REFS" ]; then
    echo -e "${YELLOW}These files still import '$MODULE_NAME'; remove those imports and route registrations:${NC}"
    echo "$REFS" | sed 's/^/  - /'
fi

# Step 4: Run go mod tidy
echo -e "${GREEN}Running go mod tidy to clean up dependencies...${NC}"
if go mod tidy; then
    echo -e "${GREEN}Dependencies tidied successfully${NC}"
else
    echo -e "${YELLOW}go mod tidy reported issues (see output above)${NC}"
fi

echo ""
echo -e "${GREEN}Module '$MODULE_NAME' has been deleted!${NC}"
echo ""
echo "Next steps:"
echo "  1. Remove any remaining imports/usages listed above"
echo "     (usually internal/server/router.go and internal/server/server.go)"
echo "  2. Drop the module's tables with a new migration if they exist in your database"
echo "  3. Run 'make build' to verify everything works"
