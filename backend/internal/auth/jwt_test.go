package auth

import (
	"errors"
	"testing"
	"time"
)

func TestJWTManagerVerifyAcceptsCorrectlySignedNonExpiredToken(t *testing.T) {
	issuer := NewJWTManager("sec-001-valid-secret", time.Hour, fixedJWTIssueTime)
	token, err := issuer.Issue(Seller{ID: "seller-1", Email: "seller@example.com"})
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}
	verifier := NewJWTManager("sec-001-valid-secret", time.Hour, fixedJWTValidVerificationTime)

	claims, err := verifier.Verify(token)

	if err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}
	if claims.SellerID != "seller-1" {
		t.Fatalf("SellerID = %q, want seller-1", claims.SellerID)
	}
	if claims.Email != "seller@example.com" {
		t.Fatalf("Email = %q, want seller@example.com", claims.Email)
	}
}

func TestJWTManagerVerifyUsesInjectedClockForExpiration(t *testing.T) {
	issuer := NewJWTManager("sec-001-clock-secret", time.Hour, fixedJWTIssueTime)
	token, err := issuer.Issue(Seller{ID: "seller-1", Email: "seller@example.com"})
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	validVerifier := NewJWTManager("sec-001-clock-secret", time.Hour, fixedJWTValidVerificationTime)
	if _, err := validVerifier.Verify(token); err != nil {
		t.Fatalf("Verify before fixed expiration returned error: %v", err)
	}

	expiredVerifier := NewJWTManager("sec-001-clock-secret", time.Hour, fixedJWTExpiredVerificationTime)
	if _, err := expiredVerifier.Verify(token); !errors.Is(err, ErrInvalidAccessToken) {
		t.Fatalf("Verify after fixed expiration error = %v, want ErrInvalidAccessToken", err)
	}
}

func TestJWTManagerVerifyRejectsInvalidTokens(t *testing.T) {
	issuer := NewJWTManager("sec-001-issuer-secret", time.Hour, fixedJWTIssueTime)
	tokenSignedWithDifferentSecret, err := issuer.Issue(Seller{ID: "seller-1", Email: "seller@example.com"})
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}
	expiredIssuer := NewJWTManager("sec-001-valid-secret", time.Hour, fixedJWTIssueTime)
	expiredToken, err := expiredIssuer.Issue(Seller{ID: "seller-1", Email: "seller@example.com"})
	if err != nil {
		t.Fatalf("Issue expired token returned error: %v", err)
	}

	manager := NewJWTManager("sec-001-valid-secret", time.Hour, fixedJWTExpiredVerificationTime)
	tests := []struct {
		name     string
		rawToken string
	}{
		{name: "token signed using different secret", rawToken: tokenSignedWithDifferentSecret},
		{name: "malformed token", rawToken: "not-a-jwt"},
		{name: "expired token", rawToken: expiredToken},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := manager.Verify(tt.rawToken)

			if !errors.Is(err, ErrInvalidAccessToken) {
				t.Fatalf("error = %v, want ErrInvalidAccessToken", err)
			}
		})
	}
}

func fixedJWTIssueTime() time.Time {
	return time.Date(2035, 1, 1, 12, 0, 0, 0, time.UTC)
}

func fixedJWTValidVerificationTime() time.Time {
	return time.Date(2035, 1, 1, 12, 30, 0, 0, time.UTC)
}

func fixedJWTExpiredVerificationTime() time.Time {
	return time.Date(2035, 1, 1, 14, 0, 0, 0, time.UTC)
}
