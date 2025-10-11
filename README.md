DMS (Delivery Management System)

A backend service for managing users, authentication, and orders, built with Go, PostgreSQL, and Redis.

Features
User signup and login with JWT authentication
Role-based access (customer, admin)
Order creation, status tracking, and cancellation
Admin endpoints for managing all orders
Dockerized for easy deployment

Project Structure
```.
├── cmd/                 # Application entry point
├── internal/            # Application logic
│   ├── migrations/      # Database migration SQL
│   ├── tests/           # Unit and integration tests
├── tests/               # API and service tests
├── Dockerfile           # Docker build instructions
├── docker-compose.yml   # Multi-service orchestration
├── go.mod, go.sum       # Go dependencies
```

Getting Started
Prerequisites
Docker
Docker Compose
Running with Docker Compose
```
docker-compose up --build
```

The backend will be available at http://localhost:8080
PostgreSQL at localhost:5432 (user: postgres, password: MyDb!5432, db: dms)
Redis at localhost:6379
Database Migration
The initial schema is in init.sql. The backend will expect the tables defined there.

API Endpoints
POST /signup — Register a new user (email, password, role)
POST /login — Obtain JWT token
POST /orders — Create a new order (customer only)
POST /orders/{id}/cancel — Cancel an order (customer only)
GET /admin/orders — List all orders (admin only)
GET /orders/{id}/track — Track order status

Environment Variables
Set in docker-compose.yml or .env:

PG_DSN — PostgreSQL DSN
REDIS_ADDR — Redis address
JWT_SECRET — Secret for JWT signing

Running Tests
```
go test ./... -v
```


Copy and save this as README.md in your project root. Let me know if you want to customize any section!

