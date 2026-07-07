# PROD-001 Completion Summary

## What Changed

- Added `published` product status and repeat-publish conflict error.
- Added seller publish service, repository, handler, and route.
- Added public published-product list/detail service, repository, handlers, and routes.
- Added migration `000002` to allow published status and index public catalog ordering.
- Added public catalog frontend routes `/catalog` and `/catalog/:productID`.
- Added seller detail publish button for draft products.
- Updated product API types to separate seller and public product responses.
- Updated docs for architecture, domain map, data ownership, security boundaries, and product scope.
- Updated Makefile Go module cache path to a GOPATH-shaped temp path that validates with this Go toolchain.

## Files Affected

- `Makefile`
- `backend/internal/app/server.go`
- `backend/internal/products/model.go`
- `backend/internal/products/service.go`
- `backend/internal/products/postgres.go`
- `backend/internal/products/handler.go`
- `backend/internal/products/service_test.go`
- `backend/internal/products/postgres_integration_test.go`
- `backend/internal/products/handler_test.go`
- `backend/migrations/000002_publish_products.up.sql`
- `backend/migrations/000002_publish_products.down.sql`
- `frontend/src/shared/api.ts`
- `frontend/src/app/App.tsx`
- `frontend/src/app/App.test.tsx`
- `frontend/src/features/products/ProductDetailPage.tsx`
- `frontend/src/features/products/ProductListPage.tsx`
- `frontend/src/features/products/PublicProductListPage.tsx`
- `frontend/src/features/products/PublicProductDetailPage.tsx`
- `docs/architecture.md`
- `docs/domain-map.md`
- `docs/data-ownership.md`
- `docs/security-boundaries.md`
- `docs/product-scope.md`
- `evidence/PROD-001/*`

## Commands Executed

- `cd backend && gofmt -w internal/products internal/app/server.go`
- `cd backend && GOCACHE=/private/tmp/marketplace-go-cache GO111MODULE=on go test ./internal/products`
- `cd frontend && npm test -- App.test.tsx`
- `docker compose exec -T postgres psql -v ON_ERROR_STOP=1 -U marketplace -d marketplace_dev -f /dev/stdin < backend/migrations/000002_publish_products.up.sql`
- `cd backend && GOCACHE=/private/tmp/marketplace-go-cache GOMODCACHE=/private/tmp/marketplace-go-path/pkg/mod TEST_DATABASE_URL='postgres://marketplace:marketplace@localhost:5432/marketplace_dev?sslmode=disable' go test ./...`
- `make frontend-typecheck`
- `make frontend-test`
- `gofmt -l backend`
- `make validate`
- Local API smoke script against `http://localhost:18080`.

## Assumptions

- Public routes use `/catalog` in the frontend to keep seller `/products` routes unchanged.
- Public responses do not need seller identity in this task.
- Repeated publish returns `409` rather than idempotent success.

## Remaining Risks And Deferred Work

- Migration tooling remains file-based without version tracking.
- Down migration cannot safely preserve published rows and fails if any exist.
- Search, pagination, images, buyer accounts, orders, payments, shipping, unpublish, and product editing are deferred.

## Follow-Up: Test Isolation Correction

The product repository integration test was updated after a shared-database validation failure. The previous assertion expected the public published product list to contain exactly one row, but public list behavior intentionally returns all published products. A prior smoke product in `marketplace_dev` made that assertion nondeterministic.

The follow-up changed only `backend/internal/products/postgres_integration_test.go`:

- use run-unique seller IDs, emails, and product IDs;
- cleanup only the test-owned seller/product rows;
- assert public list includes the test's published product;
- assert public list excludes the test's draft product;
- assert every public list row has status `published`;
- stop asserting exact public list length.

No production behavior was changed.

## Evidence Produced

- `plan.md`
- `acceptance-criteria.md`
- `api-examples.md`
- `test-summary.md`
- `validation-summary.md`
- `review-report.md`
- `completion-summary.md`
