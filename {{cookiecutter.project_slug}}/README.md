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
- [Email](#email)
- [Deployment](#deployment)
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

## Email

{% if cookiecutter.email_service == "None" -%}
No production email provider was selected — email sending falls back to **Mailpit**, a local SMTP catcher for
development only (view caught mail at `http://localhost:8025`). To send real email, set `EMAIL_PROVIDER=smtp` in
`.env` and fill in `SMTP_HOST`/`SMTP_PORT`/`SMTP_USERNAME`/`SMTP_PASSWORD` for your provider —
`pkg/emailer/smtp` works with Google, Mailgun, Zoho, SendGrid, Amazon SES and Postmark unchanged.
{%- else -%}
This project sends email through **{{ cookiecutter.email_service }}** over SMTP (`pkg/emailer/smtp`).
`generate_env_file` in the post-gen hook already pre-filled `SMTP_HOST`/`SMTP_PORT` for
{{ cookiecutter.email_service }} in `.env` — open it and fill in `SMTP_USERNAME`/`SMTP_PASSWORD` (see the comments
above each field for what {{ cookiecutter.email_service }} expects there).

Locally, `USE_MAILPIT=true` and an empty `EMAIL_PROVIDER` still route through Mailpit so you don't burn real sends
while developing. Set `EMAIL_PROVIDER=smtp` when you're ready to send through {{ cookiecutter.email_service }}.
{%- endif %}

## Deployment

This project is set up to deploy to **{{ cookiecutter.hosting_provider }}**.

{% if cookiecutter.hosting_provider == "Railway" -%}
Railway builds straight from the `Dockerfile` using `railway.toml`. Provision a PostgreSQL plugin in your Railway
project, then set `DB_SOURCE` and the other variables from `env.example` in the Railway dashboard (or
`railway variables set`), and deploy:

```bash
railway up
```
{%- elif cookiecutter.hosting_provider == "DigitalOcean" -%}
`.do/app.yaml` defines a DigitalOcean App Platform spec (API service + managed Postgres). Create the app with:

```bash
doctl apps create --spec .do/app.yaml
```

Set the `TOKEN_SYMMETRIC_KEY`/`SESSION_SECRET` secrets in the App Platform dashboard after creation — they're
declared as `SECRET` type and not stored in the spec file.
{%- elif cookiecutter.hosting_provider == "AWS" -%}
`deploy/aws/apprunner-service.json` deploys the Dockerfile to AWS App Runner via ECR. See
[`deploy/aws/README.md`](deploy/aws/README.md) for the full push-image-then-create-service walkthrough, including
where to provision Postgres (RDS) and store secrets (Secrets Manager).
{%- endif %}

License

MIT License © 2026
```
