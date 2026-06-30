# Code Review Checklist

Use this checklist for pull requests, patches, or agent-produced changes.

## Scope

- The change matches the accepted task.
- No excluded marketplace features were added accidentally.
- No new dependencies were added without a clear validation need.

## Backend

- Handlers translate HTTP input/output only.
- Services own business rules.
- Repositories own SQL and database error translation.
- Errors preserve the existing JSON error response shape.
- New schema behavior has up and down migrations.

## Authentication And Authorization

- Protected routes use seller JWT middleware.
- Seller-owned data paths include authenticated seller ID.
- Cross-seller access is denied at the backend.
- Frontend auth state is not treated as authorization.

## Database

- Module ownership is respected: `auth` owns sellers, `products` owns products.
- Integration tests clean fixed records before and after running.
- Migration commands fail fast on SQL errors.

## Frontend

- Routes remain in `frontend/src/app`.
- Feature UI remains under `frontend/src/features`.
- API wire types remain in `frontend/src/shared/api.ts`.
- User-visible behavior and error states are covered when changed.

## Validation

- Read-only validation was run.
- PostgreSQL-backed validation states whether the database was migrated or reset.
- Any skipped checks are explained.
- Evidence is concise and reproducible.
