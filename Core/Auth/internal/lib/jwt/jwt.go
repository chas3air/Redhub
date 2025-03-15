package jwt

import (
	"auth/internal/domain/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewTokens(user models.User, accessDuration, refreshDuration time.Duration) (string, string, error) {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["uid"] = user.Id
	accessClaims["role"] = user.Role
	accessClaims["exp"] = time.Now().Add(accessDuration).Unix()

	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("AUTH_SECRET")))
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["uid"] = user.Id
	refreshClaims["exp"] = time.Now().Add(refreshDuration).Unix()

	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("AUTH_SECRET")))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
