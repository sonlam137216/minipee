# Security Boundaries

This document separates implemented security controls from current MVP limitations.

## Implemented Controls

### Password Handling

- Registration accepts a plaintext password only as request input.
- `auth.Service.Register` hashes passwords with `bcrypt.GenerateFromPassword` and `bcrypt.DefaultCost`.
- `auth.Service.Login` verifies passwords with `bcrypt.CompareHashAndPassword`.
- `sellers.password_hash` stores the bcrypt hash.
- Auth API responses use `sellerResponse`, which excludes `password_hash`.
- Plaintext passwords and password hashes must not be logged or included in evidence.

### JWT Creation And Validation

- `auth.JWTManager` signs access tokens with HS256.
- Tokens include seller ID, email, subject, issued-at time, and expiration time.
- JWT expiration is required during verification.
- Tokens are signed with `JWT_SECRET` from configuration.
- Invalid tokens return the auth invalid-token path and are translated to authentication failure by middleware.

### Authentication Middleware

- Seller product routes are inside a route group protected by `auth.RequireSeller`.
- The middleware requires an `Authorization: Bearer <token>` header.
- The middleware verifies the JWT before calling the next handler.
- Missing, malformed, expired, invalid, or wrongly signed tokens return `401` with a generic authentication message.

### Seller Identity Propagation

- After token verification, `auth.RequireSeller` stores seller ID in request context.
- Product handlers read seller ID with `auth.SellerIDFromContext`.
- If seller ID is missing from context, product handlers return `401`.

### Product Ownership Enforcement

- Product creation uses the seller ID from request context.
- Product list queries filter by `seller_id`.
- Seller product detail queries require matching product ID and matching seller ID.
- Product publish queries require matching product ID, matching seller ID, and draft status.
- Cross-seller product detail access returns not found.
- Cross-seller publish attempts return not found and leave the product unchanged.
- Frontend state is not trusted for backend authorization.

### Public Product Catalog Boundary

- Public catalog routes are mounted outside `auth.RequireSeller`.
- Public catalog list and detail queries filter by `status = 'published'`.
- Draft products are not returned by public list or detail routes.
- Public product responses exclude `sellerId`, seller email, password hashes, access tokens, JWT claims, and auth-owned fields.

### Secret Configuration

- `backend/internal/config` requires `DATABASE_URL`, `JWT_SECRET`, `JWT_EXPIRATION_MINUTES`, and `FRONTEND_ORIGIN`.
- `JWT_EXPIRATION_MINUTES` must parse as a positive integer.
- `.env` is ignored by Git.
- Secrets and access tokens must not be committed in task evidence or logs.

### Error Sanitization

- HTTP handlers map domain errors to controlled JSON error responses.
- Authentication failures use generic invalid-credential or authentication-required messages.
- Unknown backend errors return `internal_error` with `Internal server error`.
- JSON decoding rejects unknown fields through `httpapi.DecodeJSON`.

### Logging Restrictions

- The current request logger records method, path, and duration.
- It does not intentionally log request bodies, passwords, password hashes, JWT secrets, or access tokens.
- Future logging changes must preserve that restriction.

### Frontend Token Storage

- The frontend stores the JWT access token and seller snapshot in `localStorage`.
- The frontend uses that token to send `Authorization: Bearer <token>` requests.
- This is browser UI state only. Backend authorization remains the source of truth.

## Known Limitations And MVP Trade-Offs

- Access tokens are stored in `localStorage`, which is acceptable only for the local MVP trade-off documented in this repository.
- Refresh tokens are not implemented.
- Token revocation is not implemented.
- Dedicated invalid and expired JWT tests are still listed as known gaps.
- More advanced browser storage hardening is deferred.
- There is no per-test database isolation beyond cleanup in current integration tests.
- Migration version tracking is not implemented.
- Product unpublish, moderation, and public seller profiles are not implemented.
- Production-grade logging, rate limiting, CSRF strategy, and deployment secret management are not implemented.
