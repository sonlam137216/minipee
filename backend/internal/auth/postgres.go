package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateSeller(ctx context.Context, seller Seller) (Seller, error) {
	_, err := r.db.Exec(ctx, `
		INSERT INTO sellers (id, email, password_hash, display_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, seller.ID, seller.Email, seller.PasswordHash, seller.DisplayName, seller.CreatedAt, seller.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return Seller{}, ErrEmailAlreadyExists
		}
		return Seller{}, err
	}
	return seller, nil
}

func (r *PostgresRepository) FindSellerByEmail(ctx context.Context, email string) (Seller, error) {
	return r.findOne(ctx, `
		SELECT id, email, password_hash, display_name, created_at, updated_at
		FROM sellers
		WHERE email = $1
	`, email)
}

func (r *PostgresRepository) FindSellerByID(ctx context.Context, id string) (Seller, error) {
	return r.findOne(ctx, `
		SELECT id, email, password_hash, display_name, created_at, updated_at
		FROM sellers
		WHERE id = $1
	`, id)
}

func (r *PostgresRepository) findOne(ctx context.Context, query string, arg string) (Seller, error) {
	var seller Seller
	err := r.db.QueryRow(ctx, query, arg).Scan(
		&seller.ID,
		&seller.Email,
		&seller.PasswordHash,
		&seller.DisplayName,
		&seller.CreatedAt,
		&seller.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Seller{}, ErrInvalidCredentials
	}
	if err != nil {
		return Seller{}, err
	}
	return seller, nil
}
