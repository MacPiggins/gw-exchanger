# gw-exchanger

A lightweight Go gRPC service for retrieving exchange rates from a PostgreSQL-backed data store.

> Part of the **GW stack** — see the [general deployment repository](https://github.com/MacPiggins/gw-deploy) for deployment configuration and infrastructure.

## Overview

`gw-exchanger` exposes exchange-rate queries over gRPC and loads rates from a `Rates` table in PostgreSQL. On startup it creates the database connection, runs automatic migrations, and serves the gRPC API on the configured port.

## Features

- gRPC API for single and bulk exchange-rate lookups
- PostgreSQL persistence with automatic schema migration
- Health-aware database connection maintenance
- Docker and Docker Compose support for local development and deployment
- Structured JSON logging

## Project structure

- `cmd/main.go` — application entry point
- `internal/app/app.go` — bootstraps storage, migrations, and gRPC server
- `internal/server/server.go` — gRPC handlers
- `internal/service/exchange.go` — exchange business logic
- `internal/storage/postgres/postgres.go` — PostgreSQL access layer
- `internal/storage/postgres/migrations/` — migration files
- `docker-compose.yaml` — local runtime stack with Postgres and the app
- `config.env.example` — sample environment configuration

## Requirements

- Go 1.26+
- PostgreSQL 16+
- Docker and Docker Compose (optional, for local containerized setup)

## Environment configuration

Copy the example environment file:

```bash
cp config.env.example config.env
```

Example `config.env`:

```env
CONNSTR=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
ADDRESS=:8080
```

### Variables

- `CONNSTR` — PostgreSQL connection string used by the app
- `ADDRESS` — address the gRPC server listens on (default `:8080`; use `:8080` in the app container)
- `APP_PORT` — host port published by Docker Compose (default `8181`; override it in the shell or Compose `.env` file)
- `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` — Compose database settings (override them in the shell or Compose `.env` file)

## Run locally

### With Go

```bash
go run ./cmd/main.go
```

Make sure PostgreSQL is running and that `CONNSTR` points to a reachable database.

### With Docker Compose

Copy `config.env.example` to `config.env`, then start the stack:

```bash
cp config.env.example config.env
docker compose up --build
```

This starts:

- PostgreSQL on port `5432`
- the app on `localhost:${APP_PORT}` (with the example config, `localhost:8181`)

In Compose, the app connects to the database using the service hostname `postgres`; `localhost` inside the app container would refer to the app container itself. The app listens on container port `8080`, while `APP_PORT` selects the host-side port. Compose reads these overrides from the shell or project `.env`, not from `config.env`.

## API

The service implements the `ExchangeService` gRPC interface with the following methods:

- `GetExchangeRate`
- `GetExchangeRates`

The app reads exchange rates from Postgres and returns them as typed gRPC responses. Ensure the process is running and that the client targets the published host port (for example, `192.168.16.6:8181` when `APP_PORT=8181`).

## Database migrations

The app automatically runs migrations at startup via Goose. Migration SQL files live under:

- `internal/storage/postgres/migrations/`

## Docker image

The service also includes a `Dockerfile` for building a minimal container image:

```bash
docker build -t gw-exchanger .
```

## Testing

Run the project test suite with:

```bash
go test ./...
```

## Notes

- The app loads environment variables using `godotenv`, so a local `config.env` file is supported.
- The application binds the gRPC server with `net.Listen("tcp4", conf.Address)` and serves requests using gRPC.
- The Postgres client includes reconnection logic and auto-migration during startup.
