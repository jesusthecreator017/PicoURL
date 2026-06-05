# PicoURL

A fast, minimal URL shortener with a Go backend and a React/TypeScript frontend.
It pairs a [chi](https://github.com/go-chi/chi) HTTP API with PostgreSQL for
durable storage and Redis as a write-through cache, and serves a small Vite SPA
for the web UI. Shortcodes are **deterministic** (same URL → same code), which
makes creates idempotent and collisions detectable.

> Live demo: [picourl.xyz](https://picourl.xyz)

---

## Features

- 🔗 **Shorten any URL** — `POST /api/shorten` returns a 7-character code.
- ♻️ **Deterministic codes** — SHA-256 → base62, first 7 chars. The same URL
  always maps to the same shortcode, so re-shortening is a no-op.
- ✅ **Validation** — checks URL format/scheme and performs a lightweight HTTP
  `HEAD` reachability check before storing.
- ⚡ **Redis cache** — write-through on create, read-through on redirect with a
  fallback to Postgres and cache backfill.
- 📈 **Click stats** — per-code click counts plus a global total.
- 🖥️ **Bundled SPA** — the production Go binary serves the built React app from
  `/static`.

## API

| Method   | Path                      | Body / Params         | Response                                  |
| -------- | ------------------------- | --------------------- | ----------------------------------------- |
| `POST`   | `/api/shorten`            | `{"url": "..."}`      | `201 {"short_url": "<code>"}`             |
| `GET`    | `/api/stats/{shortcode}`  | —                     | `200 {"short_url": "<code>", "click_count": N}` |
| `GET`    | `/api/total`              | —                     | `200 {"total": N}`                        |
| `DELETE` | `/api/{shortcode}`        | —                     | `200 {"message": "deleted"}`              |
| `GET`    | `/{shortcode}`            | —                     | `301` redirect to the original URL (increments clicks) |
| `GET`    | `/`                       | —                     | Serves the React SPA (when `/static` exists) |

`short_url` is the bare shortcode; build the full link as
`https://<host>/<short_url>`.

CORS is controlled by `CORS_ORIGIN` (default `*`), so the API can be embedded in
another site (e.g. a portfolio live demo). Browser preflight (`OPTIONS`) is
answered with `204`.

## Architecture

**Request flow:** `main.go` → `api.Application` (chi router) →
`service.URLService` (business logic) → `store.CachedStore` (Redis + Postgres).

- **`cmd/server/api/`** — HTTP layer: routes & CORS (`api.go`), handlers
  (`handlers.go`), request/response types (`types.go`), JSON helpers
  (`helpers/json.go`).
- **`internal/service/`** — business logic behind the `URLService` interface:
  validation, deterministic shortcode generation, collision detection. Domain
  errors live in `errors.go`.
- **`internal/store/`** — data layer behind the `Store` interface.
  `PostgresStore` (via sqlc) is the source of truth; `RedisCache` provides
  get/set/incr/del; `CachedStore` composes both.
- **`internal/shortcode/`** — SHA-256 hash → base62, 7 chars.
- **`internal/utils/`** — URL format + reachability validation.
- **`internal/config/`, `internal/env/`** — env-var configuration.
- **`frontend/`** — React 19 + Vite SPA (TanStack Router/Query). Built output is
  served by the Go server from `/static`.

## Getting started

### Prerequisites

- Go 1.25+
- Docker (for Postgres 17 + Redis 7), or your own Postgres/Redis
- Node 22+ and pnpm (only to work on the frontend)

### 1. Start infrastructure

```bash
docker compose up -d        # Postgres on :5432, Redis on :6379
```

The Postgres container auto-applies `internal/store/migrations/001_create_url_table.sql`
on first boot. (Using your own Postgres? Apply that file manually with `psql`.)

### 2. Configure the environment

The server loads `.env` via `godotenv`. A working local config:

```env
PORT=8080
CORS_ORIGIN=*

POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB=pico_url
POSTGRES_SSLMODE=disable

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_CACHE_TTL_MINUTES=1440
```

> Make sure `POSTGRES_USER`/`POSTGRES_PASSWORD` match whatever Postgres you point
> at (the bundled `docker-compose.yml` uses `postgres` / `password`).

### 3. Run the server

```bash
go run cmd/server/main.go     # serves on :8080
```

### 4. (Optional) Run the frontend in dev

```bash
cd frontend
pnpm install
pnpm dev                      # Vite dev server
```

## Try it

```bash
# Shorten
curl -s -XPOST localhost:8080/api/shorten \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com"}'
# => {"short_url":"<code>"}

# Follow the redirect
curl -i localhost:8080/<code>            # => 301 Location: https://example.com

# Stats / total
curl -s localhost:8080/api/stats/<code>  # => {"short_url":"<code>","click_count":1}
curl -s localhost:8080/api/total         # => {"total":1}
```

## Configuration

| Variable                  | Default          | Description                              |
| ------------------------- | ---------------- | ---------------------------------------- |
| `PORT`                    | `8080`           | HTTP listen port                         |
| `CORS_ORIGIN`             | `*`              | Allowed CORS origin for the API          |
| `POSTGRES_HOST`           | `localhost`      | Postgres host                            |
| `POSTGRES_PORT`           | `5432`           | Postgres port                            |
| `POSTGRES_USER`           | `postgres`       | Postgres user                            |
| `POSTGRES_PASSWORD`       | `password`       | Postgres password                        |
| `POSTGRES_DB`             | `pico_url`       | Postgres database                        |
| `POSTGRES_SSLMODE`        | `disable`        | Postgres SSL mode                        |
| `REDIS_ADDR`              | `localhost:6379` | Redis address                            |
| `REDIS_PASSWORD`          | _(empty)_        | Redis password                           |
| `REDIS_DB`                | `0`              | Redis database index                     |
| `REDIS_CACHE_TTL_MINUTES` | `1440`           | Cache TTL in minutes (default 24h)       |

## Development

```bash
go build ./...        # build
go test ./...         # test
go vet ./...          # vet
gofmt -l .            # formatting check

# Regenerate sqlc code after editing internal/store/sqlc/queries.sql
sqlc generate
```

CI runs `go build`, `go vet`, `go test`, a `gofmt` check, and `golangci-lint run`.

## Deployment

A multi-stage `Dockerfile` builds the frontend, compiles a static Go binary, and
produces a small Alpine runtime image that serves both the API and the SPA.

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

The production compose file expects a `.env.prod` (used by the app, Postgres, and
Redis) and runs all three services on an internal network with health checks.

## Tech stack

Go · chi · PostgreSQL (pgx + sqlc) · Redis · React 19 · TypeScript · Vite ·
TanStack Router/Query · Docker
