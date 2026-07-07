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

func (r *memoryProductRepository) FindBySellerID(ctx context.Context, sellerID string, productID string) (Product, error) {
	product, ok := r.products[productID]
	if !ok || product.SellerID != sellerID {
		return Product{}, ErrProductNotFound
	}
	return product, nil
}

func (r *memoryProductRepository) PublishDraftBySeller(ctx context.Context, sellerID string, productID string, now time.Time) (Product, error) {
	product, ok := r.products[productID]
	if !ok || product.SellerID != sellerID || product.Status != StatusDraft {
		return Product{}, ErrProductNotFound
	}
	product.Status = StatusPublished
	product.UpdatedAt = now.UTC()
	r.products[productID] = product
	return product, nil
}

func (r *memoryProductRepository) ListPublished(ctx context.Context) ([]Product, error) {
	var products []Product
	for _, product := range r.products {
		if product.Status == StatusPublished {
			products = append(products, product)
		}
	}
	return products, nil
}

func (r *memoryProductRepository) FindPublishedByID(ctx context.Context, productID string) (Product, error) {
	product, ok := r.products[productID]
	if !ok || product.Status != StatusPublished {
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

	_, err := service.GetSellerProduct(context.Background(), "seller-2", "product-1")

	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("error = %v, want ErrProductNotFound", err)
	}
}

func TestOwnerPublishesDraftProduct(t *testing.T) {
	createdAt := fixedTime()
	publishAt := createdAt.Add(5 * time.Minute)
	repository := &memoryProductRepository{products: map[string]Product{
		"product-1": {
			ID:        "product-1",
			SellerID:  "seller-1",
			Name:      "Draft product",
			Status:    StatusDraft,
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
		},
	}}
	service := NewService(repository, fixedID("unused"), func() time.Time { return publishAt })

	product, err := service.PublishSellerProduct(context.Background(), "seller-1", "product-1")

	if err != nil {
		t.Fatalf("PublishSellerProduct returned error: %v", err)
	}
	if product.Status != StatusPublished {
		t.Fatalf("Status = %q, want %q", product.Status, StatusPublished)
	}
	if !product.UpdatedAt.Equal(publishAt) {
		t.Fatalf("UpdatedAt = %v, want %v", product.UpdatedAt, publishAt)
	}
}

func TestNonOwnerCannotPublishProduct(t *testing.T) {
	repository := &memoryProductRepository{products: map[string]Product{
		"product-1": {ID: "product-1", SellerID: "seller-1", Name: "Draft product", Status: StatusDraft},
	}}
	service := NewService(repository, fixedID("unused"), fixedTime)

	_, err := service.PublishSellerProduct(context.Background(), "seller-2", "product-1")

	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("error = %v, want ErrProductNotFound", err)
	}
	if repository.products["product-1"].Status != StatusDraft {
		t.Fatalf("Status = %q, want %q", repository.products["product-1"].Status, StatusDraft)
	}
}

func TestAlreadyPublishedProductCannotBePublishedAgain(t *testing.T) {
	updatedAt := fixedTime()
	repository := &memoryProductRepository{products: map[string]Product{
		"product-1": {
			ID:        "product-1",
			SellerID:  "seller-1",
			Name:      "Published product",
			Status:    StatusPublished,
			UpdatedAt: updatedAt,
		},
	}}
	service := NewService(repository, fixedID("unused"), func() time.Time { return updatedAt.Add(time.Hour) })

	_, err := service.PublishSellerProduct(context.Background(), "seller-1", "product-1")

	if !errors.Is(err, ErrProductAlreadyPublished) {
		t.Fatalf("error = %v, want ErrProductAlreadyPublished", err)
	}
	if !repository.products["product-1"].UpdatedAt.Equal(updatedAt) {
		t.Fatalf("UpdatedAt changed to %v, want %v", repository.products["product-1"].UpdatedAt, updatedAt)
	}
}

func TestPublicProductsOnlyIncludePublishedProducts(t *testing.T) {
	repository := &memoryProductRepository{products: map[string]Product{
		"draft":     {ID: "draft", SellerID: "seller-1", Name: "Draft product", Status: StatusDraft},
		"published": {ID: "published", SellerID: "seller-1", Name: "Published product", Status: StatusPublished},
	}}
	service := NewService(repository, fixedID("unused"), fixedTime)

	products, err := service.ListPublishedProducts(context.Background())

	if err != nil {
		t.Fatalf("ListPublishedProducts returned error: %v", err)
	}
	if len(products) != 1 || products[0].ID != "published" {
		t.Fatalf("products = %#v, want only published product", products)
	}
}

func TestPublishedProductDetailRejectsDrafts(t *testing.T) {
	repository := &memoryProductRepository{products: map[string]Product{
		"draft":     {ID: "draft", SellerID: "seller-1", Name: "Draft product", Status: StatusDraft},
		"published": {ID: "published", SellerID: "seller-1", Name: "Published product", Status: StatusPublished},
	}}
	service := NewService(repository, fixedID("unused"), fixedTime)

	product, err := service.GetPublishedProduct(context.Background(), "published")
	if err != nil {
		t.Fatalf("GetPublishedProduct returned error: %v", err)
	}
	if product.ID != "published" {
		t.Fatalf("ID = %q, want published", product.ID)
	}
	_, err = service.GetPublishedProduct(context.Background(), "draft")
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("draft error = %v, want ErrProductNotFound", err)
	}
}

func fixedID(id string) func() string {
	return func() string {
		return id
	}
}

func sequenceID(ids ...string) func() string {
	index := 0
	return func() string {
		id := ids[index]
		index++
		return id
	}
}

func fixedTime() time.Time {
	return time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)
}
