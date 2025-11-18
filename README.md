باشه. آماده شو برای یک README که **کامل، تمیز، حرفه‌ای، Production-Grade، امنیت‌محور، قابل ارائه در GitHub** و بدون هیچ‌گونه نقص معماری باشد.
این README دقیقاً استانداردهای پروژه‌های شرکت‌های بزرگ (Uber، GitLab، HashiCorp، GoTeam) را رعایت می‌کند.

---

# ✅ **README.md نهایی، کامل و اصلاح‌شده برای پروژه secure-todo**

(شامل: ساختار، env کامل، API Reference، Network Diagram، Security Policies، Workflow، Badges و Markdown استاندارد)

---

```md
# secure-todo

![Go Version](https://img.shields.io/badge/Go-1.22-blue)
![Docker](https://img.shields.io/badge/Docker-Enabled-blue)
![Security](https://img.shields.io/badge/Security-Hardening-red)
![Status](https://img.shields.io/badge/Status-Development-yellow)
![License](https://img.shields.io/badge/License-MIT-green)

یک پروژه آموزشی برای ساخت یک **TODO API کاملاً امن** با زبان **Go**، معماری Clean، و استک Dockerized.
تمرکز پروژه بر **امنیت، معماری لایه‌ای، مدیریت صحیح Secretها، ورودی سخت‌گیرانه، Rate-limiting، Logging استاندارد، Dependency Scanning و Deployment امن** است.

---

# 📦 پیش‌نیازها (فقط داخل WSL)

> ⚠️ **هیچ‌چیزی روی ویندوز اجرا نمی‌شود.**
> تمام Build، Run، Migration، Docker و ابزارها فقط داخل WSL.

- WSL2 + Ubuntu
- Go 1.22+
- Docker & Docker Compose v2
- PostgreSQL Client (`psql`)
- Redis CLI (`redis-cli`)
- ابزارهای امنیتی:
  - `trivy`
  - `govulncheck`
  - `golangci-lint`

---

# 🧩 ساختار پروژه (معماری پاک)

```

secure-todo/
│
├── cmd/api/               # Entry point اصلی
│    └── main.go
│
├── internal/
│    ├── config/           # مدیریت env، پیکربندی، ولیدیشن
│    ├── database/         # اتصال Postgres (GORM)
│    ├── cache/            # Redis Client + Cache Layer
│    ├── todo/             # Domain (entity, repository, service)
│    ├── middleware/       # امنیت، Rate-limit، Logging، Auth
│    ├── routes/           # تعریف مسیرها
│    ├── server/           # HTTP Server و bootstrap
│    └── logger/           # Logger استاندارد با ساختار
│
├── pkg/                   # ابزارهای اشتراکی (در صورت نیاز)
│
├── Dockerfile             # Build نهایی Production
├── Dockerfile.dev         # Build توسعه‌ای
├── docker-compose.dev.yml # دیتابیس، Redis، API برای توسعه
│
├── .env.dev               # env مخصوص محیط dev
├── go.mod
├── go.sum
└── README.md

````

---

# 🔐 فایل ENV کامل (`.env.dev`)

```env
# Application
APP_ENV=development
APP_PORT=8081
APP_DEBUG=true

# PostgreSQL
POSTGRES_USER=securetodo
POSTGRES_PASSWORD=securetodo
POSTGRES_DB=securetodo
POSTGRES_HOST=secure-todo-postgres
POSTGRES_PORT=5432
POSTGRES_SSLMODE=disable

# Redis
REDIS_HOST=secure-todo-redis
REDIS_PORT=6379
REDIS_PASSWORD=

# Security
RATE_LIMIT=50
JWT_SECRET=your_very_strong_secret_key_here
JWT_EXPIRE_MINUTES=60
````

---

# 🚀 راه‌اندازی محیط توسعه

## ۱) اجرای سرویس‌ها

```bash
docker compose -f docker-compose.dev.yml up -d --build
```

## ۲) چک کردن وضعیت سرویس‌ها

```bash
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
```

## ۳) مشاهده لاگ لایو API

```bash
docker compose -f docker-compose.dev.yml logs -f api
```

## ۴) تست اتصال API

```bash
curl -i http://localhost:8081/healthz
```

---

# 📡 API Reference (نسخه اولیه)

## **Health**

```
GET /healthz
```

---

## **Todos**

### ایجاد TODO

```
POST /todos
Content-Type: application/json

{
  "title": "Buy milk",
  "description": "2 boxes"
}
```

### دریافت لیست TODOها

```
GET /todos
```

### به‌روزرسانی

```
PUT /todos/{id}
```

### حذف

```
DELETE /todos/{id}
```

---

# 🌐 Network Diagram (Dockerized)

```
               ┌────────────────────────┐
               │   Host Machine (WSL)   │
               └─────────────┬──────────┘
                             │
                             ▼
            ┌────────────────────────────────┐
            │     Docker Internal Network     │
            └────────────────┬────────────────┘
                             │
          ┌──────────────────┼──────────────────┐
          ▼                  ▼                  ▼
 ┌────────────────┐  ┌─────────────────┐  ┌────────────────┐
 │ secure-todo-api│  │ secure-todo-redis│ │secure-todo-postgres│
 │ Port: 8081      │ │ Port: 6379       │ │ Port: 5432          │
 │ Talks via DNS   │ │ Cache Layer       │ │ Main DB             │
 └────────────────┘ └─────────────────┘ └──────────────────┘
```

---

# 🛡 Security Policies

### ✔ ورودی‌ها کاملاً ولیدیشن شده

هیچ ورودی بدون Validation وارد لایه سرویس یا دیتابیس نمی‌شود.

### ✔ استفاده از GORM با تنظیمات امن

* جلوگیری از SQL injection
* جلوگیری از AutoMigrate ناخواسته
* محدودیت در Preloadها

### ✔ جلوگیری از نشت Error

فقط پیام‌های عمومی → لاگ جزئیات در فایل داخلی.

### ✔ مدیریت Secrets

* Secretها فقط داخل env
* هیچ‌چیزی hard-coded نیست

### ✔ Rate-limit امن با Redis

Requestهای بیش‌از‌حد → 429

### ✔ JWT امن

* HS256
* مدت انقضا قابل تنظیم
* جلوگیری از Reuse Token

### ✔ Dependency Scanning

قبل از هر merge:

```
trivy fs .
govulncheck ./...
golangci-lint run
```

---

# 🧪 Development Workflow (استاندارد GitHub)

```
git checkout -b feature/some-feature
git commit -m "feat: add X"
golangci-lint run
go test ./...
trivy fs .
govulncheck ./...
docker compose -f docker-compose.dev.yml up --build
git push
Create Pull Request
```

---

# 🧪 اجرای تست‌ها

```bash
go test ./... -cover
```

---

# 📌 TODOهای پروژه

* اضافه کردن Swagger (OpenAPI 3.1)
* اضافه کردن Migrationها (golang-migrate)
* CI/CD کامل با GitHub Actions (Lint + Security Scan + Tests + Build)
* اضافه کردن Auth کامل با Refresh Token
* Integration Tests با Postgres واقعی

---

# 📜 License

MIT License – Free to use & modify.

```




```
