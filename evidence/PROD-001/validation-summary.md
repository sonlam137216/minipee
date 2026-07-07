# PROD-001 Validation Summary

## Database State

- Local Docker PostgreSQL was already running.
- Existing local database had migration `000001`.
- Applied only `backend/migrations/000002_publish_products.up.sql` to preserve local data.
- Did not run `make reset-dev-db`.

Migration command result:

```text
ALTER TABLE
ALTER TABLE
CREATE INDEX
```

## Full Validation

Command:

```sh
make validate
```

Result: exited `0`.

Included checks:

- Go formatting check.
- Frontend TypeScript typecheck.
- Backend Go tests against local PostgreSQL.
- Frontend Vitest suite.
- Frontend production build.
- Backend build.

## Follow-Up Validation: Shared Database Test Isolation

Initial failing command:

```sh
cd backend && GOCACHE=/private/tmp/marketplace-go-cache GOMODCACHE=/private/tmp/marketplace-go-path/pkg/mod TEST_DATABASE_URL="postgres://marketplace:marketplace@localhost:5432/marketplace_dev?sslmode=disable" go test ./internal/products -count=1
```

Initial result: exited `1`.

Cause: `TestPostgresProductRepository` expected the public list to contain exactly one published fixture, but the shared database also contained a previously published smoke product.

Post-fix validation:

```sh
gofmt -w backend/internal/products/postgres_integration_test.go
cd backend && GOCACHE=/private/tmp/marketplace-go-cache GOMODCACHE=/private/tmp/marketplace-go-path/pkg/mod TEST_DATABASE_URL="postgres://marketplace:marketplace@localhost:5432/marketplace_dev?sslmode=disable" go test ./internal/products -count=1
cd backend && GOCACHE=/private/tmp/marketplace-go-cache GOMODCACHE=/private/tmp/marketplace-go-path/pkg/mod TEST_DATABASE_URL="postgres://marketplace:marketplace@localhost:5432/marketplace_dev?sslmode=disable" go test ./internal/products -count=1
make validate
```

Results:

- First focused products test: exited `0`.
- Second focused products test: exited `0`.
- `make validate`: exited `0`.

No database reset was used.
