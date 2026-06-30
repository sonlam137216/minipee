# Agent Guide

This repository is a small marketplace MVP. Keep agent work scoped, boring, and verifiable.

## Repository Layout

- `backend/`: Go REST API modular monolith.
- `backend/internal/auth/`: seller registration, login, JWT issuing and verification, seller request context, and seller persistence.
- `backend/internal/products/`: seller-owned draft product rules, HTTP handlers, and product persistence.
- `backend/migrations/`: PostgreSQL schema migrations.
- `frontend/`: React, TypeScript, and Vite seller UI.
- `docs/`: durable project documentation.
- `tasks/`: lightweight task lifecycle records.
- `evidence/`: guidance for validation evidence and task reports.

## Setup And Validation

- Install frontend dependencies with `make frontend-install`.
- Start PostgreSQL with `make dev-db`, then confirm readiness with `make wait-db`.
- Apply migrations with `make migrate-up`.
- Use `make reset-dev-db` only for local validation when it is acceptable to destroy local `marketplace_dev` data.
- Run validation with `make validate`; it must not modify tracked source files or database data, but may generate ignored build artifacts such as `frontend/dist`.
- Use `make format` only when intentionally rewriting Go formatting.

## Go Conventions

- Keep backend code under `backend/internal`.
- Put HTTP translation in handlers, business rules in services, SQL in repositories, and shared HTTP helpers in `httpapi`.
- Keep module boundaries explicit; do not reach across packages for another module's storage details.
- Format Go with `gofmt` through `make format`; verification should use read-only checks.

## React And TypeScript Conventions

- Keep route composition in `frontend/src/app`.
- Keep feature UI under `frontend/src/features`.
- Keep API types and request helpers in `frontend/src/shared/api.ts`.
- Preserve strict TypeScript and test user-visible behavior with Testing Library.

## Modular-Monolith Boundaries

- `auth` owns sellers, password hashes, JWTs, and seller identity in request context.
- `products` owns products and seller-scoped product access.
- Cross-module behavior should go through explicit service or repository contracts.

## Authentication And Authorization Safety

- Protected backend routes must use seller JWT middleware.
- Seller-owned reads and writes must include seller ID in the authorization path.
- Frontend auth state is only UI state; backend authorization must not rely on it.

## Database Migrations

- Every schema change needs an up migration and a down migration.
- Tests that touch PostgreSQL must run against a migrated schema.
- Do not make ad hoc database changes outside migrations.

## Scope Discipline

Do not add orders, payments, inventory, shipping, product publishing, buyer storefronts, microservices, AI infrastructure, orchestration, automated deployment, automated push, or automated merge behavior unless an accepted task explicitly changes scope.

## Test Expectations

- Add unit tests for business rules.
- Add integration tests for SQL behavior, ownership checks, and security-sensitive database paths.
- Add frontend tests when changing critical user flows, auth state, routing, or API error handling.

## Definition Of Done

- Changes are scoped to the task.
- Validation commands were run and reported.
- Database state for integration validation is clear.
- Evidence is recorded in the task or completion report.
- Remaining risks or skipped checks are stated plainly.
- Completion reports must cover: what changed, why it changed, files affected, commands executed, test and validation results, acceptance-criteria mapping, assumptions, remaining risks, deferred work, and evidence produced.
