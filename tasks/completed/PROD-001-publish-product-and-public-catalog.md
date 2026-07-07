# PROD-001: Publish product and expose public catalog

## Status

Active

## Background

The current marketplace foundation allows a seller to register, log in, create draft products, list their own products and retrieve their own products.

However, products cannot yet be published, and unauthenticated buyers cannot browse public products.

## Why

A marketplace needs a clear transition from seller-owned draft products to publicly visible products.

This task introduces the smallest coherent public catalog feature:

- Seller publishes a draft product.
- Public users can list published products.
- Public users can view published product details.
- Draft products remain private to the owning seller.

## Desired outcome

A seller can publish one of their own draft products.

Once published, the product becomes visible through public catalog endpoints and the frontend public catalog pages.

Draft products must remain hidden from public users.

## In scope

### Backend

- Add product publish behavior.
- Add seller endpoint to publish an owned draft product.
- Add public endpoint to list published products.
- Add public endpoint to retrieve one published product by ID.
- Add service-layer business rules for status transition.
- Add repository queries needed for published products.
- Add backend tests for publish behavior and public visibility.

### Frontend

- Add seller UI action to publish a draft product.
- Add a public product list page.
- Add a public product detail page.
- Keep UI simple and accessible.
- Add focused frontend tests for the new behavior where practical.

### Documentation and evidence

- Update README or docs only if commands, routes or documented current scope change.
- Create evidence under `evidence/PROD-001/`.

## Out of scope

- Product images
- Product variants
- Categories
- Inventory
- Cart
- Orders
- Payments
- Shipping
- Promotions
- Reviews
- Admin moderation
- Search ranking
- Pagination beyond a simple initial implementation
- Seller product editing
- Unpublish behavior
- Buyer accounts
- Wishlist
- Recommendation systems
- Microservices
- Event-driven architecture
- AI orchestrator
- Codex skills
- Autonomous Git push, merge or deployment

## Business rules

- Only the owning seller can publish their own product.
- A seller cannot publish another seller's product.
- Only a `draft` product can transition to `published`.
- Publishing must be performed by the backend, not by trusting client-provided status.
- Public catalog endpoints must only return products with status `published`.
- Draft products must never appear in public list or public detail responses.
- Public users do not need authentication to browse published products.
- Publishing a product should update `updated_at`.
- Product ownership must continue to be enforced on seller endpoints.
- Public responses must not expose private seller authentication data.
- Product status values must remain explicit and validated.

## Product status policy

The product status lifecycle for this task is:

```text
draft → published
```

No other transitions are introduced in this task.

Repeated publish behavior:

- Publishing an already published product should return a clear client error, such as `409 Conflict`, with a safe public error response.
- It must not create duplicate side effects.
- It must not change ownership or create a new product.

## Acceptance criteria

1. An authenticated seller can publish their own draft product.
2. Publishing changes the product status from `draft` to `published`.
3. Publishing updates the product `updated_at` timestamp.
4. A seller cannot publish another seller's product.
5. Publishing an already published product returns a clear safe error and does not create duplicate side effects.
6. Public product list returns only published products.
7. Public product detail returns a published product by ID.
8. Public product detail does not return draft products.
9. Public product responses do not expose password hashes, JWT data or private seller auth fields.
10. Seller product list/detail behavior continues to work.
11. Frontend seller product UI allows publishing a draft product.
12. Frontend public catalog can list published products.
13. Frontend public product detail can show a published product.
14. Draft products are not shown in public frontend catalog.
15. Backend tests cover publish success, ownership rejection, repeated publish and public visibility.
16. Frontend tests cover the most important new UI behavior.
17. `make validate` passes.
18. No unrelated marketplace features are added.
19. No AI orchestration, repository skills or autonomous Git behavior is added.

## API impact

Expected new or changed endpoints:

### Seller endpoint

```text
POST /api/v1/seller/products/{productID}/publish
```

Expected behavior:

- Requires seller authentication.
- Publishes only the authenticated seller's own draft product.
- Returns the published product or a concise success response.
- Returns a safe error for unauthorized ownership or invalid state.

### Public catalog endpoints

```text
GET /api/v1/products
GET /api/v1/products/{productID}
```

Expected behavior:

- Does not require authentication.
- Returns only products with status `published`.
- Does not expose private seller authentication data.

Endpoint names may be adjusted only if the existing route conventions clearly require it.

## Database impact

No new table is expected.

A migration is only allowed if the current schema cannot safely support explicit product status values or public catalog queries.

If a migration is required:

- Explain why.
- Add both up and down migrations where safe.
- Preserve existing data.
- Ensure SQL failures return non-zero.
- Add tests or verification for migration behavior.

## Security considerations

- Backend authorization is mandatory for publish.
- Frontend checks are not security controls.
- Public catalog responses must be sanitized.
- Do not expose internal SQL errors.
- Do not expose private seller fields.
- Do not log access tokens, passwords, password hashes or JWT secrets.
- Do not trust client-provided product status.

## Technical constraints

- Keep the backend as a modular monolith.
- Keep HTTP handlers thin.
- Business rules belong in product service logic.
- SQL belongs in repositories.
- Avoid speculative abstractions.
- Do not add new production dependencies unless the approved plan proves they are necessary.
- Use existing error-response conventions.
- Use existing frontend routing and API-client conventions.
- Keep frontend styling minimal.

## Required tests

### Backend tests

Cover:

1. Owner can publish own draft product.
2. Non-owner cannot publish another seller's product.
3. Published product appears in public list.
4. Draft product does not appear in public list.
5. Public detail returns a published product.
6. Public detail does not return a draft product.
7. Already published product cannot be republished.
8. Public response does not include private seller auth fields.

### Frontend tests

Cover the most important feasible cases:

1. Seller can trigger publish action for a draft product.
2. Public catalog renders published products.
3. Public detail renders a published product.
4. Draft products are not rendered in the public catalog when API returns only published products, or the frontend correctly relies on the public API contract.

Do not create broad brittle UI tests.

## Required evidence

Create evidence under:

```text
evidence/PROD-001/
```

Recommended files:

- `plan.md`
- `acceptance-criteria.md`
- `api-examples.md`
- `test-summary.md`
- `validation-summary.md`
- `review-report.md`
- `completion-summary.md`

Evidence must include:

- Exact commands executed.
- Result of each command.
- Acceptance-criteria mapping.
- API request and response examples for publish and public catalog behavior.
- Database state used for integration validation.
- Remaining risks.
- Confirmation that no secrets, passwords, password hashes or reusable tokens are stored.
- Confirmation that no out-of-scope marketplace features were added.

## Assumptions

- The current product table already has a status field.
- Public catalog can start with a simple list without advanced pagination or search.
- Buyer authentication is not required for this task.
- Product images are not required.
- Seller display information may be included only if the current schema supports it safely.

## Open questions

None currently blocking.

## Definition of done

The task is complete only when:

- Publish behavior is implemented and tested.
- Public catalog behavior is implemented and tested.
- Draft products remain private.
- Seller ownership is enforced.
- Public responses are sanitized.
- Frontend supports the new basic flow.
- Evidence maps acceptance criteria to commands and results.
- `make validate` passes.
- Independent review has no Blocker or High findings.
- Human approval is still required for commit and merge.

## Completion record

To be filled after implementation and review.
