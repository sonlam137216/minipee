package auth

import (
	"errors"
	"net/mail"
	"strings"
	"time"
)

var (
	ErrInvalidAuthInput   = errors.New("invalid authentication input")
	ErrEmailAlreadyExists = errors.New("email is already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUnauthenticated    = errors.New("authentication required")
	ErrInvalidAccessToken = errors.New("invalid access token")
)

type Seller struct {
	ID           string
	Email        string
	DisplayName  string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RegisterInput struct {
	Email       string
	Password    string
	DisplayName string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResult struct {
	Seller      Seller
	AccessToken string
}

func NormalizeEmail(email string) (string, error) {
	parsed, err := mail.ParseAddress(strings.TrimSpace(email))
	if err != nil {
		return "", ErrInvalidAuthInput
	}
	normalized := strings.ToLower(parsed.Address)
	if normalized == "" {
		return "", ErrInvalidAuthInput
	}
	return normalized, nil
}
