package auth

import (
	"time"

	"inventory/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	JWTSecret        string
	jwtRefreshSecret string
}

func NewAuthService(jwtSecret, jwtRefreshSecret string) *AuthService {
	return &AuthService{
		JWTSecret:        jwtSecret,
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
	return token.SignedString([]byte(s.JWTSecret))
}

func (s *AuthService) GenerateRefreshToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, domain.Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(24 * time.Hour)),
		},
	})
	return token.SignedString([]byte(s.JWTSecret))
}
