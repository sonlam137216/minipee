package auth

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresSellerRepositoryAndLogin(t *testing.T) {
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

	repository := NewPostgresRepository(db)
	service := NewService(repository, NewJWTManager("test-secret", time.Hour, fixedTime), fixedID("00000000-0000-4000-8000-000000000101"), fixedTime)
	email := "seller-auth-" + time.Now().Format("150405.000000000") + "@example.com"

	result, err := service.Register(ctx, RegisterInput{Email: email, Password: "password123", DisplayName: "Seller"})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if result.Seller.PasswordHash == "password123" {
		t.Fatal("password was stored in plaintext")
	}

	_, err = service.Register(ctx, RegisterInput{Email: email, Password: "password123", DisplayName: "Seller"})
	if err != ErrEmailAlreadyExists {
		t.Fatalf("duplicate register error = %v, want ErrEmailAlreadyExists", err)
	}

	loginResult, err := service.Login(ctx, LoginInput{Email: email, Password: "password123"})
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}
	if loginResult.AccessToken == "" {
		t.Fatal("Login access token is empty")
	}

	_, err = service.Login(ctx, LoginInput{Email: email, Password: "wrong-password"})
	if err != ErrInvalidCredentials {
		t.Fatalf("wrong password error = %v, want ErrInvalidCredentials", err)
	}
}
