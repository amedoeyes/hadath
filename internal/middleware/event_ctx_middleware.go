package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/amedoeyes/hadath/internal/api/response"
	"github.com/amedoeyes/hadath/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func EventCtx(repo *repository.EventRepository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := uuid.Parse(chi.URLParam(r, "id"))
			fmt.Printf("id: %v\n", id)
			if err != nil {
				response.HandleError(w, err)
				return
			}
			event, err := repo.Get(r.Context(), id)
			if err != nil {
				response.WriteJSONError(w, http.StatusNotFound, "Not found", nil)
				return
			}
			ctx := context.WithValue(r.Context(), "event", event)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
