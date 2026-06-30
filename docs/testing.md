# Testing And Validation

This project uses Go tests for backend behavior and Vitest with Testing Library for frontend behavior.

## Test Types

Backend unit tests:

- run with `go test ./...`;
- use in-memory repositories for service behavior;
- do not require PostgreSQL.

Backend integration tests:

- live next to repository implementations as `*_integration_test.go`;
- require `TEST_DATABASE_URL`;
- skip when `TEST_DATABASE_URL` is not set;
- expect the target database to be migrated.

Frontend tests:

- run with `npm test` in `frontend/`;
- use jsdom, Testing Library, and a test `localStorage` implementation;
- mock `fetch` for API behavior.

## Validation Commands

Fast local checks:

```bash
gofmt -l backend
cd backend && GOCACHE=/private/tmp/marketplace-go-cache GOMODCACHE=/private/tmp/marketplace-go-mod go test ./...
make frontend-typecheck
make frontend-test
```

Repository validation:

```bash
make validate
```

`make validate` must not modify tracked source files or database data, but it may generate ignored build artifacts such as `frontend/dist`.

Full local validation with PostgreSQL:

```bash
make frontend-install
make dev-db
make wait-db
make reset-dev-db
make validate
```

`make reset-dev-db` is destructive for the local `marketplace_dev` database.

## Formatting

`make lint` must remain read-only. It checks Go formatting and frontend type safety.

Use this only when intentionally rewriting Go files:

```bash
make format
```

## Migration Behavior

Migrations live under `backend/migrations`.

- `make migrate-up` applies the current schema.
- `make migrate-down` rolls it back.
- `make reset-dev-db` runs down and then up with `ON_ERROR_STOP=1`.

Tests that require PostgreSQL should not assume an unmigrated or empty database unless the validation command explicitly reset it first.

## Test Isolation Expectations

Integration tests that create fixed database records must clean those records before and after the test. Randomized records should still be cleaned up when practical.

The current approach intentionally uses a shared local test database. Do not add per-test containers, per-test schemas, or transaction harnesses unless a future task accepts that scope.

## Known Gaps

These Phase 1 audit items are known but are not part of the Phase 2 enablement fix:

- broader frontend coverage for register, create, and list happy paths;
- dedicated invalid and expired JWT tests;
- malformed localStorage hardening;
- more advanced integration test isolation;
- migration tooling with version tracking.
