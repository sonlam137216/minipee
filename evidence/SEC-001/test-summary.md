# SEC-001 Test Summary

## Added Coverage

- `TestJWTManagerVerifyAcceptsCorrectlySignedNonExpiredToken`
  - Issues a correctly signed, non-expired seller token.
  - Verifies the token returns seller ID and email claims.
- `TestJWTManagerVerifyUsesInjectedClockForExpiration`
  - Issues a token at a fixed UTC time.
  - Verifies the same token is accepted before its fixed expiration.
  - Verifies the same token is rejected after its fixed expiration by moving only the injected verifier time.
- `TestJWTManagerVerifyRejectsInvalidTokens`
  - Rejects a token signed with a different secret.
  - Rejects a malformed token string.
  - Rejects a token whose expiration is before the injected verifier time.
  - Asserts the public service error is `ErrInvalidAccessToken`.
- `TestRequireSellerAcceptsValidBearerToken`
  - Verifies the middleware accepts a valid bearer token.
  - Verifies the seller ID reaches request context.
- `TestRequireSellerRejectsMissingAndInvalidBearerTokensWithGenericResponse`
  - Rejects missing, wrong-secret, malformed, and expired bearer tokens.
  - Verifies the protected handler is not called on rejection.
  - Verifies each rejection returns HTTP `401` with `unauthenticated` and `Authentication required`.
- `TestAssertExactAuthRejectionBodyRejectsUnexpectedDetails`
  - Verifies the strict SEC-001 response assertion rejects an unexpected JWT detail field.
- `TestAssertExactAuthRejectionBodyRejectsTrailingJSON`
  - Verifies the strict SEC-001 response assertion rejects trailing JSON content.

## Review Findings Addressed

- Finding 1: Rejection response shape was not asserted strictly.
  - Correction: middleware rejection tests now use `assertExactAuthRejectionBody`, which uses `json.Decoder.DisallowUnknownFields`, validates public code and message, and checks for EOF after the first JSON value.
- Finding 2: JWT expiration verification used the process wall clock instead of `JWTManager.now`.
  - Correction: `JWTManager.Verify` now passes the existing injected clock to the JWT parser with `jwt.WithTimeFunc(m.now)`.

## Determinism

- No test uses `time.Sleep`.
- Valid and expired token outcomes are relative to explicit fixed UTC verifier times.
- Moving only the injected verifier time changes token validity as expected.
- Fixed test-only secrets are used in tests only.
- No complete access tokens are recorded in evidence.
