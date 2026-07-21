# Simple Makefile for a Go project
SHELL := /bin/bash
PROJECTNAME := "blog"

help:
	@echo "Available make commands:"
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'
	@echo
	@echo "Usage: make <target> [VARIABLE=value]"
	@echo "Example: make module name=users"


# Build the application
all: build test ## Build and test the application

build: ## Build the Go binary
	@echo "Building..."
	
	@go build -o main cmd/api/main.go

deps: ## Download and tidy Go dependencies
	@echo "Downloading Go modules..."
	go mod tidy
	@echo "Go modules ready"

install-dependencies: ## Install all required dependencies for the project
# 	sudo apt install ca-certificates
	@$(MAKE) deps
	@echo "  → Installing development tools..."
	@# Install golangci-lint for linting
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@# Install mockgen for mocking (optional)
	@which mockgen > /dev/null || (echo "Installing mockgen..." && go install github.com/golang/mock/mockgen@latest)
#	gomigrate
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/hibiken/asynq@latest
	go get github.com/gorilla/sessions
#	for Gin
	go get github.com/swaggo/gin-swagger
	go get github.com/swaggo/files
# 	Install Goth
	@echo "Installing Goth and providers..."
	@go get github.com/markbates/goth
	@go get github.com/markbates/goth/providers/google
	@go get github.com/markbates/goth/providers/github
	@go mod tidy
	@echo "Goth dependencies installed!"
#	firebase oauth
	go get firebase.google.com/go/v4
	go get google.golang.org/api/option
	@# Install mockgen for mocking (optional)
	@which mockgen > /dev/null || (echo "Installing mockgen..." && go install github.com/golang/mock/mockgen@latest)
	@echo "All dependencies installed!"



# Run the application
run: ## Run the Go application
	@go run cmd/api/main.go

docker-run: ## Build and start Docker containers
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down: ## Stop Docker containers
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test: ## Run all tests
	@echo "Testing..."
	@go test ./... -v

# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi


create-module: ## Create a new module. Usage: make create-module name=<module_name>
	@if [ -z "$(name)" ]; then \
		echo "Please provide a module name. Example: make create-module name=users"; \
		exit 1; \
	fi
	@bash scripts/create_sqlc_module.sh "$(name)" "$(PROJECTNAME)"

delete-module: ## Delete a module and its SQLC config. Usage: make delete-module name=<module_name>
	@if [ -z "$(name)" ]; then \
		echo "Please provide a module name. Example: make delete-module name=users"; \
		exit 1; \
	fi
	@python3 scripts/delete_module.py "$(name)"


# to tidy the code base by installing dependencies
tidy:
	go mod tidy

migration: ## Create a new migration. Usage: make migration module=<module> name=<migration_name>
	@if [ -z "$(module)" ]; then \
		echo "Please provide a module. Example: make migration module=users name=add_profile_table"; \
		exit 1; \
	fi
	@if [ -z "$(name)" ]; then \
		echo "Please provide a migration name. Example: make migration module=users name=add_profile_table"; \
		exit 1; \
	fi
	@if [ ! -d "internal/$(module)" ]; then \
		echo "❌ Module '$(module)' does not exist in internal/"; \
		exit 1; \
	fi

	@BASE_DIR=internal/$(module); \
	MIGRATION_DIR=$$BASE_DIR/migration; \
	mkdir -p $$MIGRATION_DIR; \
	TIMESTAMP=$$(date +"%Y%m%d%H%M%S"); \
	COUNT=$$(ls -1 $$MIGRATION_DIR/*.up.sql 2>/dev/null | wc -l | tr -d ' '); \
	SEQ=$$(printf "%06d" $$((COUNT + 1))); \
	FILENAME="$${TIMESTAMP}_$${SEQ}_$(name)"; \
	echo "🚀 Creating migration $$FILENAME for module '$(module)'"; \
	touch $$MIGRATION_DIR/$$FILENAME.up.sql; \
	touch $$MIGRATION_DIR/$$FILENAME.down.sql; \
	echo "✅ Created:"; \
	echo "   - $$MIGRATION_DIR/$$FILENAME.up.sql"; \
	echo "   - $$MIGRATION_DIR/$$FILENAME.down.sql";

sqlc:
	sqlc generate

db_docs:
	dbdocs build internal/doc/db.dbml

db_schema:
	dbml2sql internal/doc/db.dbml -o internal/doc/db.sql

swagger-doc:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/api/main.go --parseDependency --parseInternal --output ./docs
	@echo "Swagger documentation generated at ./docs"

mock:
	mockgen -package mock_db -destination internal/users/mock/users.go github.com/michaelassa01/blog/internal/users/domain Repository

.PHONY: all build run test clean watch docker-run docker-down itest install-dependencies test tidy swagger-doc createmigration module sqlc migration db_docs db_schema mock
