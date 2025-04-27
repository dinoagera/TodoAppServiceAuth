package jwt

import (
	"time"

	"github.com/dinoagera/api-auth/internal/config"
	"github.com/dinoagera/api-auth/internal/domain/models"
	"github.com/golang-jwt/jwt"
)

func NewToken(user models.User, duration time.Duration) (string, error) {
	sekretKey := config.GetConfig().JWTSecret
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Username
	claims["exp"] = time.Now().Add(duration).Unix()
	tokenString, err := token.SignedString([]byte(sekretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
