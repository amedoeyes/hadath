package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/amedoeyes/hadath/internal/repository"
	"github.com/amedoeyes/hadath/internal/utility"
)

func Auth(repo *repository.APIKeyRepository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}
			key := strings.TrimPrefix(authHeader, "Bearer ")

			apiKey, err := repo.GetByKey(r.Context(), utility.HashAPIKey(key))
			if err != nil {
				http.Error(w, "invalid API key", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "apiKey", apiKey)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
