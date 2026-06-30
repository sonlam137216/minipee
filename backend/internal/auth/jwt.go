package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	SellerID string `json:"seller_id"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret []byte
	ttl    time.Duration
	now    func() time.Time
}

func NewJWTManager(secret string, ttl time.Duration, now func() time.Time) *JWTManager {
	return &JWTManager{secret: []byte(secret), ttl: ttl, now: now}
}

func (m *JWTManager) Issue(seller Seller) (string, error) {
	now := m.now().UTC()
	claims := Claims{
		SellerID: seller.ID,
		Email:    seller.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   seller.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTManager) Verify(rawToken string) (Claims, error) {
	token, err := jwt.ParseWithClaims(rawToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidAccessToken
		}
		return m.secret, nil
	}, jwt.WithExpirationRequired())
	if err != nil {
		return Claims{}, ErrInvalidAccessToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid || claims.SellerID == "" {
		return Claims{}, ErrInvalidAccessToken
	}
	return *claims, nil
}
