package middleware

import (
	tokenparser "apigateway/internal/lib/jwt/tokenParser"
	"net/http"
	"os"
	"strings"
	"time"
)

type Middleware struct {
}

func New() *Middleware {
	return &Middleware{}
}

// Middleware для проверки токена
func (m *Middleware) ValidateToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))

		secretKey := []byte(os.Getenv("AUTH_SECRET"))
		claims, err := tokenparser.ParseToken(tokenString, secretKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims.Exp.Before(time.Now()) {
			http.Error(w, "Token has expired", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
