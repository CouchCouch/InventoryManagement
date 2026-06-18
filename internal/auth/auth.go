package auth

import (
	"errors"
	"time"

	"inventory/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	jwtSecret        string
	jwtRefreshSecret string
}

func NewAuthService(jwtSecret, jwtRefreshSecret string) *AuthService {
	return &AuthService{
		jwtSecret:        jwtSecret,
		jwtRefreshSecret: jwtRefreshSecret,
	}
}

func (s *AuthService) GenerateAccessToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, domain.Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(24 * time.Hour)),
		},
	})
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) GenerateRefreshToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, domain.Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(24 * time.Hour)),
		},
	})
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*domain.Claims); ok && token.Valid {
		return claims.Email, nil
	}

	return "", errors.New("invalid token")
}
