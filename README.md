# Cookiecutter GoCore — Golang REST API Boilerplate

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791?style=flat&logo=postgresql)
![License](https://img.shields.io/badge/License-BSD--3--Clause-blue)
![CI](https://img.shields.io/badge/CI-GitHub%20Actions-success)

**GoCore** is a clean, modular **REST API boilerplate for Golang**, built for engineers who want structure, clarity, and long-term maintainability — without unnecessary abstractions.

It focuses **only on REST APIs**.

No frontend.  
No magic.  
No over-engineering.

---

## Why GoCore?

Building REST APIs in Go is powerful — but structuring them well at scale is hard.

GoCore solves this by providing:

- A predictable REST architecture
- Clear separation of responsibilities
- Strong typing and compile-time safety
- A modular design that grows with your codebase
- Tooling that supports real production workloads

If you care about **clean code, scalability, and maintainability**, GoCore is for you.

---

## Key Features

- REST-first architecture
- Gin-powered HTTP server
- Domain-driven modular structure
- SQLC for type-safe SQL queries
- PostgreSQL integration
- Docker & Docker Compose support
- Swagger / OpenAPI documentation
- Cookiecutter template for fast project generation
- Environment-based configuration
- Middleware-ready (auth, logging, rate limiting)

---

## Design Principles

- REST only
- Explicit over implicit
- Composition over inheritance
- Modules over monoliths
- Easy to test
- Easy to reason about

GoCore is intentionally boring — because boring systems are reliable.

---

## Project Structure

```text
gocore/
├── cmd/api/
│   └── main.go              # Application entry point
├── internal/
│   ├── server/              # HTTP server & middleware
│   ├── database/            # Database connection
│   ├── users/               # Example module
│   │   ├── domain.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── handler.go
│   │   ├── router.go
│   │   ├── query/           # SQLC queries
│   │   └── migration/       # Module migrations
│   └── ...
├── migration/               # Global migrations
├── scripts/                 # Helper scripts
├── docker-compose.yml
├── sqlc.yaml
├── Makefile
├── go.mod
└── README.md
```

---
## Architecture Flow

- **Router:** Route definitions  
- **Handler:** HTTP input/output  
- **Service:** Business logic  
- **Repository:** Data access  
- **Domain:** Core entities and rules  

Each layer has **one responsibility — nothing more**, keeping your code clean and maintainable.

---

## Requirements

- Go 1.22+  
- PostgreSQL 14+  
- Docker & Docker Compose (recommended)

---
**Go Modules / Packages:**

- `firebase.google.com/go v3.13.0+incompatible`  
- `github.com/cloudinary/cloudinary-go/v2 v2.14.1`  
- `github.com/gin-gonic/gin v1.11.0`  
- `github.com/go-playground/validator/v10 v10.28.0`  
- `github.com/golang-jwt/jwt v3.2.2+incompatible`  
- `github.com/google/uuid v1.6.0`  
- `github.com/jackc/pgx/v5 v5.7.6`  
- `github.com/o1egl/paseto v1.0.0`  
- `github.com/spf13/viper v1.21.0`  
- `github.com/stretchr/testify v1.11.1`  
- `github.com/swaggo/files v1.0.1`  
- `github.com/swaggo/gin-swagger v1.6.1`  
- `golang.org/x/crypto v0.47.0`  
- `golang.org/x/time v0.14.0`  
- `google.golang.org/api v0.264.0` 

## Getting Started

### Install Cookiecutter (Recommended)
### Windows
- **Option 1 :** Using Python and pip  

```bash
python --version  # Ensure Python is installed
```
```
pip install --user cookiecutter
```
```
cookiecutter --version
```

- **Option 2 :** Using Chocolatey  

```bash
choco install cookiecutter
```
```
cookiecutter --version
```

### macOS
- **Option 1 :** Using Homebrew  

```bash
brew install cookiecutter
```
```
cookiecutter --version
```

- **Option 2 :** Using Python and pip

```bash
python3 --version
```
```
pip3 install --user cookiecutter
```
```
cookiecutter --version
```


### Ubuntu / Linux
- **Option 1 :** Using apt + pip 

```bash
sudo apt update
```
```
sudo apt install -y python3 python3-pip
```
```
pip3 install --user cookiecutter
```
```
cookiecutter --version
```

- **Option 2 :** Using pipx (isolated install)

```bash
python3 -m pip install --user pipx
```
```
python3 -m pipx ensurepath
```
```
pipx install cookiecutter
```
```
cookiecutter --version
```

## Generate a New Project
```bash
cookiecutter https://github.com/michaelassa01/cookiecutter-gocore
```
You'll be prompted for some values. Provide them, then a Golang project will be created for you.

Warning: After this point, change 'Michael Assanama', 'michaelassa01', etc to your own information.

Answer the prompts with your own desired options. For example:

```
Cloning into 'cookiecutter-gocore'...
remote: Counting objects: 550, done.
remote: Compressing objects: 100% (310/310), done.
Receiving objects: 100% (550/550), done.
Resolving deltas: 100% (283/283), done.

project_name [gomacbot]: gomacbot
project_slug [gomacbot]: gomacbot
description [Behold My Awesome Goland Project!]: Behold My Awesome Goland Project!
author_name [Michael Assanama]: Michael Assanama
github_username [michaelassa01]: michaelassa01
module_path [github.com/michaelassa01/gomacbot]: github.com/michaelassa01/gomacbot
email [michaelassanama@yourdomain.com]: michaelassanama@yourdomain.com
go_version [1.25.4]: 1.25.4
Select postgresql_version:
1 - 17
2 - 16
3 - 15
4 - 14
Choose from 1, 2, 3, 4 [1]: 2
Select cloud_provider:
1 - AWS
2 - Railway
3 - Azure
4 - Digitalocean
5 - None
Choose from 1, 2, 3, 4, 5 [1]: 3
Select mail_service:
1 - Mailgun
2 - Amazon SES
3 - Mailtrap
4 - Postmark
5 - Sendgrid
6 - Other SMTP
Choose from 1, 2, 3, 4, 5, 6 [1]: 1
Select ci_tool:
1 - None
2 - Travis
3 - Gitlab
4 - Github
5 - Drone
Choose from 1, 2, 3, 4, 5 [1]: 4
Select license:
1 - MIT
2 - BSD
3 - GPLv3
4 - Apache 2.0
5 - Not open source
Choose from 1, 2, 3, 4, 5 [1]: 1
```

Enter the project and take a look around:
```bash
cd gomacbot/
ls
```

Install Dependencies
```bash
make install-dependencies
go mod tidy
```

Run Postgres with Docker (Optional)
```bash
make docker-run
```

Run the API
```bash
make run
```
Or for live reload
```
make watch
```

Create a module:
```bash
make module name=users
```

Create a migration:
```bash
make migration module=users name=add_profile_table
```

Generate SQLC and Swagger Docs
```bash
make sqlc
```
```bash
make swagger-doc
```

View Swagger at:
```
http://localhost:8080/swagger/index.html
```

Run Tests
```bash
make test       # All tests
```
```bash
make itest      # Integration tests
```


## Makefile Commands
Run ```make help ``` to see all commands. Examples:

- `make build                  # Build Go binary`
- ` make run `                  # Run API
- ` make test  `                 # Run all tests
- ` make tidy  `                 # Go module tidy
- ` make docker-run `            # Start Docker containers
- ` make docker-down   `         # Stop Docker containers
- ` make module name=users  `    # Create a new module
- ` make migration module=users name=add_profile_table ` # Create a migration
- ` make sqlc  `                 # Generate SQLC code
- ` make swagger-doc  `          # Generate Swagger docs
- ` make watch    `              # Live reload with Air




Congratulations! Your Gocore REST API project is ready, fully modular, and production-ready.

powered by [Cookiecutter](https://awesomeopensource.com/project/elangosundar/awesome-README-templates).
