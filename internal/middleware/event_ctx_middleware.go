package middleware

import (
	"context"
	"net/http"

	"github.com/amedoeyes/hadath/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func EventCtx(repo *repository.EventRepository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := uuid.Parse(chi.URLParam(r, "id"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			event, err := repo.Get(r.Context(), id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "event", event)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
