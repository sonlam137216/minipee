package products

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, product Product) (Product, error) {
	_, err := r.db.Exec(ctx, `
		INSERT INTO products (id, seller_id, name, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, product.ID, product.SellerID, product.Name, product.Description, product.Status, product.CreatedAt, product.UpdatedAt)
	if err != nil {
		return Product{}, err
	}
	return product, nil
}

func (r *PostgresRepository) ListBySeller(ctx context.Context, sellerID string) ([]Product, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, seller_id, name, description, status, created_at, updated_at
		FROM products
		WHERE seller_id = $1
		ORDER BY created_at DESC, id DESC
	`, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(
			&product.ID,
			&product.SellerID,
			&product.Name,
			&product.Description,
			&product.Status,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *PostgresRepository) FindSellerDraftByID(ctx context.Context, sellerID string, productID string) (Product, error) {
	var product Product
	err := r.db.QueryRow(ctx, `
		SELECT id, seller_id, name, description, status, created_at, updated_at
		FROM products
		WHERE id = $1 AND seller_id = $2 AND status = 'draft'
	`, productID, sellerID).Scan(
		&product.ID,
		&product.SellerID,
		&product.Name,
		&product.Description,
		&product.Status,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Product{}, ErrProductNotFound
	}
	if err != nil {
		return Product{}, err
	}
	return product, nil
}
