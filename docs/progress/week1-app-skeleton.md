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

## Tasks done

### Git & branching

* [x] Updated `develop` from origin inside WSL
* [x] Created feature branch `feature/week1-app-skeleton`
* [x] Ensured all development happens only inside WSL
* [x] Opened Draft PR and added self-review

### App skeleton (initial code)

* [x] Added Go application entrypoint in `cmd/api/main.go`
* [x] Added config loader in `internal/config/config.go` (`AppConfig{HTTPPort: "8080"}`)
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

### Sanity checks

* [x] Ran `go mod tidy`
* [x] Ran `go list ./...`
* [x] Verified successful build with `go build ./cmd/api`

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

Expected:

* All commands succeed with no errors
* PR #2 contains all skeleton files, layout placeholders, and this document
