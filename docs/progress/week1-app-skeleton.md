# Week 1 – App skeleton (api, config, domain, base layout)

## Scope

Initial application skeleton for secure-todo:

* HTTP API entrypoint
* Basic config layer
* Todo domain model
* Todo service layer
* HTTP router layer
* Deploy directory placeholders
* Base directory layout for future components
* Typed config module with env loading and validation
* Placeholder files (no logic yet)
* Verified successful build inside WSL
* Draft PR opened

## Git evidence

* Base branch: `develop`
* Feature branch: `feature/week1-app-skeleton`
* Pull Request: `#2 – Week 1 – App skeleton (api, config, domain)` (Draft)

### Key commits (milestones)

* `59b306` – Initial app skeleton (api, config, domain, service, router)
* `d08c4c` – Added base directory layout and placeholders
* `7a7f1e` – Added typed config module with env loading and validation

## Tasks done

### Git & branching

* [x] Updated `develop` from origin inside WSL
* [x] Created feature branch `feature/week1-app-skeleton`
* [x] Ensured all development happens only inside WSL
* [x] Opened Draft PR and added self-review

### App skeleton (initial code)

* [x] Added Go application entrypoint in `cmd/api/main.go`
* [x] Added initial config loader in `internal/config/config.go`
* [x] Added domain entity in `internal/domain/todo.go`
* [x] Added service layer skeleton in `internal/service/todo_service.go`
* [x] Added initial HTTP router in `internal/http/router.go`
* [x] Created `deploy/dev` and `deploy/prod` directories with `.gitkeep`
* [x] Verified build with `go build ./cmd/api`

### Base directory layout (structure only)

Created the following directories:

```
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

All placeholder files contain only a `package` declaration.

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

    * `getRequiredEnv(key string) (string, error)` – fails if env is missing/empty
    * `parseIntEnv(key string) (int, error)` – uses `strconv.Atoi`
    * `parseDurationEnv(key string) (time.Duration, error)` – uses `time.ParseDuration`
* [x] Enforced strict validation:

    * `DB_DSN`, `JWT_SECRET`, `JWT_ACCESS_TTL`, `JWT_REFRESH_TTL`, `JWT_ISSUER`, server timeouts and DB limits must be present and valid
    * Fails fast with clear error messages if env is missing or invalid (no default secrets)
* [x] Updated `cmd/api/main.go` to:

    * Call `config.LoadConfig()`
    * Log port and environment: `starting secure-todo api on :<port> (env=<env>)`
    * Use `cfg.Server.Port` for HTTP listen address

### Sanity checks

* [x] Ran `go mod tidy`
* [x] Ran `go list ./...`
* [x] Verified successful build with `go build ./cmd/api`
* [x] Verified runtime config behavior with `go run ./cmd/api` and a valid `.env.dev`:

    * `APP_ENV=development`
    * `SERVER_PORT=8082`
    * Valid duration values for timeouts and TTLs
    * Non-empty `DB_DSN` and `JWT_SECRET`

## Current directory overview

```
cmd/api/main.go

internal/
  config/config.go
  domain/{domain.go,todo.go}
  http/
    router.go
    handler/health.go
    router/router.go
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

## How to verify

Inside WSL, from repo root:

```bash
cd ~/projects/secure-todo
git checkout feature/week1-app-skeleton
go mod tidy
go list ./...
go build ./cmd/api
```

With a valid `.env.dev`:

```env
APP_ENV=development

SERVER_PORT=8082
SERVER_READ_TIMEOUT=5s
SERVER_WRITE_TIMEOUT=10s

DB_DSN=postgres://user:pass@localhost:5432/secure_todo?sslmode=disable
DB_MAX_OPEN_CONNS=10
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=30m

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

JWT_SECRET=super-secret-key-change-this
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h
JWT_ISSUER=secure-todo-api
```

Run:

```bash
go run ./cmd/api
```

Expected:

* Log line: `starting secure-todo api on :8082 (env=development)`
* No config/validation errors
* All commands succeed with no errors
* PR #2 contains all skeleton files, layout placeholders, config module changes, and this updated document
