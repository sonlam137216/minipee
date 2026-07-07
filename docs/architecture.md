# Architecture

This repository contains a Phase 1 marketplace foundation: a Go backend modular monolith, a PostgreSQL database, and a separately runnable React seller frontend.

## Backend

The backend entrypoint is `backend/cmd/api/main.go`. It loads configuration, creates the application server, starts HTTP serving, and handles graceful shutdown.

`backend/internal/app/server.go` wires the application together:

- creates a `pgxpool` database pool;
- creates shared clock and ID functions;
- wires `auth` and `products` services to PostgreSQL repositories;
- configures Chi routes, request logging, CORS, health checks, and seller JWT protection.

Current API surface:

- `GET /health`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/products`
- `GET /api/v1/products/{productID}`
- `POST /api/v1/seller/products`
- `GET /api/v1/seller/products`
- `GET /api/v1/seller/products/{productID}`
- `POST /api/v1/seller/products/{productID}/publish`

## Backend Modules

`auth` owns seller identity:

- seller registration and login rules;
- password hashing and credential checks;
- JWT issue and verification;
- seller ID storage in request context;
- `sellers` table persistence.

`products` owns seller products and public published catalog reads:

- draft product creation rules;
- draft-to-published transition rules;
- product list and detail retrieval for the authenticated seller;
- published product list and detail retrieval for public catalog users;
- product validation;
- `products` table persistence.

`httpapi` owns shared JSON and error response helpers. `config` owns required environment loading. `platform/id` owns UUID generation.

## Frontend

The frontend is a React and TypeScript Vite app under `frontend/`.

- `frontend/src/app`: application shell, routes, and protected route behavior.
- `frontend/src/features/auth`: seller login, registration, and browser auth state.
- `frontend/src/features/products`: seller product list, create, detail, and public catalog pages.
- `frontend/src/shared/api.ts`: API types, request helpers, and API error parsing.
- `frontend/src/test`: Vitest and Testing Library setup.

The frontend stores the access token and seller snapshot in `localStorage` for local MVP development. This is UI state only; authorization remains a backend responsibility.

## Data Ownership

PostgreSQL uses one local database and one schema. Ownership is by backend module:

- `auth` owns `sellers`.
- `products` owns `products`.

`products.seller_id` references `sellers.id` with `ON DELETE CASCADE`. Product queries that read seller-owned product data must include the authenticated seller ID.

## Authentication And Authorization Boundaries

Public auth routes issue JWT access tokens after registration or login.

Seller product routes are protected by `auth.RequireSeller`, which:

- requires an `Authorization: Bearer <token>` header;
- verifies the JWT using the configured secret;
- stores the seller ID in request context.

Product handlers read the seller ID from context. Product repository methods enforce ownership by including `seller_id` in seller product queries. Cross-seller seller product access and publish attempts return not found. Public product routes are unauthenticated and only return rows with status `published`.

## Current Request Flows

Seller registration:

1. Frontend posts email, password, and display name.
2. Auth handler decodes JSON and calls the auth service.
3. Auth service normalizes email, hashes the password, creates the seller, and issues a JWT.
4. Frontend stores the JWT and seller snapshot in localStorage.

Seller login:

1. Frontend posts email and password.
2. Auth service normalizes the email, loads the seller, checks bcrypt, and issues a JWT.
3. Frontend stores the JWT and seller snapshot.

Draft product creation:

1. Frontend posts name and description with a bearer token.
2. Middleware verifies the token and places seller ID in context.
3. Product service validates name and forces status to `draft`.
4. Product repository inserts a row owned by the seller.

Seller product read:

1. Frontend requests list or detail with a bearer token.
2. Middleware establishes seller identity.
3. Product repository queries by seller ID; detail lookup returns owned draft or published products.

Product publish:

1. Frontend posts to the publish endpoint with a bearer token.
2. Middleware establishes seller identity.
3. Product service loads the seller-owned product, rejects already published products, and publishes only draft rows.
4. Product repository updates status to `published` and refreshes `updated_at`.

Public catalog read:

1. Frontend requests catalog list or detail without a bearer token.
2. Product repository filters by `status = 'published'`.
3. Public product responses omit seller ID and auth-owned fields.
