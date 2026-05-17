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

# Step 2: Remove SQLC configuration for this module
if [ -f "sqlc.yaml" ]; then
    echo -e "${GREEN}Removing SQLC configuration for module: $MODULE_NAME${NC}"
    
    # Create a temporary file without the module's SQLC config
    awk -v module="$MODULE_NAME" '
    BEGIN { skip=0; printed_header=0 }
    /^sql:/ { print; next }
    /^  - engine:/ { 
        # Store current line
        current=$0
        # Look ahead to see if this is our module
        getline next_line
        if (next_line ~ "schema: ./internal/" module "/migration") {
            skip=1
            # Skip until we find the next engine or end of file
            while (getline > 0) {
                if ($0 ~ /^  - engine:/) {
                    # Found next engine, reprocess this line
                    skip=0
                    print current
                    print $0
                    break
                }
            }
            if (skip) {
                # End of file reached
                break
            }
        } else {
            # Not our module, print both lines
            print current
            print next_line
        }
        next
    }
    {
        if (!skip) print
    }' sqlc.yaml > sqlc.yaml.tmp
    
    # Replace original file with cleaned version
    mv sqlc.yaml.tmp sqlc.yaml
    
    # Clean up empty sqlc.yaml or remove if only the header remains
    if [ ! -s sqlc.yaml ] || [ "$(wc -l < sqlc.yaml)" -le 1 ]; then
        echo -e "${YELLOW}No modules left in sqlc.yaml, removing file${NC}"
        rm sqlc.yaml
    fi
    
    echo -e "${GREEN}SQLC configuration removed for module: $MODULE_NAME${NC}"
else
    echo -e "${YELLOW}sqlc.yaml not found, skipping SQLC cleanup${NC}"
fi

# Step 3: Remove any references in router initialization
echo -e "${GREEN}Checking for references in code...${NC}"

# Check main.go for module references
if [ -f "cmd/api/main.go" ]; then
    if grep -q "internal/$MODULE_NAME" cmd/api/main.go; then
        echo -e "${YELLOW}Found references to '$MODULE_NAME' in cmd/api/main.go${NC}"
        echo "You may need to manually remove these imports and initialization code."
    fi
fi

# Check router.go if it exists
if [ -f "internal/router/router.go" ]; then
    if grep -q "internal/$MODULE_NAME" internal/router/router.go; then
        echo -e "${YELLOW}Found references to '$MODULE_NAME' in internal/router/router.go${NC}"
        echo "You may need to manually remove these route registrations."
    fi
fi

# Step 4: Run go mod tidy to clean up dependencies
echo -e "${GREEN}Running go mod tidy to clean up dependencies...${NC}"
go mod tidy
if [ $? -eq 0 ]; then
    echo -e "${GREEN}Dependencies tidied successfully${NC}"
else
    echo -e "${YELLOW}go mod tidy had issues (this is normal if no code changes were needed)${NC}"
fi

# Step 5: Remove migration files from the database (optional)
echo ""
read -p "Do you want to also drop the migration table for this module from the database? (y/N): " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${GREEN}Dropping migration table for module: $MODULE_NAME${NC}"
    echo "Please run the appropriate migration down command manually."
    echo "Example: migrate -path internal/$MODULE_NAME/migration -database \"your_database_url\" down"
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
echo "  3. Run database migrations down if needed"
echo "  4. Run 'make build' to verify everything works"