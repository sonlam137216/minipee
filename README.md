# Marketplace MVP

A small marketplace application focused on the first seller product-management slice. The Phase 1 implementation is a modular monolith backend and a separately runnable React frontend.

## Current Scope

Implemented:

- Seller registration and login.
- JWT access-token authentication.
- Create draft product.
- List the authenticated seller's products.
- Retrieve the authenticated seller's own draft product.
- PostgreSQL schema migrations.

Deferred:

- Refresh tokens and production-grade browser token storage.
- Product publishing.
- Buyer storefront.
- Orders, payments, inventory, shipping, reviews and promotions.
- Admin tooling, microservices, Redis, Kafka, Elasticsearch and AI-related infrastructure.

## Prerequisites

- Go 1.23 or newer.
- Node.js 22 or newer.
- npm.
- Docker and Docker Compose.

## Development Guides

- [Agent guide](AGENTS.md): durable rules for scoped, verifiable agent-assisted work.
- [Architecture](docs/architecture.md): current backend, frontend, data ownership and auth boundaries.
- [Testing and validation](docs/testing.md): test types, database prerequisites and validation commands.
- [Agent workflow](docs/agent-workflow.md): task lifecycle, evidence and completion-report expectations.
- [Code review](docs/code-review.md): review checklist for repository changes.

## Environment Setup

Copy the example environment file and set a real local secret:

```bash
cp .env.example .env
```

Variables:

- `APP_PORT`: backend HTTP port.
- `DATABASE_URL`: PostgreSQL connection string.
- `JWT_SECRET`: required signing secret for JWT access tokens.
- `JWT_EXPIRATION_MINUTES`: access-token lifetime.
- `FRONTEND_ORIGIN`: allowed CORS origin.
- `VITE_API_BASE_URL`: frontend API base URL.

The backend validates required environment variables during startup and fails fast if any are missing.

## Start PostgreSQL

```bash
make dev-db
```

## Run Migrations

```bash
make migrate-up
```

Rollback:

```bash
make migrate-down
```

## Run Backend

```bash
make backend-run
```

Backend health check:

```bash
curl -i http://localhost:8080/health
```

## Run Frontend

```bash
make frontend-install
make frontend-dev
```

Open `http://localhost:5173`.

For this MVP, the frontend stores the JWT access token in `localStorage`. This is a deliberate local-development trade-off, not the final production security design.

## Run Tests

```bash
make backend-test
make frontend-test
```

## Run Validation

```bash
make validate
```

Validation runs Go formatting checks, frontend type checks, backend tests, frontend tests and builds for both applications. It must not modify tracked source files or database data, but it may generate ignored build artifacts such as `frontend/dist`.

Full validation requires frontend dependencies and a running, migrated PostgreSQL database:

```bash
make frontend-install
make dev-db
make wait-db
make reset-dev-db
make validate
```

`make reset-dev-db` is destructive for the local `marketplace_dev` database. Use `make format` when you intentionally want to rewrite Go formatting.

## API Endpoints

```text
GET  /health
POST /api/v1/auth/register
POST /api/v1/auth/login
POST /api/v1/seller/products
GET  /api/v1/seller/products
GET  /api/v1/seller/products/{productID}
```

Protected endpoints require:

```text
Authorization: Bearer <access-token>
```

## Example API Workflow

Register:

```bash
curl -s http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"seller@example.com","password":"password123","displayName":"Seller"}'
```

Login:

```bash
curl -s http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"seller@example.com","password":"password123"}'
```

Create a draft product:

```bash
curl -s http://localhost:8080/api/v1/seller/products \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"name":"Draft product","description":"Optional description"}'
```

List seller products:

```bash
curl -s http://localhost:8080/api/v1/seller/products \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```
