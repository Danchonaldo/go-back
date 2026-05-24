# 📋 TaskBoard — Trello-like Task Management System

Full-stack project built with **Go (Gin)** backend and **React** frontend.

## Architecture

```
taskboard/
├── main_service/          # Main REST API (Go + Gin + GORM)
│   ├── handlers/          # HTTP handlers (auth, boards, tasks, users)
│   ├── middleware/        # JWT auth middleware, logger
│   ├── models/            # GORM models
│   ├── db/                # DB connection + migrations
│   │   └── migrations/    # SQL migration files (golang-migrate)
│   ├── tests/             # 10 unit tests
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── notification_service/  # Microservice (Go + Gin)
│   ├── main.go
│   ├── go.mod
│   └── Dockerfile
├── frontend/              # React SPA
│   ├── src/
│   │   ├── pages/         # Login, Register, Dashboard, BoardDetail
│   │   ├── services/      # Axios API service
│   │   └── context/       # Auth context
│   ├── Dockerfile
│   └── nginx.conf
├── docker-compose.yml
└── Makefile
```

## ✅ Requirements Coverage

| Requirement | Implementation |
|---|---|
| REST API (Gin) | `main_service` with 20 endpoints |
| Postman | 20 endpoints documented below |
| PostgreSQL + GORM | `db/database.go`, all models use GORM |
| JWT Authentication | `handlers/auth_handler.go` + `middleware/auth_middleware.go` |
| Migrations (golang-migrate) | `db/migrations/*.sql` |
| 10 Unit Tests | `tests/handlers_test.go` |
| Middleware | JWT Auth + CORS + Logger |
| Microservices (Resty v2) | `handlers/notification_client.go` → calls notification_service |
| Docker + Docker Compose | `Dockerfile` × 3, `docker-compose.yml` |
| Frontend | React + React Router, Kanban UI |

## 🚀 Quick Start

```bash
# Start everything
make up

# Or manually:
docker-compose up --build
```

Services:
- **Frontend**: http://localhost:3000
- **Main API**: http://localhost:8082
- **Notification Service**: http://localhost:8083

## 📡 API Endpoints (20 total)

### Auth (public)
| # | Method | URL | Description |
|---|--------|-----|-------------|
| 1 | POST | `/auth/register` | Register new user |
| 2 | POST | `/auth/login` | Login, returns JWT |
| 3 | GET | `/auth/me` | Get current user profile |

### Users (protected)
| # | Method | URL | Description |
|---|--------|-----|-------------|
| 4 | GET | `/api/users` | Get all users |
| 5 | GET | `/api/users/:id` | Get user by ID |
| 6 | PUT | `/api/users/:id` | Update user |
| 7 | DELETE | `/api/users/:id` | Delete user |

### Boards (protected)
| # | Method | URL | Description |
|---|--------|-----|-------------|
| 8 | POST | `/api/boards` | Create board |
| 9 | GET | `/api/boards` | Get all boards |
| 10 | GET | `/api/boards/:id` | Get board by ID |
| 11 | PUT | `/api/boards/:id` | Update board |
| 12 | DELETE | `/api/boards/:id` | Delete board |
| 13 | GET | `/api/boards/:id/tasks` | Get tasks by board |

### Tasks (protected)
| # | Method | URL | Description |
|---|--------|-----|-------------|
| 14 | POST | `/api/tasks` | Create task |
| 15 | GET | `/api/tasks` | Get all tasks |
| 16 | GET | `/api/tasks/:id` | Get task by ID |
| 17 | PUT | `/api/tasks/:id` | Update task |
| 18 | DELETE | `/api/tasks/:id` | Delete task |
| 19 | PATCH | `/api/tasks/:id/status` | Update task status only |

### Utility
| # | Method | URL | Description |
|---|--------|-----|-------------|
| 20 | POST | `/notify` | Trigger manual notification |
| — | GET | `/ping` | Health check |

### Notification Service
| Method | URL | Description |
|--------|-----|-------------|
| GET | `/health` | Health check |
| POST | `/notify` | Receive notification |
| GET | `/notifications` | Get all notifications |

## 🔐 Authentication

All `/api/*` routes require the JWT token in the header:
```
Authorization: Bearer <token>
```

## 🧪 Running Tests

```bash
make test
# or:
cd main_service && go test ./tests/... -v
```

10 tests covering:
1. Register success
2. Register duplicate email
3. Register missing fields
4. Login wrong password
5. Login success (JWT returned)
6. Get tasks unauthorized
7. Create board missing title
8. Create board success
9. Create task success
10. Update task status invalid value

## 🏗️ Microservices Communication

`main_service` calls `notification_service` via **Resty v2** HTTP client every time a task is created, updated, or deleted. The call is non-blocking (goroutine).

```go
// notification_client.go
resp, err := restyClient.R().
    SetHeader("Content-Type", "application/json").
    SetBody(map[string]string{"message": message}).
    Post(getNotificationURL() + "/notify")
```

## 🗄️ Database Schema

```sql
users    (id, name, email, password, created_at)
boards   (id, title, description, user_id → users, created_at)
tasks    (id, title, content, status, priority, board_id → boards, user_id → users, created_at, updated_at)
```
