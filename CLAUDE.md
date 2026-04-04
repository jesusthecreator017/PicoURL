# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

PicoURL is a URL shortener service written in Go. It uses Redis for URL/count caching and PostgreSQL for persistent storage. The project is in early development — `cmd/server/main.go` is a placeholder and the API layer (`cmd/server/api/api.go`) is empty.

## Build & Run

```bash
go build ./...              # build all packages
go run cmd/server/main.go   # run the server
go test ./...               # run all tests
go test ./internal/store/   # run tests for a single package
```

## Architecture

- **`cmd/server/`** — Application entrypoint and HTTP API (`api/` subpackage)
- **`internal/store/`** — Data access layer behind a `Store` interface with `SaveURL`, `GetOriginalURL`, `IncrementCount`, `GetCount`. Currently only `RedisStore` implements it.
- **`internal/config/`** — Loads app config from environment variables (port, CORS, Postgres, Redis)
- **`internal/env/`** — Helpers for reading env vars with defaults (`GetString`, `GetInt`)
- **`frontend/`** — Frontend directory (empty, placeholder)

## Key Design Details

- **Module path**: `github.com/jesusthecreator017/PicoURL` (note: `redis.go` uses a different import path `github.com/jesusgonzalez07/URLShort` — this mismatch needs fixing)
- **Store interface pattern**: All data access goes through `store.Store` interface so implementations (Redis, Postgres) can be swapped
- **Redis key scheme**: URLs stored as `shortURL → originalURL`, click counts as `shortURL:count`
- **Config via env vars**: `PORT`, `CORS_ORIGIN`, `REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB`, `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`
