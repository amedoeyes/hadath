package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/amedoeyes/hadath/internal/repository"
	"github.com/go-chi/chi/v5"
)

func EventCtx(repo *repository.EventRepository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			event, err := repo.GetByID(r.Context(), uint32(id))
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "event", event)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
