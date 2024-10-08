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
	eventRepo := repository.NewEventRepository()

	authService := service.NewAuthService(apiKeyRepo, userRepo)
	eventService := service.NewEventService(eventRepo)

	authHandler := handler.NewAuthHandler(authService)
	eventHandler := handler.NewEventHandler(eventService)

	authMiddleware := middleware.Auth(apiKeyRepo)
	userCtxMiddleware := middleware.UserCtx(userRepo)
	eventCtxMiddleware := middleware.EventCtx(eventRepo)

	r.Use(chiMiddleware.Logger)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authHandler.SignUp)
		r.Post("/signin", authHandler.SignIn)
		r.With(authMiddleware).Post("/signout", authHandler.SignOut)
	})

	r.Route("/events", func(r chi.Router) {
		r.Use(authMiddleware)
		r.Use(userCtxMiddleware)
		r.Post("/", eventHandler.Create)
		r.Get("/", eventHandler.GetAll)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(eventCtxMiddleware)
			r.Get("/", eventHandler.GetByID)
			r.Put("/", eventHandler.UpdateByID)
			r.Delete("/", eventHandler.DeleteByID)
		})
	})

	return r
}
