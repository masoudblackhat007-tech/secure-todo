# secure-todo

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go\&logoColor=white)](https://go.dev)
[![WSL2](https://img.shields.io/badge/WSL2-Ubuntu%2022.04-orange?logo=ubuntu\&logoColor=white)](https://learn.microsoft.com/en-us/windows/wsl/)
[![Docker](https://img.shields.io/badge/Docker-Required-2496ED?logo=docker\&logoColor=white)](https://www.docker.com/)
[![Security-Focused](https://img.shields.io/badge/Secure%20By%20Design-Yes-success)](#)

A fully security-focused TODO application built in Go with PostgreSQL, Docker, and a hardened development workflow under WSL2.

---

## Features

* 🔐 **Security-first architecture**
* 🐳 **Dockerized development** (PostgreSQL, Redis, App)
* 🏗 **Clean folder structure (internal/***)**
* 🧪 **Health-check endpoints**
* 🧵 **Strict GitFlow branching**
* 📦 **Makefile automation**
* 🔎 **WSL2-only development** (no Windows binaries)

---

## Requirements

* **WSL2 Ubuntu 22.04**
* **Docker Engine (WSL2 backend)**
* **docker-compose**
* **Go 1.22+**

---

## Project Structure

```
secure-todo/
│── cmd/
│── internal/
│   ├── api/
│   ├── config/
│   ├── core/
│   ├── database/
│   ├── handlers/
│   ├── listener/
│   ├── logger/
│   ├── models/
│   ├── repository/
│   ├── services/
│   ├── usecases/
│── deployments/
│── Dockerfile
│── docker-compose.yml
│── Makefile
│── .env.example
│── README.md
```

---

## Environment Variables

Copy:

```bash
cp .env.example .env
```

`.env` contains:

```
APP_PORT=8080
DATABASE_URL=postgres://postgres:postgres@db:5432/secure_todo?sslmode=disable
REDIS_ADDR=redis:6379
```

---

## Running in Docker

### Build + Up (app + db + redis)

```bash
docker compose up --build
```

### Detached mode

```bash
docker compose up -d --build
```

### Stop

```bash
docker compose down
```

---

## Makefile Shortcuts

```bash
make up          # docker compose up --build
make down        # docker compose down
make logs        # docker logs -f <container>
make fmt         # go fmt ...
make tidy        # go mod tidy
```

---

## Health Check

```bash
curl http://localhost:8080/health
```

Response:

```json
{ "status": "ok" }
```

---

## GitFlow Rules

* **never** commit directly to `develop` یا `main`
* هر ویژگی جدید:

```bash
git checkout develop
git pull
git checkout -b feature/<name>
```

بعد از اتمام:

```bash
git checkout develop
git merge --no-ff feature/<name>
git push
git branch -d feature/<name>
git push origin --delete feature/<name>
```

---

## Dockerfile

```dockerfile
FROM golang:1.22 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o secure-todo ./cmd/main.go

FROM alpine:3.19
WORKDIR /app

COPY --from=builder /app/secure-todo .
COPY .env .env

EXPOSE 8080
CMD ["./secure-todo"]
```

---

## docker-compose.yml

```yaml
version: '3.9'

services:
  app:
    build: .
    container_name: secure-todo-app
    depends_on:
      - db
      - redis
    ports:
      - "8080:8080"
    env_file:
      - .env
    restart: always

  db:
    image: postgres:15
    container_name: secure-todo-db
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: secure_todo
    ports:
      - "5432:5432"

  redis:
    image: redis:7
    container_name: secure-todo-redis
    ports:
      - "6379:6379"
```

---

## License

MIT License.
