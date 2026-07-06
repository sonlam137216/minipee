# FE-001 Plan

## Goal

Recover safely when frontend auth `localStorage` contains malformed, partial, or invalid state, without changing backend behavior or marketplace scope.

## Implementation

- Add app-level regression tests in `frontend/src/app/App.test.tsx` for malformed JSON, missing token, missing seller data, invalid seller shape, valid stored auth, and logout.
- Keep login storage behavior covered by the existing login flow test, extended to assert the serialized seller is stored.
- Update `frontend/src/features/auth/AuthContext.tsx` so stored token and seller data are read and validated together.
- Clear both `marketplace.accessToken` and `marketplace.seller` when stored auth is malformed, partial, or invalid.
- Treat invalid stored auth as unauthenticated so protected routes redirect to login and no protected API request is sent during recovery.

## Scope Controls

- No backend behavior changes.
- No database or migration changes.
- No new marketplace endpoints, routes, or features.
- No new production dependencies.
- No token format validation in the frontend; backend JWT verification remains authoritative.
