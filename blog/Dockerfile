# Build stage
FROM golang:1.25.4-alpine3.22 AS builder

# Install certificates (not strictly needed in builder, but okay)
RUN apk update && apk add --no-cache ca-certificates curl

WORKDIR /app

# Copy only go.mod and go.sum first (for caching dependencies)
COPY go.mod go.sum ./
RUN go mod download

# Install additional dependencies (add this)
RUN go get github.com/markbates/goth && \
    go get github.com/markbates/goth/providers/google && \
    go get github.com/markbates/goth/providers/github && \
    go get github.com/gorilla/sessions && \
    go mod tidy

# Install Air for live reloading
RUN if ! command -v air >/dev/null 2>&1; then \
    go install github.com/air-verse/air@latest; \
    fi

# install swag for swagger docs generation
RUN if ! command -v swag >/dev/null 2>&1; then \
    go install github.com/swaggo/swag/cmd/swag@latest; \
    fi

# Now copy the rest of the code
COPY . .

# generate swagger docs and fix compatibility
RUN swag init -g cmd/api/main.go --parseDependency --parseInternal --output ./docs

RUN apk add curl

# Build the app
RUN go build -o main cmd/api/main.go

# Download golang-migrate binary
RUN if [ ! -f /usr/local/bin/migrate ]; then \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz \
    | tar xvz && mv migrate /usr/local/bin/migrate; \
    fi


# Run stage
FROM debian:12-slim
WORKDIR /app


# Install CA certificates here
RUN apt-get update && apt-get install -y \
    ca-certificates \
    netcat-openbsd \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*


RUN apt-get update && apt-get install -y netcat-openbsd \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/main .
COPY --from=builder /usr/local/bin/migrate ./migrate

# Install Air binary
COPY --from=builder /go/bin/air /usr/local/bin/air

# Copy the generated swagger docs
COPY --from=builder /app/docs ./docs

# COPY .env .
COPY pkg/emailer/templates pkg/emailer/templates
COPY start.sh .
COPY wait-for.sh .
COPY .air.toml .
COPY internal/*/migration ./migration
COPY pkg/firebase/serviceAccountKey.json /app/pkg/firebase/serviceAccountKey.json

RUN chmod +x /app/start.sh /app/wait-for.sh /app/main

EXPOSE 8080
ENTRYPOINT ["/app/start.sh"]
CMD ["/app/main"]
