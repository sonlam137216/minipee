# SEC-001: Add automated JWT rejection tests

## Goal

Add focused automated backend tests proving that protected seller authentication rejects invalid, malformed, expired, and missing JWT access tokens while continuing to accept a correctly signed, non-expired seller token.

This task converts the Phase 1 manual JWT smoke-test coverage into deterministic regression coverage without changing production authentication behavior unless the tests reveal a confirmed defect.

## Scope

- Inspect the current JWT issuing and verification implementation in `backend/internal/auth/jwt.go`.
- Inspect the seller authentication middleware in `backend/internal/auth/middleware.go`.
- Add focused Go tests at the most meaningful existing boundary, preferably `JWTManager.Verify` and/or `RequireSeller`.
- Verify invalid-signature rejection.
- Verify malformed-token rejection.
- Verify expired-token rejection.
- Verify valid-token acceptance.
- Verify missing authentication token rejection if existing coverage is not already adequate.
- Verify authentication rejection responses remain generic and do not expose JWT parser or verification details.
- Update `docs/testing.md` only if the real test command or prerequisite changes.
- Add concise task evidence under `evidence/SEC-001/`.

## Acceptance Criteria

- A correctly signed, non-expired seller token is accepted.
- A token signed with a different secret is rejected.
- A malformed token is rejected.
- An expired token is rejected.
- A missing authentication token is rejected if this behavior is not already adequately covered.
- Rejection does not expose internal JWT parsing or verification details.
- Tests do not rely on sleep-based expiration timing.
- Existing authentication tests continue to pass.
- `make validate` passes.
- No production behavior changes are made unless a test reveals a confirmed defect.
- No marketplace feature or API endpoint is added.

## Explicit Exclusions

- Do not change the authentication API.
- Do not add refresh tokens.
- Do not change the JWT signing algorithm.
- Do not change token storage.
- Do not change seller registration or login behavior.
- Do not add role-based authorization.
- Do not add frontend authentication features.
- Do not add new production dependencies.
- Do not refactor unrelated authentication code.
- Do not add product, order, payment, or marketplace feature changes.
- Do not add repository skills.
- Do not add orchestrator code.
- Do not perform autonomous Git push or merge.

## Current Context

- `docs/security-boundaries.md` documents JWT creation and validation as an implemented control.
- `docs/security-boundaries.md` also lists dedicated invalid and expired JWT tests as a known gap.
- `docs/testing.md` says backend unit tests run with `go test ./...` and do not require PostgreSQL.
- `docs/testing.md` says backend integration tests require `TEST_DATABASE_URL`, skip when it is not set, and expect a migrated schema.
- `docs/agent-workflow.md` requires scoped changes, validation evidence, and concise completion reporting.
- `backend/internal/auth/jwt.go` defines `JWTManager` with HS256 signing, expiration-required verification, seller ID claims, and an injectable `now func() time.Time`.
- `backend/internal/auth/middleware.go` defines `RequireSeller`, which requires `Authorization: Bearer <token>`, verifies the token, stores seller ID in request context, and returns `401` with `unauthenticated` / `Authentication required` on missing or invalid tokens.
- `backend/internal/app/server.go` protects `/api/v1/seller/products` routes with `auth.RequireSeller(jwtManager)`.
- Existing auth tests in `backend/internal/auth/service_test.go` cover registration, login, password hashing behavior, and incorrect password failure.
- Existing auth integration tests in `backend/internal/auth/postgres_integration_test.go` cover PostgreSQL-backed registration and login behavior.
- No dedicated JWT verification or seller-authentication middleware test file exists yet.

## Implementation Notes

- Do not begin from HTTP product behavior unless middleware or JWT service tests cannot cover the acceptance criteria cleanly.
- Prefer a new focused unit test file in `backend/internal/auth`, such as `jwt_test.go` or `middleware_test.go`, following the current Go package test style.
- Use deterministic clock control through the existing `NewJWTManager(secret, ttl, now)` constructor.
- Do not use `time.Sleep` to test token expiration.
- For expired-token coverage, issue a token with a fixed clock and verify it with a manager whose clock or token TTL makes the token expired deterministically. If the current JWT library validation does not consult the injected clock during verification, identify the smallest coherent production change needed before changing behavior.
- For invalid-signature coverage, issue a token with one secret and verify it with a manager using a different secret.
- For malformed-token coverage, use a deliberately non-JWT string and assert the public error path only.
- For middleware coverage, use `httptest` with a minimal protected handler that asserts `SellerIDFromContext` is populated only on valid tokens.
- Do not log or persist full access tokens in evidence. If a token must be referenced, describe it by scenario only, such as "token signed with alternate secret".
- Keep test assertions on public behavior: accepted request reaches the protected handler; rejected requests return `401` with generic `unauthenticated` / `Authentication required`; verification failures return `ErrInvalidAccessToken` where testing the JWT service directly.
- Do not duplicate the same rejection cases across JWT service, middleware, and full route tests unless the second layer proves a distinct behavior.

## Business Rules

- Protected seller routes require a valid seller access token.
- Invalid tokens must not authenticate a seller.
- Expired tokens must not authenticate a seller.
- Malformed tokens must not authenticate a seller and must not panic.
- Valid seller tokens must authenticate as the token seller.
- No marketplace business feature behavior changes are intended.

## Security Considerations

- Tests must not log JWT secrets.
- Tests must not log plaintext passwords or password hashes.
- Evidence must not contain complete reusable access tokens.
- Rejection responses must not expose internal JWT parsing, signing, expiration, or verification details.
- Frontend auth state is irrelevant to this task; backend authorization remains the source of truth.
- Seller-owned protected routes must continue to depend on backend JWT middleware, not client-side state.

## Required Tests

- Add backend Go tests using existing conventions.
- Cover valid-token acceptance.
- Cover wrong-secret or invalid-signature rejection.
- Cover malformed-token rejection.
- Cover expired-token rejection.
- Cover missing-token rejection if existing coverage is not already adequate.
- Cover generic rejection response behavior if testing middleware.
- Do not add frontend tests; no frontend behavior is in scope.
- Do not add migration tests; no schema change is in scope.
- Do not add PostgreSQL integration tests unless the chosen boundary unexpectedly requires database behavior. JWT and middleware behavior should be testable without PostgreSQL.

## Validation Plan

- Run targeted backend auth tests first:

```bash
cd backend && GOCACHE=/private/tmp/marketplace-go-cache GOMODCACHE=/private/tmp/marketplace-go-mod go test ./internal/auth
```

- Run all backend tests:

```bash
cd backend && GOCACHE=/private/tmp/marketplace-go-cache GOMODCACHE=/private/tmp/marketplace-go-mod TEST_DATABASE_URL="postgres://marketplace:marketplace@localhost:5432/marketplace_dev?sslmode=disable" go test ./...
```

- Run repository validation:

```bash
make validate
```

- PostgreSQL should be running and migrated for full validation if integration tests are expected to run instead of skip:

```bash
make dev-db
make wait-db
make migrate-up
```

- Do not run `make reset-dev-db` unless destructive local database reset is explicitly acceptable for validation.

## Definition Of Done

- Implementation remains within this task scope.
- Automated tests cover the required JWT behavior.
- Tests are deterministic and do not sleep for expiration.
- Existing authentication tests continue to pass.
- `make validate` passes.
- The complete diff receives independent review.
- Valid review findings are resolved.
- Acceptance criteria are mapped to evidence.
- Assumptions and remaining risks are documented.
- Human approval is still required for commit and merge.
- Follow-up work is separated from this task.

## Evidence

- Create concise evidence under `evidence/SEC-001/`.
- Include, where relevant:
  - `plan.md`
  - `acceptance-criteria.md`
  - `test-summary.md`
  - `validation-summary.md`
  - `review-report.md`
  - `completion-summary.md`
- Commands run:
  - Record exact commands executed.
- Results:
  - Record pass, fail, skipped, or not-run result for each command.
- Database state:
  - Record whether no database was required, integration tests skipped due to missing `TEST_DATABASE_URL`, or a migrated PostgreSQL database was used.
- Skipped checks:
  - Record any check that could not be run and why.
- Do not store:
  - JWT secrets.
  - Passwords.
  - Password hashes.
  - Complete reusable access tokens.
  - Large generated logs.

## Completion Record

- Change summary:
- Reason for change:
- Files changed:
- Commands executed:
- Test and validation results:
- Acceptance criteria result:
- Assumptions:
- Required tests result:
- Residual risks:
- Deferred work:
- Evidence produced:

## Review Notes

- Risks:
  - Expired-token behavior must be deterministic. Do not introduce sleep-based tests.
  - If verification-time clock injection is missing for expiration validation, document the confirmed defect and implement only the smallest needed production change.
  - Keep evidence sanitized; do not paste full tokens or secrets.
- Follow-up tasks:
  - None identified at task creation time.
