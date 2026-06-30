package products

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, product Product) (Product, error)
	ListBySeller(ctx context.Context, sellerID string) ([]Product, error)
	FindSellerDraftByID(ctx context.Context, sellerID string, productID string) (Product, error)
}

type Service struct {
	repository Repository
	newID      func() string
	now        func() time.Time
}

func NewService(repository Repository, newID func() string, now func() time.Time) *Service {
	return &Service{repository: repository, newID: newID, now: now}
}

func (s *Service) CreateDraft(ctx context.Context, sellerID string, input CreateProductInput) (Product, error) {
	product, err := NewDraftProduct(s.newID(), sellerID, input, s.now())
	if err != nil {
		return Product{}, err
	}
	return s.repository.Create(ctx, product)
}

func (s *Service) ListSellerProducts(ctx context.Context, sellerID string) ([]Product, error) {
	return s.repository.ListBySeller(ctx, sellerID)
}

func (s *Service) GetSellerDraft(ctx context.Context, sellerID string, productID string) (Product, error) {
	return s.repository.FindSellerDraftByID(ctx, sellerID, productID)
}
