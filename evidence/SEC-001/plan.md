# SEC-001 Plan

## Scope

Implement only automated backend tests for JWT authentication acceptance and rejection behavior.

## Approach

- Add focused unit coverage for `JWTManager.Verify`.
- Add focused middleware coverage for `RequireSeller` using `httptest`.
- Use fixed test-only secrets and deterministic issue times.
- Create expired tokens with an expiration timestamp already in the past.
- Avoid logging or recording complete JWTs, secrets, passwords, or hashes.
- Preserve production behavior unless a direct defect blocks the acceptance criteria.

## Validation Plan

- Run `gofmt` on changed Go test files.
- Run focused auth package tests.
- Run repeat focused auth package tests with `-count=1`.
- Run `make validate`.
- Inspect status, diff, dependency files, generated artifacts, and evidence.
