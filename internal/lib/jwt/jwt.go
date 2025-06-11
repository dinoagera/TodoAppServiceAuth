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

// func ParseToken(tokenString string) (*jwt.Token, error) {
// 	secretKey := config.GetConfig().JWTSecret
// 	if strings.Count(tokenString, ".") != 2 {
// 		return nil, fmt.Errorf("invalid token format")
// 	}
// 	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(secretKey), nil
// 	})
// }
