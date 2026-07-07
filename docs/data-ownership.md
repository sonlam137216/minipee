# Data Ownership

This document maps the current database schema to backend module ownership.

## Current Tables

The current schema starts in `backend/migrations/000001_init.up.sql` and is extended by `backend/migrations/000002_publish_products.up.sql`.

| Table | Owning module | Repository | Write ownership |
| --- | --- | --- | --- |
| `sellers` | `backend/internal/auth` | `auth.PostgresRepository` | Seller registration writes through `CreateSeller`. |
| `products` | `backend/internal/products` | `products.PostgresRepository` | Draft product creation writes through `Create`; publish writes through `PublishDraftBySeller`. |

There are no current tables for buyers, orders, payments, inventory, shipping, carts, or AI features.

## Table Details

### `sellers`

Columns:

- `id UUID PRIMARY KEY`
- `email TEXT NOT NULL UNIQUE`
- `password_hash TEXT NOT NULL`
- `display_name TEXT NOT NULL`
- `created_at TIMESTAMPTZ NOT NULL`
- `updated_at TIMESTAMPTZ NOT NULL`

Owned by `auth`.

Allowed writes:

- `auth.PostgresRepository.CreateSeller` inserts sellers during registration.

Allowed reads:

- `auth.PostgresRepository.FindSellerByEmail` for login.
- `auth.PostgresRepository.FindSellerByID` for seller lookup.

Sensitive fields:

- `password_hash` must remain inside persistence and authentication boundaries.
- API responses use `sellerResponse`, which excludes `password_hash`.
- Plaintext passwords are accepted only as request input for registration and login and must not be persisted or logged.

### `products`

Columns:

- `id UUID PRIMARY KEY`
- `seller_id UUID NOT NULL REFERENCES sellers(id) ON DELETE CASCADE`
- `name TEXT NOT NULL`
- `description TEXT NOT NULL DEFAULT ''`
- `status TEXT NOT NULL CHECK (status IN ('draft', 'published'))`
- `created_at TIMESTAMPTZ NOT NULL`
- `updated_at TIMESTAMPTZ NOT NULL`

Indexes:

- `products_seller_created_idx` on `(seller_id, created_at DESC)`
- `products_seller_id_id_draft_idx` on `(seller_id, id)` where `status = 'draft'`
- `products_published_created_idx` on `(created_at DESC, id DESC)` where `status = 'published'`

Owned by `products`.

Allowed writes:

- `products.PostgresRepository.Create` inserts products with a seller ID supplied by the authenticated request path.
- `products.PostgresRepository.PublishDraftBySeller` updates a seller-owned draft product to `published` and refreshes `updated_at`.

Allowed reads:

- `products.PostgresRepository.ListBySeller` lists products where `seller_id = $1`.
- `products.PostgresRepository.FindBySellerID` reads one seller-owned product where `id = $1 AND seller_id = $2`.
- `products.PostgresRepository.ListPublished` lists products where `status = 'published'`.
- `products.PostgresRepository.FindPublishedByID` reads one product where `id = $1 AND status = 'published'`.

## Foreign-Key Relationships

`products.seller_id` references `sellers.id` with `ON DELETE CASCADE`.

The relationship means:

- every product row belongs to one seller row;
- deleting a seller deletes that seller's products at the database level;
- product ownership is represented by `products.seller_id`.

## Seller-To-Product Ownership

The authenticated seller ID comes from `auth.RequireSeller`, is stored in request context, and is read by product handlers through `auth.SellerIDFromContext`.

Product ownership is enforced by:

- route protection in `backend/internal/app/server.go`;
- product handlers requiring a seller ID in context;
- product repository queries including `seller_id` for seller-owned reads;
- product publish using `seller_id` from context and draft-only SQL;
- product creation using the authenticated seller ID rather than trusting a client-supplied owner.

Public catalog reads do not carry seller identity. They must filter by `status = 'published'` and must use response shapes that omit `seller_id` and auth-owned fields.

## Allowed Cross-Module Access

Allowed:

- Product handlers may read seller identity from auth request context.
- The products table may reference `sellers.id` through the foreign key.
- Application wiring may create both auth and products services and repositories.

Not allowed:

- `products` must not read or write `sellers.password_hash`.
- `products` must not validate credentials or issue tokens.
- `auth` must not own product SQL or product business rules.
- Handlers must not bypass services and repositories for domain behavior.

## Ownership Bypass Rules

- Seller-owned product reads and writes must carry authenticated seller ID through the backend path.
- Product SQL that returns seller-owned data must include `seller_id` in the predicate.
- Public product SQL must include `status = 'published'` in the predicate.
- Frontend auth state is not an authorization boundary.
- Clients must not be allowed to choose product ownership.
- Direct database changes outside migrations and repositories are not part of the application workflow.
