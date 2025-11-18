# Week 1 – App skeleton (api, config, domain)

## Scope

Initial application skeleton for secure-todo:
- HTTP API entrypoint
- Basic config layer
- Todo domain model
- Todo service skeleton
- HTTP router skeleton
- Deploy directory placeholders

## Git evidence

- Base branch: `develop`
- Feature branch: `feature/week1-app-skeleton`
- Pull Request: `#2 – Week 1 – App skeleton (api, config, domain)` (Draft)

### Key commit

- SHA: 59b306   <!-- این رو بعداً با SHA واقعی خودت عوض کن -->

## Tasks done

- [x] Created feature branch from up-to-date `develop` in WSL.
- [x] Added Go module entrypoint in `cmd/api/main.go`.
- [x] Added config loader in `internal/config/config.go` with HTTPPort field.
- [x] Added `Todo` domain model in `internal/domain/todo.go`.
- [x] Added `TodoService` skeleton in `internal/service/todo_service.go`.
- [x] Added HTTP router skeleton in `internal/http/router.go`.
- [x] Added deploy directories `deploy/dev` and `deploy/prod` with `.gitkeep`.
- [x] Verified `go build ./cmd/api` passes inside WSL.
- [x] Created Draft PR and added self-review on GitHub.

## How to verify

Inside WSL, from repo root:

```bash
git checkout feature/week1-app-skeleton
go build ./cmd/api
