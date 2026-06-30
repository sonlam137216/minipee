package products

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

type memoryProductRepository struct {
	products map[string]Product
}

func (r *memoryProductRepository) Create(ctx context.Context, product Product) (Product, error) {
	if r.products == nil {
		r.products = map[string]Product{}
	}
	r.products[product.ID] = product
	return product, nil
}

func (r *memoryProductRepository) ListBySeller(ctx context.Context, sellerID string) ([]Product, error) {
	var products []Product
	for _, product := range r.products {
		if product.SellerID == sellerID {
			products = append(products, product)
		}
	}
	return products, nil
}

func (r *memoryProductRepository) FindSellerDraftByID(ctx context.Context, sellerID string, productID string) (Product, error) {
	product, ok := r.products[productID]
	if !ok || product.SellerID != sellerID || product.Status != StatusDraft {
		return Product{}, ErrProductNotFound
	}
	return product, nil
}

func TestProductNameShorterThanThreeCharactersIsRejected(t *testing.T) {
	service := NewService(&memoryProductRepository{}, fixedID("00000000-0000-4000-8000-000000000001"), fixedTime)

	_, err := service.CreateDraft(context.Background(), "seller-1", CreateProductInput{Name: "ab"})

	if !errors.Is(err, ErrInvalidProductName) {
		t.Fatalf("error = %v, want ErrInvalidProductName", err)
	}
}

func TestProductNameLongerThanTwoHundredCharactersIsRejected(t *testing.T) {
	service := NewService(&memoryProductRepository{}, fixedID("00000000-0000-4000-8000-000000000001"), fixedTime)

	_, err := service.CreateDraft(context.Background(), "seller-1", CreateProductInput{Name: strings.Repeat("a", 201)})

	if !errors.Is(err, ErrInvalidProductName) {
		t.Fatalf("error = %v, want ErrInvalidProductName", err)
	}
}

func TestNewProductStatusIsAlwaysDraft(t *testing.T) {
	service := NewService(&memoryProductRepository{}, fixedID("00000000-0000-4000-8000-000000000001"), fixedTime)

	product, err := service.CreateDraft(context.Background(), "seller-1", CreateProductInput{Name: "Draft product"})

	if err != nil {
		t.Fatalf("CreateDraft returned error: %v", err)
	}
	if product.Status != StatusDraft {
		t.Fatalf("Status = %q, want %q", product.Status, StatusDraft)
	}
}

func TestClientProvidedStatusCannotOverrideDraft(t *testing.T) {
	service := NewService(&memoryProductRepository{}, fixedID("00000000-0000-4000-8000-000000000001"), fixedTime)

	product, err := service.CreateDraft(context.Background(), "seller-1", CreateProductInput{
		Name:   "Draft product",
		Status: "published",
	})

	if err != nil {
		t.Fatalf("CreateDraft returned error: %v", err)
	}
	if product.Status != StatusDraft {
		t.Fatalf("Status = %q, want %q", product.Status, StatusDraft)
	}
}

func TestProductOwnershipIsEnforced(t *testing.T) {
	repository := &memoryProductRepository{products: map[string]Product{
		"product-1": {ID: "product-1", SellerID: "seller-1", Name: "Draft product", Status: StatusDraft},
	}}
	service := NewService(repository, fixedID("unused"), fixedTime)

	_, err := service.GetSellerDraft(context.Background(), "seller-2", "product-1")

	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("error = %v, want ErrProductNotFound", err)
	}
}

func fixedID(id string) func() string {
	return func() string {
		return id
	}
}

func fixedTime() time.Time {
	return time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)
}
