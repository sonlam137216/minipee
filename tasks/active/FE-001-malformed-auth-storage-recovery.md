# FE-001: Recover safely from malformed frontend auth storage

## Goal

Make the seller frontend recover safely when browser `localStorage` contains malformed or invalid authentication state, so the app does not crash during startup and protected routes return the user to the login flow.

This task addresses the documented malformed `localStorage` hardening gap without changing backend authorization, marketplace domain behavior, or the current local-MVP token storage trade-off.

## Scope

- Inspect and update frontend authentication state initialization in `frontend/src/features/auth/AuthContext.tsx`.
- Preserve the existing login, registration, logout, and protected-route flows.
- Add focused frontend tests for malformed and invalid auth storage behavior, likely in `frontend/src/app/App.test.tsx` or a colocated auth test file using the existing Vitest and Testing Library setup.
- Update documentation only if the implemented behavior changes documented setup, validation, or security-boundary text.
- Record concise validation evidence in this task file or under `evidence/FE-001/` if more detail is needed.

## Acceptance Criteria

- If `marketplace.seller` contains malformed JSON, the app renders without throwing.
- If stored auth state is incomplete, invalid, or inconsistent, the frontend clears the stored auth keys and treats the user as unauthenticated.
- If only one of `marketplace.accessToken` or `marketplace.seller` is present, the frontend clears the partial auth state and treats the user as unauthenticated.
- Protected seller routes redirect to the login page after malformed or invalid stored auth state is recovered.
- Valid stored auth state continues to authenticate the frontend UI and allow protected seller pages to render.
- Login and registration still store `marketplace.accessToken` and `marketplace.seller` in the existing format.
- Logout still removes both auth storage keys.
- Recovery does not send malformed, partial, or invalid stored tokens to protected API requests during startup.
- No backend authentication or authorization behavior changes are made.
- No new marketplace feature, API endpoint, dependency, or storage mechanism is added.

## Explicit Exclusions

- Do not change backend JWT verification, seller middleware, or product authorization.
- Do not add refresh tokens, token revocation, session rotation, cookies, or CSRF changes.
- Do not replace `localStorage` with another storage mechanism.
- Do not add orders, payments, inventory, shipping, publishing, storefronts, or other marketplace features.
- Do not add new frontend routes or authentication screens.
- Do not add new production dependencies.
- Do not perform database migrations or ad hoc database changes.
- Do not run destructive database reset commands unless explicitly accepted for validation.
- Do not perform autonomous Git push, merge, or deployment.

## Current Context

- `docs/security-boundaries.md` documents frontend token storage as browser UI state only; backend authorization remains the source of truth.
- `docs/security-boundaries.md` lists `localStorage` token storage as an MVP trade-off and defers more advanced browser storage hardening.
- `docs/testing.md` lists malformed `localStorage` hardening as a known gap.
- `docs/testing.md` says frontend tests run with `npm test` in `frontend/`, using jsdom, Testing Library, and test `localStorage`.
- `docs/agent-workflow.md` requires scoped changes, validation evidence, and concise completion reporting.
- `evidence/README.md` requires acceptance criteria to map to observable evidence and prohibits secrets in evidence.
- `frontend/src/features/auth/AuthContext.tsx` currently initializes `token` from `localStorage.getItem("marketplace.accessToken")`.
- `frontend/src/features/auth/AuthContext.tsx` currently initializes `seller` by directly parsing `localStorage.getItem("marketplace.seller")` with `JSON.parse`.
- Direct parsing means malformed seller storage can throw during `AuthProvider` initialization before the app can recover.
- `frontend/src/app/ProtectedRoute.tsx` currently gates seller routes on `auth.token === null`.
- `frontend/src/app/App.test.tsx` currently covers login storage, protected-route login redirection, and product form validation with valid stored auth state.

## Implementation Notes

- Keep the recovery logic small and local to the frontend auth boundary.
- Prefer a single helper for reading stored auth state so token and seller are validated together.
- Treat frontend auth state as present only when both the access token and a seller object with the expected shape are valid.
- On recovery from malformed, invalid, or partial state, remove both `marketplace.accessToken` and `marketplace.seller`.
- Avoid exposing parser errors or stored token values in UI, logs, tests, or evidence.
- Preserve the existing serialized seller shape from `AuthResponse`.
- Prefer user-visible behavior tests over implementation-only tests: render the app with prepared `localStorage` and assert route behavior, storage cleanup, and absence of protected API calls where relevant.
- Do not over-validate JWT structure in the frontend; backend JWT verification remains authoritative.

## Business Rules

- No marketplace business rules are added or changed.
- Seller UI access remains a frontend convenience only; backend authorization must still rely on seller JWT middleware.
- A browser with invalid saved auth state must be treated as unauthenticated.
- Valid saved seller auth state should continue the current seller UI session behavior.

## Security Considerations

- Frontend stored auth data is untrusted browser state.
- Malformed, partial, or invalid stored auth state must not be used to authorize protected API requests.
- Recovery must not log or record complete JWT access tokens.
- Tests and evidence must not include real secrets, passwords, password hashes, or reusable access tokens.
- Backend authorization remains the security boundary; frontend recovery only prevents crashes and stale UI authentication.

## Required Tests

- Add frontend tests for malformed `marketplace.seller` JSON.
- Add frontend tests for partial auth storage where only the token or only the seller exists.
- Add frontend tests for invalid seller shape or otherwise inconsistent stored auth state.
- Add or preserve a frontend test proving valid stored auth state still reaches a protected seller page.
- Add or preserve tests proving login stores the token and seller in the existing format.
- Add or preserve tests proving unauthenticated protected routes redirect to login.
- Do not add backend tests; no backend behavior is in scope.
- Do not add PostgreSQL integration tests; no schema or database behavior is in scope.
- Do not add migration tests; no schema change is in scope.

## Validation Plan

- Run targeted frontend tests:

```bash
cd frontend && npm test
```

- Run frontend type checking:

```bash
make frontend-typecheck
```

- Run repository validation:

```bash
make validate
```

- PostgreSQL is not required for the frontend-specific tests.
- PostgreSQL should be running and migrated if `make validate` is expected to execute backend integration tests instead of skipping them:

```bash
make dev-db
make wait-db
make migrate-up
```

- Do not run `make reset-dev-db` unless destructive local database reset is explicitly acceptable for validation.

## Definition Of Done

- Acceptance criteria are met.
- Required frontend tests are added or updated.
- Existing login, registration, logout, and protected-route behavior is preserved.
- Required validation commands were run or explicitly skipped with a reason.
- Evidence records exact commands, results, database state, and skipped checks.
- Complete diff receives independent review.
- Valid review findings are resolved.
- Follow-up work is separated from this task.

## Evidence

- Commands run:
- Results:
- Database state:
- Skipped checks:

## Completion Record

- Change summary:
- Reason for change:
- Files changed:
- Commands executed:
- Test and validation results:
- Acceptance criteria result:
- Assumptions:
  - The prompt referenced pasted FE-001 task content, but no concrete task contract text was included after that placeholder. This task record was created from the user-provided objective and repository instructions.
- Required tests result:
- Residual risks:
- Deferred work:
- Evidence produced:

## Review Notes

- Risks:
  - Avoid making frontend validation look like a security guarantee; backend JWT middleware remains authoritative.
  - Avoid clearing valid stored sessions due to overly strict seller validation that rejects the current API response shape.
  - Ensure malformed storage recovery does not trigger protected product API calls before redirecting to login.
- Follow-up tasks:
  - None identified at task creation time.
