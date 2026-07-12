package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
	ErrWrongTokenType = errors.New("unexpected token type")
)


type Claims struct {
	UserID int64     `json:"user_id"`
	Role   string    `json:"role"`
	Type   TokenType `json:"type"`
	jwt.RegisteredClaims
}
// --- convenience wrappers, added for handler/middleware use ---

// GenerateAccessToken issues a short-lived token used to authenticate API calls.
func GenerateAccessToken(secret string, userID int64, role string, ttl time.Duration) (string, error) {
	return GenerateToken(secret, userID, role, AccessToken, ttl)
}

// GenerateRefreshToken issues a long-lived token, meant to live in an
// HttpOnly cookie and be exchanged for a new access token.
func GenerateRefreshToken(secret string, userID int64, role string, ttl time.Duration) (string, error) {
	return GenerateToken(secret, userID, role, RefreshToken, ttl)
}

// ParseAccessToken validates a token presented in the Authorization header.
func ParseAccessToken(secret, tokenString string) (*Claims, error) {
	return ParseToken(secret, tokenString, AccessToken)
}

// ParseRefreshToken validates a token presented via the refresh cookie.
func ParseRefreshToken(secret, tokenString string) (*Claims, error) {
	return ParseToken(secret, tokenString, RefreshToken)
}

func GenerateToken(secret string, userID int64, role string, tokenType TokenType, ttl time.Duration) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}


func ParseToken(secret, tokenString string, wantType TokenType) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.Type != wantType {
		return nil, ErrWrongTokenType
	}

	return claims, nil
}