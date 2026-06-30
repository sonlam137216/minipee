# Domain Map

This document describes the domains implemented in the current repository. It does not treat future marketplace capabilities as current architecture.

## Auth

Package: `backend/internal/auth`

Responsibilities:

- Register sellers.
- Normalize seller email addresses.
- Hash seller passwords with bcrypt.
- Verify login credentials.
- Issue and verify seller JWT access tokens.
- Store and read seller identity in request context.
- Persist seller records through `auth.PostgresRepository`.

Data and concepts owned:

- `sellers` table.
- Seller ID, email, display name, password hash, creation time, and update time.
- JWT claims for seller identity.
- Authentication errors such as invalid credentials and duplicate email.

Interfaces or modules it may call:

- `httpapi` for JSON request and error responses from auth handlers.
- `platform/id` indirectly through service wiring for seller IDs.
- PostgreSQL through `auth.PostgresRepository`.
- bcrypt and JWT libraries for password and token behavior.

Things it must not own:

- Product storage or product business rules.
- Seller product ownership queries.
- Buyer accounts, orders, payments, inventory, shipping, product publishing, or storefront behavior.

Important business rules:

- Emails are trimmed and lowercased before persistence or lookup.
- Registration requires a non-empty display name and a password of at least 8 characters.
- Passwords are stored as bcrypt hashes, not plaintext.
- Duplicate seller email registration returns the auth duplicate-email error.
- Login returns invalid credentials for unknown emails and incorrect passwords.

Authentication and authorization responsibilities:

- `auth.JWTManager` issues HS256 access tokens with seller ID, email, issued-at, expiry, and subject claims.
- `auth.RequireSeller` verifies bearer tokens and writes the seller ID to request context.
- Auth owns authentication. It does not authorize product ownership beyond establishing the authenticated seller identity.

## Products

Package: `backend/internal/products`

Responsibilities:

- Create seller-owned draft products.
- Validate draft product input.
- List products for the authenticated seller.
- Read one authenticated seller draft product by ID.
- Persist product records through `products.PostgresRepository`.

Data and concepts owned:

- `products` table.
- Product ID, seller ID, name, description, status, creation time, and update time.
- Draft product status and product validation errors.

Interfaces or modules it may call:

- `auth.SellerIDFromContext` in HTTP handlers to read authenticated seller identity.
- `httpapi` for JSON and error responses from product handlers.
- PostgreSQL through `products.PostgresRepository`.
- `platform/id` indirectly through service wiring for product IDs.

Things it must not own:

- Passwords, password hashes, JWT creation, JWT verification, or seller credential checks.
- Seller account persistence except the `seller_id` foreign-key relationship already present on products.
- Product publishing, buyer discovery, carts, orders, payments, inventory, shipping, or storefront behavior.

Important business rules:

- Product names must contain between 3 and 200 characters after trimming.
- New products are always created with status `draft`.
- Client-provided status cannot override draft creation.
- Seller product list queries are scoped by authenticated seller ID.
- Product detail lookup requires product ID, authenticated seller ID, and draft status.
- Cross-seller draft product access returns not found.

Authentication and authorization responsibilities:

- Product routes are mounted behind `auth.RequireSeller`.
- Product handlers must obtain the seller ID from request context.
- Product repository reads that expose seller-owned data must include `seller_id` in the SQL predicate.

## Shared HTTP And Platform Infrastructure

Packages:

- `backend/internal/httpapi`
- `backend/internal/config`
- `backend/internal/platform/id`
- `backend/internal/app`

Responsibilities:

- `httpapi` provides JSON response, JSON error response, and JSON request decoding helpers.
- `config` loads required environment variables and parses JWT expiration.
- `platform/id` generates UUID values.
- `app` wires repositories, services, handlers, routes, middleware, CORS, logging, database connection, and graceful server ownership.

Data and concepts owned:

- Shared HTTP error response shape.
- Runtime configuration values.
- UUID generation utility.
- Application wiring, not domain business data.

Interfaces or modules it may call:

- `app` calls `auth` and `products` constructors to wire the modular monolith.
- `httpapi` is called by HTTP handlers across domains.

Things it must not own:

- Auth business rules.
- Product business rules.
- Table-specific SQL ownership.

Authentication and authorization responsibilities:

- `app` mounts seller product routes inside the JWT-protected route group.
- `httpapi` only formats responses; it does not decide authentication or ownership.

## Deferred Marketplace Domains

The following domains are not implemented as current architecture:

- Product publishing.
- Buyer storefront or buyer product discovery.
- Carts.
- Orders.
- Payments.
- Inventory.
- Shipping.
- Promotions.
- Reviews.
- Admin tooling.
- Microservices.
- AI runtime behavior.
- Orchestration or autonomous agent infrastructure.
