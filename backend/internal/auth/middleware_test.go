package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRequireSellerAcceptsValidBearerToken(t *testing.T) {
	manager := NewJWTManager("sec-001-middleware-secret", time.Hour, fixedJWTValidVerificationTime)
	issuer := NewJWTManager("sec-001-middleware-secret", time.Hour, fixedJWTIssueTime)
	token, err := issuer.Issue(Seller{ID: "seller-1", Email: "seller@example.com"})
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}
	protected := RequireSeller(manager)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sellerID, ok := SellerIDFromContext(r.Context())
		if !ok {
			t.Fatal("seller ID missing from request context")
		}
		if sellerID != "seller-1" {
			t.Fatalf("sellerID = %q, want seller-1", sellerID)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()

	protected.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
}

func TestRequireSellerRejectsMissingAndInvalidBearerTokensWithGenericResponse(t *testing.T) {
	verifier := NewJWTManager("sec-001-middleware-secret", time.Hour, fixedJWTExpiredVerificationTime)
	otherIssuer := NewJWTManager("sec-001-other-secret", time.Hour, fixedJWTIssueTime)
	tokenSignedWithDifferentSecret, err := otherIssuer.Issue(Seller{ID: "seller-1", Email: "seller@example.com"})
	if err != nil {
		t.Fatalf("Issue alternate token returned error: %v", err)
	}
	expiredIssuer := NewJWTManager("sec-001-middleware-secret", time.Hour, fixedJWTIssueTime)
	expiredToken, err := expiredIssuer.Issue(Seller{ID: "seller-1", Email: "seller@example.com"})
	if err != nil {
		t.Fatalf("Issue expired token returned error: %v", err)
	}
	protected := RequireSeller(verifier)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("protected handler was called for rejected authentication")
	}))

	tests := []struct {
		name        string
		authHeader  string
		wantMessage string
	}{
		{name: "missing token", wantMessage: "Authentication required"},
		{name: "token signed using different secret", authHeader: "Bearer " + tokenSignedWithDifferentSecret, wantMessage: "Authentication required"},
		{name: "malformed token", authHeader: "Bearer not-a-jwt", wantMessage: "Authentication required"},
		{name: "expired token", authHeader: "Bearer " + expiredToken, wantMessage: "Authentication required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				request.Header.Set("Authorization", tt.authHeader)
			}
			response := httptest.NewRecorder()

			protected.ServeHTTP(response, request)

			if response.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
			}
			if err := assertExactAuthRejectionBody(response.Body.Bytes(), tt.wantMessage); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestAssertExactAuthRejectionBodyRejectsUnexpectedDetails(t *testing.T) {
	body := []byte(`{"error":{"code":"unauthenticated","message":"Authentication required","details":"token is expired"}}`)

	if err := assertExactAuthRejectionBody(body, "Authentication required"); err == nil {
		t.Fatal("assertExactAuthRejectionBody accepted unexpected detail field")
	}
}

func TestAssertExactAuthRejectionBodyRejectsTrailingJSON(t *testing.T) {
	body := []byte(`{"error":{"code":"unauthenticated","message":"Authentication required"}} {"error":{"code":"unauthenticated","message":"Authentication required"}}`)

	if err := assertExactAuthRejectionBody(body, "Authentication required"); err == nil {
		t.Fatal("assertExactAuthRejectionBody accepted trailing JSON content")
	}
}

func assertExactAuthRejectionBody(rawBody []byte, wantMessage string) error {
	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	decoder := json.NewDecoder(bytes.NewReader(rawBody))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&body); err != nil {
		return fmt.Errorf("decode exact authentication rejection body: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return fmt.Errorf("authentication rejection body has trailing JSON content")
	}
	if body.Error.Code != "unauthenticated" {
		return fmt.Errorf("error code = %q, want unauthenticated", body.Error.Code)
	}
	if body.Error.Message != wantMessage {
		return fmt.Errorf("error message = %q, want %q", body.Error.Message, wantMessage)
	}
	return nil
}
