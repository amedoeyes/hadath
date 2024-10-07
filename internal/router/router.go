package router

import (
	"github.com/amedoeyes/hadath/internal/handler"
	"github.com/amedoeyes/hadath/internal/middleware"
	"github.com/amedoeyes/hadath/internal/repository"
	"github.com/amedoeyes/hadath/internal/service"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func Setup() *chi.Mux {
	r := chi.NewRouter()

	apiKeyRepo := repository.NewAPIKeyRepository()
	userRepo := repository.NewUserRepository()

	authService := service.NewAuthService(apiKeyRepo, userRepo)

	authHandler := handler.NewAuthHandler(authService)

	authMiddleware := middleware.Auth(apiKeyRepo)

	r.Use(chiMiddleware.Logger)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authHandler.SignUp)
		r.Post("/signin", authHandler.SignIn)
		r.With(authMiddleware).Post("/signout", authHandler.SignOut)
	})

	return r
}
