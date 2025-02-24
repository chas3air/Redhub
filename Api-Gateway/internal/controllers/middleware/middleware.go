package middleware

import (
	usersmanagerservice "apigateway/internal/services/usersManager"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Claims struct {
	Uid  string `json:"uid"`
	Role string `json:"role"`
	Aid  string `json:"aid"`
	Exp  int64  `json:"exp"`
	jwt.StandardClaims
}

func ValidateToken(next http.Handler, user_service *usersmanagerservice.UsersManager) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer")
		claims, err := parseToken(tokenString, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		uuidUID, err := uuid.Parse(claims.Aid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := user_service.GetUserById(r.Context(), uuidUID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if user.Role != claims.Role {
			http.Error(w, "invalid token", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func parseToken(tokenString string, secretKey []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
