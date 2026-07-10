# ![{{ cookiecutter.project_name }} Logo](https://via.placeholder.com/20) {{ cookiecutter.project_name }}

> A modular Go backend boilerplate built with **Domain-Driven Design (DDD)**, **Gin**, **SQLC**, and **Postgres** — ready for multi-module projects.

[![Go](https://img.shields.io/badge/Go-1.21-blue)](https://golang.org)  
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)  
[![Docker](https://img.shields.io/badge/Docker-Enabled-blue)](https://www.docker.com/)

---

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Requirements](#requirements)
- [Getting Started](#getting-started)
- [Makefile Commands](#makefile-commands)
- [Creating Modules](#creating-modules)
- [Running Migrations](#running-migrations)
- [Swagger Documentation](#swagger-documentation)
- [Testing](#testing)
- [License](#license)

---

## Features

- Gin HTTP server with modular routing
- DDD-style structure:
  - **Domain** → business entities
  - **Repository** → SQLC + Postgres
  - **Service** → business logic
  - **Handler** → HTTP handlers
- Multi-module support (e.g., `users`, `posts`)
- SQLC for type-safe SQL queries
- Dockerized Postgres and API
- Makefile helpers for migrations, modules, testing, and live reload
- Swagger API documentation

---

## Project Structure

{{ cookiecutter.project_name }}/
├── cmd/api/main.go # Entry point
├── internal/
│ ├── database/ # DB connection
│ ├── users/ # Users module (DDD)
│ │ ├── domain.go
│ │ ├── handler.go
│ │ ├── repository.go
│ │ ├── router.go
│ │ ├── service.go
│ │ ├── query/ # SQLC queries
│ │ └── migration/ # migrations
│ ├── posts/ # Another module
│ └── ... # Additional modules
├── scripts/ # helper scripts
├── migration/ # global migrations (optional)
├── Makefile
├── sqlc.yaml
├── docker-compose.yml
├── go.mod
└── README.md

---

## Requirements

- Go 1.21+
- Postgres 15+
- Docker & Docker Compose
- Make
- `air` (optional, for live reload)

---

## Getting Started

1. **Clone the repository**

```bash
git clone https://github.com/<your-username>/{{ cookiecutter.project_name }}.git
cd {{ cookiecutter.project_name }}

Install dependencies

make install-dependencies
go mod tidy

Start Docker Postgres container

make docker-run

Run the API

make run

live reload:

make watch

Makefile Commands

Run make help to see a full list of available commands:

make help


Some examples:

make build                # Build Go binary
make run                  # Run API
make test                 # Run all tests
make tidy                 # Go mod tidy
make docker-run           # Start Docker containers
make docker-down          # Stop Docker containers
make module name=users    # Create a new module
make migration module=users name=add_profile_table  # Create a migration
make sqlc                 # Generate SQLC code
make swagger-doc          # Generate Swagger docs

Creating Modules

Modules follow the DDD structure: domain → repository → service → handler → router → migration → queries.

make module name=users


This will automatically:

Create folder structure for the module

Create empty files: handler.go, repository.go, service.go, router.go

Initialize a first migration in migration/

Update sqlc.yaml with module paths

Running Migrations

Create a new migration:

make migration module=users name=add_profile_table


Run migrations:

./migrate -path ./migration -database "$DB_SOURCE" -verbose up


Example DB_SOURCE:

postgres://user:password@localhost:5432/{{ cookiecutter.project_name }}?sslmode=disable

Swagger Documentation

Generate Swagger docs:

make swagger-doc


View at:

http://localhost:8080/swagger/index.html

Testing

Run all tests:

make test


Run integration tests:

make itest

License

MIT License © 2026
```
