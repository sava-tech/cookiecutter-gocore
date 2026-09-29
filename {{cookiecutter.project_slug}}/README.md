# {{ cookiecutter.project_name }}

> {{ cookiecutter.description }}

A modular Go REST API built with **Domain-Driven Design (DDD)**, **Gin**, **SQLC** and **PostgreSQL**. Generated from [cookiecutter-gocore](https://github.com/sava-tech/cookiecutter-gocore).

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Enabled-blue?logo=docker)](https://www.docker.com/)

---

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Requirements](#requirements)
- [Getting Started](#getting-started)
  - [Option A: Run everything with Docker](#option-a-run-everything-with-docker)
  - [Option B: Run the API locally](#option-b-run-the-api-locally)
- [Configuration](#configuration)
- [Rate Limiting](#rate-limiting)
- [Creating Modules](#creating-modules)
- [Migrations](#migrations)
- [Swagger Documentation](#swagger-documentation)
- [Testing](#testing)
- [Makefile Commands](#makefile-commands)

---

## Features

- Gin HTTP server with modular routing
- DDD layering per module: **domain → repository → service (application) → handler → router**
- Users, auth (password, email OTP, email verification, password reset) and social login (Google / Apple via Goth)
- PASETO/JWT access and refresh tokens
- SQLC for type-safe SQL, golang-migrate for migrations
- Per-IP rate limiting with a stricter limit on auth endpoints
- Dockerized Postgres, Redis, Mailpit and API
- Swagger / OpenAPI docs

---

## Project Structure

```text
{{ cookiecutter.project_name }}/
├── cmd/api/main.go            # Entry point
├── internal/
│   ├── server/                # Gin server, router, middleware (auth, rate limiting)
│   ├── users/                 # Users + auth module
│   │   ├── domain/            # Entities and repository interfaces
│   │   ├── application/       # Business logic (services)
│   │   ├── repository/        # SQLC-backed persistence
│   │   ├── handlers/          # HTTP handlers
│   │   ├── dto/               # Request/response payloads
│   │   ├── query/             # SQL queries for SQLC
│   │   ├── migration/         # Module migrations
│   │   └── router.go          # Route registration
│   ├── socialauth/            # Google / Apple login module
│   └── shared/                # Cross-module helpers
├── pkg/                       # Reusable packages (token, response, ...)
├── utils/                     # Config loading and helpers
├── scripts/                   # Module scaffolding scripts
├── docs/                      # Generated Swagger docs
├── env.example                # Every supported environment variable
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── sqlc.yaml
```

---

## Requirements

- Go 1.22+
- Docker and Docker Compose (for Postgres, Redis and Mailpit)
- Make

The CLI tools (`sqlc`, `swag`, `migrate`, `air`, ...) are installed for you in step 3 below.

---

## Getting Started

### 1. Create your `.env`

Copy the example file and fill in real values:

```bash
cp env.example .env
```

At minimum, change:

- `TOKEN_SYMMETRIC_KEY`: must be **exactly 32 characters**, or the app refuses to start.
- `SESSION_SECRET`: at least 32 random characters.

See [Configuration](#configuration) for the rest.

### 2. Initialise Git (recommended)

```bash
git init
git add .
git commit -m "Initial commit"
```

### 3. Install dependencies and tools

```bash
make install-dependencies
```

This runs `go mod tidy` and installs `migrate`, `sqlc`, `swag`, `golangci-lint` and `mockgen` into your `$GOPATH/bin`. Make sure that directory is on your `PATH`.

### 4. Generate Swagger docs

`cmd/api/main.go` imports the generated `docs` package, so the project **won't compile until you run this once**:

```bash
make swagger-doc
```

Re-run it whenever you change handler annotations.

Now pick how you want to run the API.

### Option A: Run everything with Docker

```bash
make docker-run
```

This builds the API image and starts Postgres, Redis, Mailpit and the API. On start-up the container waits for its dependencies, runs the migrations and then starts the API with live reload.

| Service | URL |
| --- | --- |
| API | http://localhost:8080 |
| Health check | http://localhost:8080/public/health |
| Swagger | http://localhost:8080/swagger/index.html |
| Mailpit (catches outgoing email) | http://localhost:8025 |

Stop everything with:

```bash
make docker-down
```

### Option B: Run the API locally

Use this when you want the API on your machine and only the backing services in Docker.

1. **Start the backing services:**

   ```bash
   docker compose up -d postgres redis mailpit
   ```

2. **Point `.env` at localhost.** The example values use Docker service names. For a local run, change them to:

   ```dotenv
   DB_SOURCE=postgres://root:secret@localhost:5432/{{ cookiecutter.project_name }}_db?sslmode=disable
   REDIS_ADDR=localhost:6379
   MAILPIT_HOST=127.0.0.1
   ```

3. **Load `.env` into your shell.** Config is read **only from environment variables**; the `.env` file isn't read automatically.

   ```bash
   set -a && source .env && set +a
   ```

4. **Run the migrations** (see [Migrations](#migrations)):

   ```bash
   migrate -path internal/users/migration -database "$DB_SOURCE" up
   migrate -path internal/socialauth/migration -database "$DB_SOURCE" up
   ```

5. **Start the API:**

   ```bash
   make run      # plain run
   make watch    # live reload with air
   ```

6. **Check it's up:**

   ```bash
   curl http://localhost:8080/public/health
   # {"status":"ok"}
   ```

---

## Configuration

All settings come from environment variables; `env.example` lists every one. The most important:

| Variable | Purpose |
| --- | --- |
| `SERVER_ADDRESS` | Address the API listens on (default `0.0.0.0:8080`) |
| `DB_SOURCE` | Postgres connection string |
| `TOKEN_SYMMETRIC_KEY` | Token signing key, exactly 32 characters |
| `ACCESS_TOKEN_DURATION` / `REFRESH_TOKEN_DURATION` | Token lifetimes, e.g. `15m`, `720h` |
| `SESSION_SECRET` | Cookie session secret used by social login |
| `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET` | Google OAuth credentials |
| `SOCIAL_CALLBACK_URL` | OAuth callback URL |
| `USE_MAILPIT`, `EMAIL_PROVIDER`, ... | Email delivery |
| `TRUSTED_PROXIES` | Proxies allowed to set `X-Forwarded-For` (see below) |
| `RATE_LIMIT_*`, `AUTH_RATE_LIMIT_*` | Rate limits (see below) |
| `PRODUCTION` | `true` runs the compiled binary in Docker instead of `air` |

---

## Rate Limiting

Every request is rate-limited per client IP by `internal/server/middleware/ratelimit.go`. There are two tiers:

| Tier | Applies to | Default | Variables |
| --- | --- | --- | --- |
| Global | All routes | 5 req/s, burst 10 | `RATE_LIMIT_RPS`, `RATE_LIMIT_BURST` |
| Auth | `register`, `login`, OTP, `verify-email`, `forgot-password`, `reset-password` | 5 req/min, burst 5, **per endpoint** | `AUTH_RATE_LIMIT_PER_MINUTE`, `AUTH_RATE_LIMIT_BURST` |

When a client goes over the limit it gets `429 Too Many Requests` with a `Retry-After` header and the standard error body:

```json
{
  "success": false,
  "message": "Too many requests, slow down.",
  "error": "too_many_requests",
  "error_code": "RATE_LIMIT_EXCEEDED"
}
```

**Behind a load balancer or reverse proxy?** Set `TRUSTED_PROXIES` to its IP or CIDR (comma-separated for several), for example `TRUSTED_PROXIES=10.0.0.0/8`. The client IP is only read from `X-Forwarded-For` when the request comes from a trusted proxy:

- **Left empty** (the default): the header is ignored, so clients can't spoof their IP to dodge the limit.
- **Empty while behind a proxy:** every request looks like it comes from the proxy, and all users share one budget.

To protect a new route with the stricter tier, add the `authRateLimit` middleware when you register the route:

```go
authGroup.POST("/login", authRateLimit, auth.Login)
```

> **Scaling note:** limiter state is kept in memory, so each API instance counts separately. With N replicas the effective limit is N× higher. Move the limiter to Redis if you need a shared limit.

---

## Creating Modules

```bash
make create-module name=posts
```

This scaffolds `internal/posts/` with the DDD folders, creates its first migration and registers the module's queries in `sqlc.yaml`. Then:

1. Write the tables in `internal/posts/migration/*.up.sql` (and the reverse in `*.down.sql`).
2. Write the queries in `internal/posts/query/*.sql`.
3. Generate the Go code: `make sqlc`.
4. Implement the repository, service and handlers.
5. Register the module's routes in `internal/server/router.go`.

To remove a module, run the command below. It deletes `internal/posts/` and its `sqlc.yaml` entry, then lists any files that still import it. Add `FORCE=1` to skip the confirmation prompt.

```bash
make delete-module name=posts
```

---

## Migrations

Each module keeps its own migrations in `internal/<module>/migration/`.

**Create a migration:**

```bash
make migration module=users name=add_profile_table
```

This creates matching `*.up.sql` and `*.down.sql` files. Fill in both.

**Apply migrations** (with `.env` loaded into your shell):

```bash
migrate -path internal/users/migration -database "$DB_SOURCE" up
```

**Roll back the last migration:**

```bash
migrate -path internal/users/migration -database "$DB_SOURCE" down 1
```

With Docker (Option A), migrations run automatically every time the API container starts.

---

## Swagger Documentation

```bash
make swagger-doc
```

Then open http://localhost:8080/swagger/index.html.

---

## Testing

```bash
make test     # all tests
make itest    # integration tests
```

---

## Makefile Commands

Run `make help` for the full list.

| Command | What it does |
| --- | --- |
| `make install-dependencies` | Tidy modules and install CLI tools |
| `make run` | Run the API |
| `make watch` | Run with live reload (`air`) |
| `make build` | Build the `main` binary |
| `make test` / `make itest` | Run all / integration tests |
| `make docker-run` / `make docker-down` | Start / stop the Docker stack |
| `make create-module name=<name>` | Scaffold a new module |
| `make delete-module name=<name>` | Delete a module and its SQLC config |
| `make migration module=<m> name=<n>` | Create a migration |
| `make sqlc` | Generate SQLC code |
| `make swagger-doc` | Generate Swagger docs |
| `make tidy` | `go mod tidy` |
| `make mock` | Regenerate the users repository mock |
