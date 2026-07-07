package products

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresProductRepository(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	runID := uint64(time.Now().UnixNano())
	sellerID := testUUID(runID)
	otherSellerID := testUUID(runID + 1)
	cleanup := func() {
		_, _ = db.Exec(ctx, "DELETE FROM products WHERE seller_id IN ($1, $2)", sellerID, otherSellerID)
		_, _ = db.Exec(ctx, "DELETE FROM sellers WHERE id IN ($1, $2)", sellerID, otherSellerID)
	}
	cleanup()
	t.Cleanup(cleanup)
	_, err = db.Exec(ctx, `
		INSERT INTO sellers (id, email, password_hash, display_name, created_at, updated_at)
		VALUES ($1, $2, 'hash', 'Owner', now(), now()),
		       ($3, $4, 'hash', 'Other', now(), now())
	`, sellerID, fmt.Sprintf("product-owner-%d@example.com", runID), otherSellerID, fmt.Sprintf("other-owner-%d@example.com", runID))
	if err != nil {
		t.Fatalf("insert sellers: %v", err)
	}

	now := fixedTime()
	nextID := sequenceID(
		testUUID(runID+2),
		testUUID(runID+3),
	)
	service := NewService(NewPostgresRepository(db), nextID, func() time.Time {
		return now
	})
	product, err := service.CreateDraft(ctx, sellerID, CreateProductInput{Name: "Draft product", Description: "Description"})
	if err != nil {
		t.Fatalf("CreateDraft returned error: %v", err)
	}
	now = now.Add(5 * time.Minute)
	published, err := service.CreateDraft(ctx, sellerID, CreateProductInput{Name: "Published product", Description: "Public description"})
	if err != nil {
		t.Fatalf("CreateDraft returned error for published fixture: %v", err)
	}

	products, err := service.ListSellerProducts(ctx, sellerID)
	if err != nil {
		t.Fatalf("ListSellerProducts returned error: %v", err)
	}
	if len(products) != 2 {
		t.Fatalf("listed products = %d, want 2", len(products))
	}

	found, err := service.GetSellerProduct(ctx, sellerID, product.ID)
	if err != nil {
		t.Fatalf("GetSellerProduct returned error: %v", err)
	}
	if found.ID != product.ID {
		t.Fatalf("found ID = %q, want %q", found.ID, product.ID)
	}

	_, err = service.GetSellerProduct(ctx, otherSellerID, product.ID)
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("other seller error = %v, want ErrProductNotFound", err)
	}

	publishAt := now.Add(5 * time.Minute)
	now = publishAt
	published, err = service.PublishSellerProduct(ctx, sellerID, published.ID)
	if err != nil {
		t.Fatalf("PublishSellerProduct returned error: %v", err)
	}
	if published.Status != StatusPublished {
		t.Fatalf("published status = %q, want %q", published.Status, StatusPublished)
	}
	if !published.UpdatedAt.Equal(publishAt) {
		t.Fatalf("published UpdatedAt = %v, want %v", published.UpdatedAt, publishAt)
	}

	foundPublished, err := service.GetSellerProduct(ctx, sellerID, published.ID)
	if err != nil {
		t.Fatalf("GetSellerProduct for published product returned error: %v", err)
	}
	if foundPublished.Status != StatusPublished {
		t.Fatalf("seller detail status = %q, want %q", foundPublished.Status, StatusPublished)
	}

	publicProducts, err := service.ListPublishedProducts(ctx)
	if err != nil {
		t.Fatalf("ListPublishedProducts returned error: %v", err)
	}
	foundPublishedFixture := false
	for _, publicProduct := range publicProducts {
		if publicProduct.Status != StatusPublished {
			t.Fatalf("public product %q status = %q, want %q", publicProduct.ID, publicProduct.Status, StatusPublished)
		}
		if publicProduct.ID == product.ID {
			t.Fatalf("publicProducts includes draft fixture %#v", publicProduct)
		}
		if publicProduct.ID == published.ID {
			foundPublishedFixture = true
		}
	}
	if !foundPublishedFixture {
		t.Fatalf("publicProducts = %#v, want published fixture %q included", publicProducts, published.ID)
	}

	publicDetail, err := service.GetPublishedProduct(ctx, published.ID)
	if err != nil {
		t.Fatalf("GetPublishedProduct returned error: %v", err)
	}
	if publicDetail.ID != published.ID {
		t.Fatalf("public detail ID = %q, want %q", publicDetail.ID, published.ID)
	}
	_, err = service.GetPublishedProduct(ctx, product.ID)
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("draft public detail error = %v, want ErrProductNotFound", err)
	}
}

func testUUID(value uint64) string {
	return fmt.Sprintf("00000000-0000-4000-8000-%012x", value&0xffffffffffff)
}
