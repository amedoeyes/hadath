package middleware

import (
	"context"
	"net/http"

	"github.com/amedoeyes/hadath/internal/model"
	"github.com/amedoeyes/hadath/internal/repository"
)

func UserCtx(repo *repository.UserRepository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Context().Value("apiKey").(*model.APIKey)
			user, err := repo.GetByID(r.Context(), apiKey.UserID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
