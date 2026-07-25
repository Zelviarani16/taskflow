# TaskFlow API

A REST API for managing projects and tasks, built with **Go**, **Gin**, **GORM**, **PostgreSQL**, and **JWT** authentication.

This is a personal backend portfolio project. It follows a layered
architecture (`entity → dto → repository → service → handler → routes`)
so business logic, data access, and HTTP concerns stay separated and
testable — the same pattern used in production Go services.

## Features

- JWT-based authentication (register / login, bcrypt password hashing)
- Projects: create, list (paginated), get, update, delete — scoped to the
  logged-in owner
- Tasks nested under a project: create, list (with `status` / `priority`
  filters + pagination), update, delete, optional assignment to a user
- Soft deletes on projects and tasks (GORM `deleted_at`)
- Consistent JSON response envelope (`success`, `message`, `data`)
- CLI-style commands for `--migrate`, `--seed`, `--rollback`, mirroring
  `php artisan migrate` / `db:seed` for anyone coming from Laravel

## Tech stack

| Layer      | Choice                          |
|------------|----------------------------------|
| Language   | Go 1.22                          |
| Router     | [Gin](https://github.com/gin-gonic/gin) |
| ORM        | [GORM](https://gorm.io) + PostgreSQL driver |
| Auth       | JWT ([golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt)) + bcrypt |
| Config     | [godotenv](https://github.com/joho/godotenv) |

## Project structure

```
taskflow-api/
├── cmd/               # CLI command dispatcher (--migrate/--seed/--rollback)
├── config/database/   # Postgres connection setup
├── constants/         # Shared string constants (messages, enums)
├── dto/                # Request/response shapes + sentinel errors
├── entity/            # GORM models (User, Project, Task)
├── handler/           # HTTP layer: bind request -> call service -> respond
├── helpers/           # Small stateless utilities (password hashing)
├── middleware/        # JWT auth middleware
├── migrations/        # AutoMigrate / seed / rollback
├── repository/        # Data access layer (interfaces + GORM implementation)
├── routes/            # Route registration, grouped by resource
├── service/           # Business logic (ownership checks, JWT, orchestration)
├── utils/             # Response envelope, error->HTTP status mapping
└── main.go            # Dependency wiring + server bootstrap
```

Request flow: **routes** register a path → **middleware** validates the
JWT and puts `user_id` on the context → **handler** binds/validates the
request body and calls a **service** → the service enforces business
rules (e.g. "only the project owner can edit it") and talks to a
**repository** → the repository runs the actual GORM query.

## Getting started

### 1. Start Postgres

Either run it locally, or with Docker:

```bash
docker compose up -d
```

### 2. Configure environment

```bash
cp .env.example .env
# edit .env if your DB credentials differ
```

### 3. Install dependencies, migrate, seed

```bash
go mod tidy
go run . --migrate
go run . --seed   # creates admin@taskflow.local / admin12345
```

### 4. Run the server

```bash
go run .
# -> server jalan di port 8000
```

> **Note on `go.mod`:** the `replace` block points a few `golang.org/x/*`,
> `gorm.io/*`, and `gopkg.in/*` modules at their GitHub mirrors. That was
> only needed because this project was built in a network-restricted
> sandbox that couldn't reach `proxy.golang.org` directly. It's harmless
> to keep, or you can delete the `replace (...)` block and run
> `go mod tidy` on a machine with normal internet access to resolve the
> canonical modules instead.

## API reference

All responses share this envelope:

```json
{ "success": true, "message": "...", "data": { ... } }
```

### Auth (public)

**POST `/api/v1/auth/register`**
```bash
curl -X POST localhost:8000/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"Zelvia","email":"zelvia@test.com","password":"password123"}'
```

**POST `/api/v1/auth/login`**
```bash
curl -X POST localhost:8000/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"zelvia@test.com","password":"password123"}'
```
Both return `{ token, user }`. Send the token as
`Authorization: Bearer <token>` on every route below.

### Projects (requires auth)

| Method | Path                  | Description                     |
|--------|-----------------------|----------------------------------|
| POST   | `/api/v1/projects`     | Create a project                |
| GET    | `/api/v1/projects`     | List your projects (`?page=&limit=`) |
| GET    | `/api/v1/projects/:id` | Get one project (must be owner) |
| PUT    | `/api/v1/projects/:id` | Update a project                |
| DELETE | `/api/v1/projects/:id` | Soft-delete a project           |

```bash
curl -X POST localhost:8000/api/v1/projects \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"name":"Portfolio Site","description":"Rebuild personal site"}'
```

### Tasks (requires auth, nested under a project)

| Method | Path                                   | Description                              |
|--------|------------------------------------------|-------------------------------------------|
| POST   | `/api/v1/projects/:id/tasks`             | Create a task                             |
| GET    | `/api/v1/projects/:id/tasks`             | List tasks (`?status=&priority=&page=&limit=`) |
| PUT    | `/api/v1/projects/:id/tasks/:taskId`     | Update a task (status, priority, assignee, ...) |
| DELETE | `/api/v1/projects/:id/tasks/:taskId`     | Soft-delete a task                        |

```bash
curl -X POST localhost:8000/api/v1/projects/$PROJECT_ID/tasks \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"title":"Design homepage","priority":"high"}'

curl -X PUT localhost:8000/api/v1/projects/$PROJECT_ID/tasks/$TASK_ID \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"status":"in_progress"}'
```

`status` is one of `todo` / `in_progress` / `done`.
`priority` is one of `low` / `medium` / `high`.

## Design notes

- **Ownership checks live in the service layer**, not the handler — every
  project/task mutation re-verifies `project.OwnerID == currentUserID`
  before touching the database, so authorization logic isn't duplicated
  per endpoint.
- **`found bool` return values in repositories** keep "record doesn't
  exist" (normal) separate from "the DB call failed" (a real error) — the
  service layer decides what each means instead of the repository
  guessing.
- **Sentinel errors in `dto`** (e.g. `ErrProjectForbidden`,
  `ErrTaskNotFound`) let `utils.StatusFromError` map any service error to
  the right HTTP status in one place, instead of scattering `if err ==
  ...` status logic across every handler.
