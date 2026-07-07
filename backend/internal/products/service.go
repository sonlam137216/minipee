package products

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, product Product) (Product, error)
	ListBySeller(ctx context.Context, sellerID string) ([]Product, error)
	FindBySellerID(ctx context.Context, sellerID string, productID string) (Product, error)
	PublishDraftBySeller(ctx context.Context, sellerID string, productID string, now time.Time) (Product, error)
	ListPublished(ctx context.Context) ([]Product, error)
	FindPublishedByID(ctx context.Context, productID string) (Product, error)
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

func (s *Service) GetSellerProduct(ctx context.Context, sellerID string, productID string) (Product, error) {
	return s.repository.FindBySellerID(ctx, sellerID, productID)
}

func (s *Service) PublishSellerProduct(ctx context.Context, sellerID string, productID string) (Product, error) {
	product, err := s.repository.FindBySellerID(ctx, sellerID, productID)
	if err != nil {
		return Product{}, err
	}
	if product.Status == StatusPublished {
		return Product{}, ErrProductAlreadyPublished
	}
	return s.repository.PublishDraftBySeller(ctx, sellerID, productID, s.now().UTC())
}

func (s *Service) ListPublishedProducts(ctx context.Context) ([]Product, error) {
	return s.repository.ListPublished(ctx)
}

func (s *Service) GetPublishedProduct(ctx context.Context, productID string) (Product, error) {
	return s.repository.FindPublishedByID(ctx, productID)
}
