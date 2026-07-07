package products

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"marketplace/backend/internal/auth"
)

func TestPublishMapsAlreadyPublishedToConflict(t *testing.T) {
	repository := &memoryProductRepository{products: map[string]Product{
		"product-1": {
			ID:        "product-1",
			SellerID:  "seller-1",
			Name:      "Published product",
			Status:    StatusPublished,
			UpdatedAt: fixedTime(),
		},
	}}
	handler := NewHandler(NewService(repository, fixedID("unused"), func() time.Time {
		return fixedTime().Add(time.Hour)
	}))
	request := httptest.NewRequest(http.MethodPost, "/seller/products/product-1/publish", nil)
	request = request.WithContext(auth.ContextWithSellerID(withProductID(request.Context(), "product-1"), "seller-1"))
	recorder := httptest.NewRecorder()

	handler.Publish(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusConflict)
	}
	if !strings.Contains(recorder.Body.String(), `"code":"already_published"`) {
		t.Fatalf("body = %s, want already_published code", recorder.Body.String())
	}
}

func TestPublicProductResponseExcludesSellerAndAuthFields(t *testing.T) {
	repository := &memoryProductRepository{products: map[string]Product{
		"product-1": {
			ID:          "product-1",
			SellerID:    "seller-1",
			Name:        "Published product",
			Description: "Visible description",
			Status:      StatusPublished,
			CreatedAt:   fixedTime(),
			UpdatedAt:   fixedTime(),
		},
	}}
	handler := NewHandler(NewService(repository, fixedID("unused"), fixedTime))
	request := httptest.NewRequest(http.MethodGet, "/products/product-1", nil)
	request = request.WithContext(withProductID(request.Context(), "product-1"))
	recorder := httptest.NewRecorder()

	handler.PublicGet(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	for _, forbidden := range []string{"sellerId", "passwordHash", "accessToken", "jwt", "email"} {
		if _, ok := body[forbidden]; ok {
			t.Fatalf("public response includes forbidden field %q in %#v", forbidden, body)
		}
	}
	if body["name"] != "Published product" {
		t.Fatalf("name = %#v, want Published product", body["name"])
	}
}

func withProductID(ctx context.Context, productID string) context.Context {
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("productID", productID)
	return context.WithValue(ctx, chi.RouteCtxKey, routeContext)
}
