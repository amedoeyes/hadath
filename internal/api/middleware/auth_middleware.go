package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/amedoeyes/hadath/internal/api/response"
	"github.com/amedoeyes/hadath/internal/database/repository"
	"github.com/amedoeyes/hadath/internal/service"
)

func Auth(repo *repository.APIKeyRepository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				response.WriteJSONError(w, http.StatusUnauthorized, "Invalid API key", nil)
				return
			}
			key := strings.TrimPrefix(authHeader, "Bearer ")

			apiKey, err := repo.GetByKey(r.Context(), service.HashAPIKey(key))
			if err != nil {
				response.WriteJSONError(w, http.StatusUnauthorized, "Invalid API key", nil)
				return
			}

			ctx := context.WithValue(r.Context(), "apiKey", apiKey)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
