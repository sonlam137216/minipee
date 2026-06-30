package auth

import (
	"context"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	CreateSeller(ctx context.Context, seller Seller) (Seller, error)
	FindSellerByEmail(ctx context.Context, email string) (Seller, error)
	FindSellerByID(ctx context.Context, id string) (Seller, error)
}

type TokenIssuer interface {
	Issue(seller Seller) (string, error)
}

type Service struct {
	repository Repository
	tokens     TokenIssuer
	newID      func() string
	now        func() time.Time
}

func NewService(repository Repository, tokens TokenIssuer, newID func() string, now func() time.Time) *Service {
	return &Service{repository: repository, tokens: tokens, newID: newID, now: now}
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (AuthResult, error) {
	email, err := NormalizeEmail(input.Email)
	if err != nil {
		return AuthResult{}, err
	}
	displayName := strings.TrimSpace(input.DisplayName)
	if displayName == "" || len(input.Password) < 8 {
		return AuthResult{}, ErrInvalidAuthInput
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResult{}, err
	}
	now := s.now().UTC()
	seller := Seller{
		ID:           s.newID(),
		Email:        email,
		DisplayName:  displayName,
		PasswordHash: string(passwordHash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	seller, err = s.repository.CreateSeller(ctx, seller)
	if err != nil {
		return AuthResult{}, err
	}
	token, err := s.tokens.Issue(seller)
	if err != nil {
		return AuthResult{}, err
	}
	return AuthResult{Seller: seller, AccessToken: token}, nil
}

func (s *Service) Login(ctx context.Context, input LoginInput) (AuthResult, error) {
	email, err := NormalizeEmail(input.Email)
	if err != nil {
		return AuthResult{}, ErrInvalidCredentials
	}
	seller, err := s.repository.FindSellerByEmail(ctx, email)
	if err != nil {
		return AuthResult{}, ErrInvalidCredentials
	}
	if bcrypt.CompareHashAndPassword([]byte(seller.PasswordHash), []byte(input.Password)) != nil {
		return AuthResult{}, ErrInvalidCredentials
	}
	token, err := s.tokens.Issue(seller)
	if err != nil {
		return AuthResult{}, err
	}
	return AuthResult{Seller: seller, AccessToken: token}, nil
}

func (s *Service) FindSeller(ctx context.Context, sellerID string) (Seller, error) {
	return s.repository.FindSellerByID(ctx, sellerID)
}
