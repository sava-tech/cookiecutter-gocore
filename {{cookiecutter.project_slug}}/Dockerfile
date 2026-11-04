# Tools shared by the dev and build stages
FROM golang:1.27.3-alpine3.24 AS tools

RUN apk add --no-cache ca-certificates curl

WORKDIR /app

# Copy only go.mod and go.sum first (for caching dependencies)
COPY go.mod go.sum ./
RUN go mod download

# Install Air for live reloading
RUN go install github.com/air-verse/air@latest

# install swag for swagger docs generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Download golang-migrate binary
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.20.1/migrate.linux-amd64.tar.gz \
    | tar xvz -C /usr/local/bin migrate


# Dev stage: docker-compose mounts the source over /app and Air rebuilds on
# change. The scripts live outside /app so the mount doesn't hide them.
FROM tools AS dev
COPY start.sh wait-for.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/start.sh /usr/local/bin/wait-for.sh
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/start.sh"]


# Build stage
FROM tools AS builder

# Now copy the rest of the code
COPY . .

# generate swagger docs and fix compatibility
RUN swag init -g cmd/api/main.go --parseDependency --parseInternal --output ./docs

# Build the app
RUN go build -o main cmd/api/main.go


# Run stage (the default target, used for deployments)
FROM debian:13-slim
WORKDIR /app


# Install CA certificates here
RUN apt-get update && apt-get install -y \
    ca-certificates \
    netcat-openbsd \
    && update-ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/main .
COPY --from=builder /usr/local/bin/migrate ./migrate

# Copy the generated swagger docs
COPY --from=builder /app/docs ./docs

# COPY .env .
COPY pkg/emailer/templates pkg/emailer/templates
COPY start.sh .
COPY wait-for.sh .
COPY internal/*/migration ./migration
COPY pkg/firebase/serviceAccountKey.json /app/pkg/firebase/serviceAccountKey.json

RUN chmod +x /app/start.sh /app/wait-for.sh /app/main

EXPOSE 8080
ENTRYPOINT ["/app/start.sh"]
CMD ["/app/main"]
