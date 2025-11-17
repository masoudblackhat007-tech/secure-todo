# secure-todo

پروژه تمرینی برای ساخت یک TODO API امن با Go.

## پیش‌نیازها (فقط داخل WSL)

- WSL2 با یک توزیع لینوکس (مثلاً Ubuntu)
- Go 1.22 یا بالاتر
- Docker و Docker Compose
- `psql` (کلاینت PostgreSQL)
- `redis-cli`
- `trivy`
- `govulncheck`
- `golangci-lint`

> همه دستورات زیر باید **داخل WSL** اجرا شوند، نه روی ویندوز اصلی.

## راه‌اندازی سرویس‌های توسعه (Postgres + Redis)

```bash
docker compose -f docker-compose.dev.yml up -d
