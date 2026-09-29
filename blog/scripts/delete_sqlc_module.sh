#!/bin/bash

# Usage: ./scripts/delete_module.sh <module_name> <project_name>
MODULE_NAME="$1"
PROJECT_NAME="$2"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if module name is provided
if [ -z "$MODULE_NAME" ]; then
    echo -e "${RED}ERROR: Module name is required${NC}"
    echo "Usage: $0 <module_name> [project_name]"
    echo "Example: $0 users"
    exit 1
fi

BASE_DIR="internal/$MODULE_NAME"

# Check if module exists
if [ ! -d "$BASE_DIR" ]; then
    echo -e "${RED}ERROR: Module '$MODULE_NAME' does not exist in internal/${NC}"
    exit 1
fi

echo -e "${YELLOW}About to delete module: $MODULE_NAME${NC}"
echo "This will delete:"
echo "  - Directory: $BASE_DIR"
echo "  - SQLC configuration for $MODULE_NAME"
echo ""
read -p "Are you sure you want to continue? (y/N): " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}Deletion cancelled${NC}"
    exit 0
fi

# Step 1: Remove the module directory
echo -e "${GREEN}Removing module directory...${NC}"
rm -rf "$BASE_DIR"
if [ $? -eq 0 ]; then
    echo -e "${GREEN}Module directory deleted: $BASE_DIR${NC}"
else
    echo -e "${RED}Failed to delete module directory${NC}"
    exit 1
fi

# Step 2: Remove SQLC configuration for this module using sed
if [ -f "sqlc.yaml" ]; then
    echo -e "${GREEN}Removing SQLC configuration for module: $MODULE_NAME${NC}"
    
    # Create a backup
    cp sqlc.yaml sqlc.yaml.bak
    
    # Use sed to remove the module section
    # This removes from "- engine:" until the next "- engine:" or end of file
    sed -i "/^  - engine:/,/^  - engine:/{
        /^  - engine:/{
            :start
            N
            /schema: .\/internal\/$MODULE_NAME\/migration/{
                :delete
                N
                /^  - engine:/!b delete
                d
            }
            /^  - engine:/!b start
        }
    }" sqlc.yaml
    
    # Clean up empty lines
    sed -i '/^[[:space:]]*$/d' sqlc.yaml
    
    # Check if sqlc.yaml is now empty or only has header
    if [ ! -s sqlc.yaml ] || [ "$(grep -c "engine:" sqlc.yaml)" -eq 0 ]; then
        # Remove file if no engines left
        rm sqlc.yaml
        echo -e "${GREEN}Removed empty sqlc.yaml${NC}"
    else
        echo -e "${GREEN}SQLC configuration removed for module: $MODULE_NAME${NC}"
    fi
    
    rm -f sqlc.yaml.bak
else
    echo -e "${YELLOW}sqlc.yaml not found, skipping SQLC cleanup${NC}"
fi

# Step 3: Check for references
echo -e "${GREEN}Checking for references in code...${NC}"

# Check main.go
if [ -f "cmd/api/main.go" ]; then
    if grep -q "internal/$MODULE_NAME" cmd/api/main.go; then
        echo -e "${YELLOW}Found references to '$MODULE_NAME' in cmd/api/main.go${NC}"
        echo "You may need to manually remove these imports and initialization code."
    fi
fi

# Check router.go
if [ -f "internal/router/router.go" ]; then
    if grep -q "internal/$MODULE_NAME" internal/router/router.go; then
        echo -e "${YELLOW}Found references to '$MODULE_NAME' in internal/router/router.go${NC}"
        echo "You may need to manually remove these route registrations."
    fi
fi

# Step 4: Run go mod tidy
echo -e "${GREEN}Running go mod tidy to clean up dependencies...${NC}"
go mod tidy
if [ $? -eq 0 ]; then
    echo -e "${GREEN}Dependencies tidied successfully${NC}"
else
    echo -e "${YELLOW}go mod tidy had issues${NC}"
fi

echo ""
echo -e "${GREEN}Module '$MODULE_NAME' has been successfully deleted!${NC}"
echo ""
echo "Summary:"
echo "  - Deleted directory: internal/$MODULE_NAME"
echo "  - Removed from sqlc.yaml"
echo "  - Ran go mod tidy"
echo ""
echo "Next steps (manual):"
echo "  1. Remove any imports from cmd/api/main.go"
echo "  2. Remove route registrations from internal/router/router.go"
echo "  3. Run 'make build' to verify everything works"