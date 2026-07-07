# PROD-001 Publish Product And Public Catalog Plan

## Accepted Scope

Implement the smallest publish and public catalog path:

- Add product status `published` alongside `draft`.
- Add protected seller publish endpoint `POST /api/v1/seller/products/{productID}/publish`.
- Add unauthenticated public catalog endpoints:
  - `GET /api/v1/products`
  - `GET /api/v1/products/{productID}`
- Keep seller-created products forced to draft.
- Keep seller-owned reads and writes scoped by authenticated seller ID.
- Ensure public responses exclude `sellerId` and auth-owned fields.
- Add frontend `/catalog` and `/catalog/:productID` pages.
- Keep existing protected seller product routes.
- Add focused backend, repository, handler, and frontend tests.

## Explicitly Excluded

Product images, variants, categories, inventory, search, pagination, cart, orders, payments, shipping, buyer accounts, moderation, unpublish, product editing, microservices, events, AI orchestration, deployment, autonomous push, and autonomous merge.
