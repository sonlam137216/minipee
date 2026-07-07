# PROD-001 Acceptance Criteria

| AC | Result | Evidence |
| --- | --- | --- |
| 1 | Pass | Seller A publish smoke returned `200`; backend service and handler tests cover owner publish. |
| 2 | Pass | Service, repository integration, and smoke output show status changes from `draft` to `published`. |
| 3 | Pass | Service and repository integration tests assert publish `updated_at` uses publish clock. |
| 4 | Pass | Smoke returned `404 Product not found` for seller B publishing seller A product; service test covers non-owner. |
| 5 | Pass | Smoke returned `409 already_published`; handler and service tests cover repeat publish. |
| 6 | Pass | Service, repository integration, frontend test, and smoke verify public list excludes drafts. |
| 7 | Pass | Repository integration and smoke verify public detail for published product returns `200`. |
| 8 | Pass | Repository integration verifies public detail rejects draft product with `ErrProductNotFound`. |
| 9 | Pass | Handler test decodes public JSON and checks forbidden fields; smoke public responses omit `sellerId`. |
| 10 | Pass | Existing seller tests pass; repository integration verifies seller can view owned published product. |
| 11 | Pass | Frontend test clicks publish button and verifies `POST` with bearer token. |
| 12 | Pass | Frontend test renders `/catalog` from `GET /api/v1/products`. |
| 13 | Pass | Frontend test renders `/catalog/:productID` from public detail API. |
| 14 | Pass | Frontend catalog test asserts draft product name is absent when public API excludes drafts. |
| 15 | Pass | Backend service, handler, and repository tests added for publish/public visibility. |
| 16 | Pass | Frontend App tests added for publish/catalog/detail. |
| 17 | Pass | `make validate` exited `0`. |
| 18 | Pass | Diff review found no images, variants, categories, inventory, cart, orders, payments, or shipping. |
| 19 | Pass | Diff review found no AI orchestration, autonomous push, merge, or deploy behavior. |
