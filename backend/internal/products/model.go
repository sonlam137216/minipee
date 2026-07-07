package products

import (
	"errors"
	"strings"
	"time"
)

const (
	StatusDraft     = "draft"
	StatusPublished = "published"
)

var (
	ErrInvalidProductName      = errors.New("product name must contain between 3 and 200 characters")
	ErrProductNotFound         = errors.New("product not found")
	ErrProductAlreadyPublished = errors.New("product is already published")
)

type Product struct {
	ID          string
	SellerID    string
	Name        string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateProductInput struct {
	Name        string
	Description string
	Status      string
}

func NewDraftProduct(id string, sellerID string, input CreateProductInput, now time.Time) (Product, error) {
	name := strings.TrimSpace(input.Name)
	if len([]rune(name)) < 3 || len([]rune(name)) > 200 {
		return Product{}, ErrInvalidProductName
	}
	return Product{
		ID:          id,
		SellerID:    sellerID,
		Name:        name,
		Description: strings.TrimSpace(input.Description),
		Status:      StatusDraft,
		CreatedAt:   now.UTC(),
		UpdatedAt:   now.UTC(),
	}, nil
}
