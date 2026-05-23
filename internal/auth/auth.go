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

func (s *AuthService) GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, domain.Claims{
		email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})
	return token.SignedString([]byte(s.JWTSecret))
}
