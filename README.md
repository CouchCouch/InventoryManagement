<div align="center">

# Inventory Management

</div>

## Description

A simple inventory tracking system with basic CRUD operations using PostgreSQL, Go, React with Vite and Tailwind CSS.

---

## Prerequisites

- **Go 1.25+** — [download](https://go.dev/dl/)
- **Bun** — package manager & JS runtime ([Bun](https://bun.sh/))
- **PostgreSQL**

## Setup

### 1. Database

Create a database:

```sh
createdb inventory_management
```

### 2. Configuration

Copy the example config and fill in your values:

```sh
cp example-config.yaml config.yaml
```

At minimum, set your database password, admin credentials, and generate JWT secrets:

```sh
# Generate two random secrets
openssl rand -base64 64   # use this for jwt_secret
openssl rand -base64 64   # use this for jwt_refresh_secret
```

> `config.yaml` is gitignored — your secrets stay local.

### 3. Install JS dependencies

```sh
cd web && bun install && cd ..
```

### 4. Run (development)

Start the API and Vite dev server concurrently (two terminals):

```sh
# Terminal 1 — Go API (auto-migrates DB schema on startup)
go run cmd/api/main.go

# Terminal 2 — Vite frontend
cd web && bunx --bun vite
```

Or use the Makefile (builds the frontend first, then starts the API):

```sh
make run
```

### 5. Open

The API serves at **http://localhost:3000** and the Vite dev server at **http://localhost:5173** (proxied by the API in production, but during dev the frontend talks directly to the API via CORS).

## Useful Commands

```sh
# Run tests
go test ./...

# Verbose tests
make testv

# Format code
make fmt

# Lint
make lint

# Production build
make build
```

## Project Structure

```
├── cmd/api/main.go        # Entry point
├── internal/
│   ├── auth/              # JWT token generation
│   ├── config/            # YAML config loading
│   ├── db/                # PostgreSQL layer, migrations, query builder
│   ├── domain/            # Data models & error types
│   ├── handlers/          # HTTP handlers & middleware
│   └── logger/            # Structured logging (zerolog via slog)
├── web/                   # React + Vite + Tailwind frontend
├── config.yaml            # Local config (gitignored)
├── example-config.yaml    # Template for config.yaml
├── Makefile               # Build, test, lint, run shortcuts
└── code-review-todo.md    # Outstanding issues tracked from code review
```

## To-do

### Items CRUD operations
- [x] Create Items
- [x] Read Items
- [x] Update Items
- [x] Delete Items

### Checkout Functionality
- [x] Checkout
- [x] Return
- [x] View all checked out items

### Other
- [ ] Search inventory
- [ ] Export Items
- [ ] Import Items
