# Week 1 – App skeleton (api, config, domain, base layout)

## Scope

Initial application skeleton for secure-todo, including:

* HTTP API entrypoint
* Basic config layer
* Todo domain model
* Todo service layer
* HTTP router layer (Chi-based)
* HTTP middlewares with security focus (headers, cache-control, CORS, logging, panic recovery)
* Deploy directory placeholders
* Base directory layout for future components
* Typed config module with env loading and validation
* Logging module using zap with security logging policy
* Placeholder files (no logic yet) for future layers
* Verified successful build inside WSL
* Draft PR opened and continuously updated with self-review

---

## Git evidence

* Base branch: `develop`
* Feature branch: `feature/week1-app-skeleton`
* Pull Request: `#2 – Week 1 – App skeleton (api, config, domain)` (Draft)

### Key commits (milestones)

* `59b3068` – Init Week 1 app skeleton (api, config, domain)
* `d08c4ca` – Add base directory layout and placeholders
* `7a7f1e4` – Add typed config module with env loading and validation
* `43f3bbd` – Keep strict config validation and add env fatal helper
* `e510e04` – Add zap caller info and Security Logging Policy
* `051adc3` – Wire Chi router into `main` and fix logger initialization
* `2f160c0` – Add secure Chi-based HTTP router with middlewares

---

## Tasks done

### Git & branching

* [x] Updated `develop` from origin inside WSL
* [x] Created feature branch `feature/week1-app-skeleton`
* [x] Ensured all development happens only inside WSL (no builds/runs on Windows host)
* [x] Opened Draft PR and added self-review comments
* [x] Used commits with meaningful messages and small focused changes

---

### App skeleton (initial code)

* [x] Added Go application entrypoint in `cmd/api/main.go`
* [x] Added initial config loader in `internal/config/config.go`
* [x] Added domain entity in `internal/domain/todo.go`
* [x] Added service layer skeleton in `internal/service/todo_service.go`
* [x] Created `deploy/dev` and `deploy/prod` directories with `.gitkeep`
* [x] Verified build with `go build ./cmd/api` inside WSL

---

### Base directory layout (structure only)

Created the following directories:

```text
internal/config
internal/domain
internal/http/handler
internal/http/router
internal/service
internal/repository
internal/store
internal/logging
pkg/utils
```

Added placeholder files:

* `internal/http/handler/health.go`
* `internal/http/router/router.go`
* `internal/service/service.go`
* `internal/repository/repository.go`
* `internal/domain/domain.go`
* `internal/store/store.go`
* `internal/logging/logger.go`
* `pkg/utils/utils.go`

All placeholder files initially contained only a `package` declaration, to be filled in incrementally.

---

### Config module (secure env-based configuration)

* [x] Defined typed config structures in `internal/config/config.go`:

    * `ServerConfig` (Port, Env, ReadTimeout, WriteTimeout)
    * `DatabaseConfig` (DSN, MaxOpenConns, MaxIdleConns, ConnMaxLifetime)
    * `RedisConfig` (Addr, Password, DB)
    * `JWTConfig` (Secret, AccessTokenTTL, RefreshTokenTTL, Issuer)
    * `Config` (Server, Database, Redis, JWT)

* [x] Implemented `LoadConfig() (*Config, error)`:

    * Loads `.env.dev` automatically when `APP_ENV` is empty or `development`
    * Reads env variables:

        * `APP_ENV`, `SERVER_PORT`, `SERVER_READ_TIMEOUT`, `SERVER_WRITE_TIMEOUT`
        * `DB_DSN`, `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_CONN_MAX_LIFETIME`
        * `REDIS_ADDR`, `REDIS_PASSWORD`, `REDIS_DB`
        * `JWT_SECRET`, `JWT_ACCESS_TTL`, `JWT_REFRESH_TTL`, `JWT_ISSUER`

* [x] Added helpers for validation and parsing:

    * `getRequiredEnv(key string) (string, error)`
    * `parseIntEnv(key string) (int, error)`
    * `parseDurationEnv(key string) (time.Duration, error)`

* [x] Enforced strict validation:

    * No defaults for secrets
    * Required: `DB_DSN`, `JWT_SECRET`, `JWT_ACCESS_TTL`, `JWT_REFRESH_TTL`, `JWT_ISSUER`
    * Required: DB connection limits and server timeouts
    * Config loading fails fast with clear error messages if any required value is missing or invalid

* [x] Added an explicit “fatal env helper”:

    * Helper (e.g. `getEnvOrFatal`) to immediately `log.Fatal` for code paths that must not continue without specific
      env values
    * Keeps `LoadConfig` strict while still allowing some call sites to fail fast by design

* [x] Updated `cmd/api/main.go` to:

    * Call `config.LoadConfig()`
    * Use `cfg.Server.Port` for HTTP listen address
    * Propagate `cfg.Server.Env` into logger initialization and router (for CORS behavior)

---

### Logging module (zap-based)

* [x] Added logging module in `internal/logging/logger.go` using `go.uber.org/zap`

* [x] Implemented global logger:

    * `Init(env string) (*zap.Logger, error)`:

        * `env == "production"` → `zap.NewProduction(zap.AddCaller())`
        * `env != "production"` → `zap.NewDevelopment(zap.AddCaller())`
        * Returns `(*zap.Logger, error)` and does **not** exit the process itself
        * Caller (e.g. `main`) is responsible for handling initialization errors

    * `L()`:

        * Returns the initialized global logger
        * Panics if `Init` has not been called (fast failure on incorrect usage)

* [x] Logs contain at least: `time`, `level`, `caller`, `msg` (via zap + `AddCaller`)

* [x] Added Security Logging Policy (documented directly in `internal/logging/logger.go`):

    * No secrets, passwords, tokens, API keys, or full DSNs in logs
    * No JWT access/refresh tokens, session IDs, one-time codes, or similar sensitive values
    * Only non-sensitive metadata allowed (e.g. length, type, high-level identifiers, hashed IDs)
    * Logger usage throughout the app must respect this policy

---

### HTTP entrypoint & logger wiring (`cmd/api/main.go`)

* [x] Updated `cmd/api/main.go` to use the new logging and router setup:

    * Loads configuration via `config.LoadConfig()`

    * Initializes logger via `logging.Init(cfg.Server.Env)` and uses the returned `*zap.Logger`

    * Fails fast with `log.Fatalf` if:

        * Config loading fails
        * Logger initialization fails

    * Creates a temporary `/healthz` handler inline using `http.HandlerFunc`:

        * Returns HTTP 200 with body `ok`
        * Used as a simple health endpoint until a dedicated handler is implemented in `internal/http/handler`

    * Calls `router.NewRouter(logger, cfg, healthHandler)` to build the HTTP handler stack

    * Computes `addr := ":" + cfg.Server.Port` and logs a startup line:

        * Message: `"starting secure-todo api"`
        * Fields: `addr`, and environment via `cfg.Server.Env` if desired

        * Starts the HTTP server with `http.ListenAndServe(addr, r)` and logs a fatal error using the logger if
          `ListenAndServe` returns an error.

---

### HTTP router & middleware (Chi-based)

* [x] Implemented Chi-based router in `internal/http/router/router.go`:

    * `func NewRouter(logger *zap.Logger, cfg *config.Config, healthHandler http.Handler) http.Handler`

* [x] Base Chi middlewares:

    * `middleware.RequestID`:

        * Attaches a unique request ID to each incoming request for tracing

    * `middleware.RealIP`:

        * Extracts the real client IP from standard headers when behind reverse proxies

* [x] Middleware for secure HTTP headers:

    * Adds the following headers to all responses:

        * `X-Content-Type-Options: nosniff`
        * `X-Frame-Options: DENY`

    * `X-XSS-Protection` is intentionally **not** set (deprecated/ignored by modern browsers)

* [x] No-cache middleware for dynamic responses:

    * Sets:

        * `Cache-Control: no-store`

    * Ensures responses are not cached on disk by browsers or proxies, reducing risk of sensitive data being stored
      locally.

* [x] Minimal CORS for development:

    * For non-production environments (`cfg.Server.Env != "production"`):

        * `Access-Control-Allow-Origin: *`
        * `Access-Control-Allow-Methods: GET,POST,PUT,DELETE,OPTIONS`
        * `Access-Control-Allow-Headers: Content-Type, Authorization`

    * Handles `OPTIONS` preflight requests by responding with `204 No Content` and returning early

    * In production, CORS is not automatically opened and must be explicitly configured in a future step

* [x] Request logging middleware (security-aware):

    * Wraps `http.ResponseWriter` with a `responseRecorder` that tracks HTTP status codes

    * Logs **only** the following fields:

        * HTTP method
        * Request path
        * Status code
        * Duration (time taken to handle the request)

    * Does **not** log:

        * Request headers
        * Request body
        * Query parameters or cookies

    * Fully compliant with the Security Logging Policy (no sensitive content in logs)

* [x] Panic recovery:

    * Uses `middleware.Recoverer` to:

        * Catch panics in handlers
        * Prevent the entire process from crashing
        * Return an HTTP 500 instead of killing the server

* [x] Routing:

    * Registers the health endpoint:

        * `GET /healthz` → `healthHandler.ServeHTTP`
        * `healthHandler` is injected from `main`, keeping router purely focused on wiring

---

### Sanity checks

* [x] Ran `go mod tidy`

* [x] Ran `go list ./...`

* [x] Verified build with `go build ./cmd/api`

* [x] Verified runtime config with valid `.env.dev`

* [x] Verified failure modes:

    * Missing `DB_DSN` → `missing required environment variable: DB_DSN`
    * Invalid duration → `invalid duration value for SERVER_READ_TIMEOUT: "xxx"`

* [x] Verified that:

    * HTTP server starts successfully on `:<port>` from config
    * Logger is initialized before any structured logging
    * `/healthz` returns 200 with a simple `"ok"` body
    * HTTP logs include only non-sensitive metadata

---

## Current directory overview

```text
cmd/api/main.go

internal/
  config/config.go
  domain/{domain.go,todo.go}
  http/
    handler/health.go         # placeholder for future health handlers
    router/router.go          # Chi-based router with security-focused middlewares
  service/{service.go,todo_service.go}
  repository/repository.go
  store/store.go
  logging/logger.go
  storage/.gitkeep

deploy/
  dev/.gitkeep
  prod/.gitkeep

pkg/utils/utils.go

docs/progress/week1-app-skeleton.md
```

---

## How to verify

```bash
cd ~/projects/secure-todo
git checkout feature/week1-app-skeleton

go mod tidy
go list ./...
go build ./cmd/api
```

Run with valid `.env.dev`:

```bash
go run ./cmd/api
```

Expected:

* Correct startup log line (including address and environment)
* No config validation errors with a valid `.env.dev`
* Successful run
* `GET /healthz` returns HTTP 200 with body `ok`
* HTTP logs include method, path, status, and duration only
* All changes are visible in PR `#2` on branch `feature/week1-app-skeleton`
