package router

import (
	"github.com/amedoeyes/hadath/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Setup() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	userHandler := handler.NewUserHandler()

	r.Post("/users", userHandler.Create)
	r.Get("/users", userHandler.List)
	r.Get("/users/{id}", userHandler.GetById)
	r.Put("/users/{id}", userHandler.Update)
	r.Delete("/users/{id}", userHandler.Delete)

	return r
}
