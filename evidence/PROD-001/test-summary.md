# PROD-001 Test Summary

## Red Checks

- `cd backend && GOCACHE=/private/tmp/marketplace-go-cache GO111MODULE=on go test ./internal/products`
  - Failed before implementation with undefined `StatusPublished`, missing service methods, and missing handler methods.
- `cd frontend && npm test -- App.test.tsx`
  - Failed before implementation because publish button and `/catalog` routes were missing.

## Passing Checks

- `cd backend && GOCACHE=/private/tmp/marketplace-go-cache GOMODCACHE=/private/tmp/marketplace-go-path/pkg/mod TEST_DATABASE_URL='postgres://marketplace:marketplace@localhost:5432/marketplace_dev?sslmode=disable' go test ./...`
  - Passed after applying migration `000002` to local dev DB.
- `make frontend-typecheck`
  - Passed.
- `make frontend-test`
  - Passed, 1 test file and 12 tests.
- `gofmt -l backend`
  - Passed with no output.
- `make validate`
  - Passed.

## Notes

The original temp module cache path `/private/tmp/marketplace-go-mod` did not load packages with this Go toolchain. The Makefile now uses `/private/tmp/marketplace-go-path/pkg/mod`, which was verified with `make validate`.

## Follow-Up Test Isolation Fix

Observed validation failure:

- `cd backend && GOCACHE=/private/tmp/marketplace-go-cache GOMODCACHE=/private/tmp/marketplace-go-path/pkg/mod TEST_DATABASE_URL="postgres://marketplace:marketplace@localhost:5432/marketplace_dev?sslmode=disable" go test ./internal/products -count=1`
  - Failed in `TestPostgresProductRepository`.
  - Failure showed public list contained the test's published fixture plus an older smoke product: `85caf760-d87d-4c68-a896-f6ac542b445c`.

Root cause:

- The integration test used the shared development database but asserted `ListPublishedProducts` returned exactly one product.
- Public list production behavior is to return all published products, so the test was assuming global database isolation that the documented setup does not provide.

Fix:

- Updated `backend/internal/products/postgres_integration_test.go` only.
- Test fixtures now use run-unique seller IDs, emails, and product IDs.
- Cleanup deletes only rows belonging to the test's seller IDs.
- Public list assertion now verifies:
  - the test's published product is included;
  - the test's draft product is not included;
  - every returned public product has status `published`;
  - no exact global count is assumed.

Verification after fix:

- Focused products test passed once:
  - `ok marketplace/backend/internal/products 0.406s`
- Focused products test passed a second time:
  - `ok marketplace/backend/internal/products 0.289s`
- `make validate` exited `0`.

Production behavior:

- No production code or API behavior was changed for this follow-up fix.
