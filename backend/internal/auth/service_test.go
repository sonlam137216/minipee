package auth

import (
	"context"
	"errors"
	"testing"
	"time"
)

type memorySellerRepository struct {
	byID    map[string]Seller
	byEmail map[string]Seller
}

func (r *memorySellerRepository) CreateSeller(ctx context.Context, seller Seller) (Seller, error) {
	if r.byID == nil {
		r.byID = map[string]Seller{}
	}
	if r.byEmail == nil {
		r.byEmail = map[string]Seller{}
	}
	if _, exists := r.byEmail[seller.Email]; exists {
		return Seller{}, ErrEmailAlreadyExists
	}
	r.byID[seller.ID] = seller
	r.byEmail[seller.Email] = seller
	return seller, nil
}

func (r *memorySellerRepository) FindSellerByEmail(ctx context.Context, email string) (Seller, error) {
	seller, ok := r.byEmail[email]
	if !ok {
		return Seller{}, ErrInvalidCredentials
	}
	return seller, nil
}

func (r *memorySellerRepository) FindSellerByID(ctx context.Context, id string) (Seller, error) {
	seller, ok := r.byID[id]
	if !ok {
		return Seller{}, ErrInvalidCredentials
	}
	return seller, nil
}

func TestRegisterNormalizesEmailAndDoesNotStorePlaintextPassword(t *testing.T) {
	repository := &memorySellerRepository{}
	service := NewService(repository, NewJWTManager("test-secret", time.Hour, fixedTime), fixedID("seller-1"), fixedTime)

	result, err := service.Register(context.Background(), RegisterInput{
		Email:       " Seller@Example.COM ",
		Password:    "password123",
		DisplayName: "Seller",
	})

	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if result.Seller.Email != "seller@example.com" {
		t.Fatalf("Email = %q, want seller@example.com", result.Seller.Email)
	}
	if repository.byEmail["seller@example.com"].PasswordHash == "password123" {
		t.Fatal("password was stored in plaintext")
	}
	if result.AccessToken == "" {
		t.Fatal("AccessToken is empty")
	}
}

func TestLoginWithIncorrectPasswordFails(t *testing.T) {
	repository := &memorySellerRepository{}
	service := NewService(repository, NewJWTManager("test-secret", time.Hour, fixedTime), fixedID("seller-1"), fixedTime)
	_, err := service.Register(context.Background(), RegisterInput{
		Email:       "seller@example.com",
		Password:    "password123",
		DisplayName: "Seller",
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	_, err = service.Login(context.Background(), LoginInput{Email: "seller@example.com", Password: "wrong-password"})

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("error = %v, want ErrInvalidCredentials", err)
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
