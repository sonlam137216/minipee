package products

import (
	"context"
	"os"
	"testing"

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

	sellerID := "00000000-0000-4000-8000-000000000201"
	otherSellerID := "00000000-0000-4000-8000-000000000202"
	cleanup := func() {
		_, _ = db.Exec(ctx, "DELETE FROM products WHERE seller_id IN ($1, $2)", sellerID, otherSellerID)
		_, _ = db.Exec(ctx, "DELETE FROM sellers WHERE id IN ($1, $2)", sellerID, otherSellerID)
	}
	cleanup()
	t.Cleanup(cleanup)
	_, err = db.Exec(ctx, `
		INSERT INTO sellers (id, email, password_hash, display_name, created_at, updated_at)
		VALUES ($1, 'product-owner@example.com', 'hash', 'Owner', now(), now()),
		       ($2, 'other-owner@example.com', 'hash', 'Other', now(), now())
	`, sellerID, otherSellerID)
	if err != nil {
		t.Fatalf("insert sellers: %v", err)
	}

	service := NewService(NewPostgresRepository(db), fixedID("00000000-0000-4000-8000-000000000301"), fixedTime)
	product, err := service.CreateDraft(ctx, sellerID, CreateProductInput{Name: "Draft product", Description: "Description"})
	if err != nil {
		t.Fatalf("CreateDraft returned error: %v", err)
	}

	products, err := service.ListSellerProducts(ctx, sellerID)
	if err != nil {
		t.Fatalf("ListSellerProducts returned error: %v", err)
	}
	if len(products) != 1 {
		t.Fatalf("listed products = %d, want 1", len(products))
	}

	found, err := service.GetSellerDraft(ctx, sellerID, product.ID)
	if err != nil {
		t.Fatalf("GetSellerDraft returned error: %v", err)
	}
	if found.ID != product.ID {
		t.Fatalf("found ID = %q, want %q", found.ID, product.ID)
	}

	_, err = service.GetSellerDraft(ctx, otherSellerID, product.ID)
	if err != ErrProductNotFound {
		t.Fatalf("other seller error = %v, want ErrProductNotFound", err)
	}
}
