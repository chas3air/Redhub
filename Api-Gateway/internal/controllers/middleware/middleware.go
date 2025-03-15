package middleware

import (
	"apigateway/internal/domain/models"
	tokenparser "apigateway/internal/lib/jwt/tokenParser"
	"apigateway/internal/storage/cache"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Middleware struct {
}

func New() *Middleware {
	return &Middleware{}
}

func (m *Middleware) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))

		claims, err := tokenparser.ParseToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if claims.Exp.Before(time.Now()) {
			parsedUIDToUUID, _ := uuid.Parse(claims.Uid)
			cache := cache.New()
			cache.Delete(parsedUIDToUUID)

			http.Error(w, "Token has expired", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) RequireUser(next http.Handler) http.Handler {
	return m.roleMiddleware("user", next)
}

func (m *Middleware) RequireUserAdmin(next http.Handler) http.Handler {
	return m.roleMiddleware("user_admin", next)
}

func (m *Middleware) RequireArticleAdmin(next http.Handler) http.Handler {
	return m.roleMiddleware("article_admin", next)
}

func (m *Middleware) roleMiddleware(requiredRole string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value("claims").(*models.Claims)

		if claims.Role != requiredRole {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) PreventAccessIfLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))

			_, err := tokenparser.ParseToken(tokenString)
			if err == nil {
				http.Error(w, "Already logged in", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
