package tokenparser

import (
	"apigateway/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ParseToken(tokenString string) (*models.Claims, error) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return models.Claims{}, nil
	})

	claims := token.Claims.(jwt.MapClaims)
	exp := time.Unix(int64(claims["exp"].(float64)), 0)

	customClaims := &models.Claims{
		Uid:  claims["uid"].(string),
		Role: claims["role"].(string),
		Exp:  exp,
	}

	return customClaims, nil
}
