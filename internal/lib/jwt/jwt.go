package jwt

import (
	"fmt"
	"time"

	"github.com/dinoagera/api-auth/internal/config"
	"github.com/dinoagera/api-auth/internal/domain/models"
	"github.com/golang-jwt/jwt"
)

func NewToken(user models.User, duration time.Duration) (string, error) {
	secretKey := config.GetConfig().JWTSecret
	if len(secretKey) < 32 {
		return "", jwt.ErrInvalidKey
	}

	claims := jwt.MapClaims{
		"uid":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(duration).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
