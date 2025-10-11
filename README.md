# DMS (Delivery Management System)

A backend service for managing users, authentication, and orders, built with Go, PostgreSQL, and Redis.

## Features

- User signup and login with JWT authentication
- Role-based access (customer, admin)
- Order creation, status tracking, and cancellation
- Admin endpoints for managing all orders
- Dockerized for easy deployment

## Project Structure

```
.
├── cmd/                 # Application entry point
├── internal/            # Application logic
│   ├── migrations/      # Database migration SQL
│   ├── tests/           # Unit and integration tests
├── tests/               # API and service tests
├── Dockerfile           # Docker build instructions
├── docker-compose.yml   # Multi-service orchestration
├── go.mod, go.sum       # Go dependencies
```

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (for local run)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### Running with Docker Compose

```sh
docker-compose up --build
```

- The backend will be available at `http://localhost:8080`
- PostgreSQL at `localhost:5432` (user: `postgres`, password: `MyDb!5432`, db: `dms`)
- Redis at `localhost:6379`

### Running Locally with Go

First, ensure Postgres and Redis are running (can use Docker Compose for just those services):

```sh
docker-compose up -d postgres redis
```

Then, in another terminal, run the backend:

```sh
go run ./cmd
```

### Database Migration

The initial schema is in `internal/migrations/init.sql`. The backend will expect the tables defined there.

## API Endpoints

- `POST /signup` — Register a new user (`email`, `password`, `role`)
- `POST /login` — Obtain JWT token
- `POST /orders` — Create a new order (customer only)
- `POST /orders/{id}/cancel` — Cancel an order (customer only)
- `GET /admin/orders` — List all orders (admin only)
- `GET /orders/{id}/track` — Track order status

## Environment Variables

Set in `docker-compose.yml` or `.env`:

- `PG_DSN` — PostgreSQL DSN
- `REDIS_ADDR` — Redis address
- `JWT_SECRET` — Secret for JWT signing

## Running Tests

```sh
go test ./... -v
```

## License

MIT
