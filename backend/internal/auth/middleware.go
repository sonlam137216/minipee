package auth

import (
	"net/http"
	"strings"

	"marketplace/backend/internal/httpapi"
)

func RequireSeller(tokens *JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			rawToken := strings.TrimPrefix(header, "Bearer ")
			if rawToken == "" || rawToken == header {
				httpapi.WriteError(w, http.StatusUnauthorized, "unauthenticated", "Authentication required", nil)
				return
			}
			claims, err := tokens.Verify(rawToken)
			if err != nil {
				httpapi.WriteError(w, http.StatusUnauthorized, "unauthenticated", "Authentication required", nil)
				return
			}
			next.ServeHTTP(w, r.WithContext(ContextWithSellerID(r.Context(), claims.SellerID)))
		})
	}
}
