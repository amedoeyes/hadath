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
	bookingRepo := repository.NewBookingRepository()

	authService := service.NewAuthService(apiKeyRepo, userRepo)
	eventService := service.NewEventService(eventRepo)
	bookingService := service.NewBookingService(bookingRepo)

	authHandler := handler.NewAuthHandler(authService)
	eventHandler := handler.NewEventHandler(eventService)
	bookingHandler := handler.NewBookingHandler(bookingService)

	authMiddleware := middleware.Auth(apiKeyRepo)
	currentUserCtxMiddleware := middleware.CurrentUserCtx(userRepo)
	eventCtxMiddleware := middleware.EventCtx(eventRepo)

	r.Use(chiMiddleware.Logger)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authHandler.SignUp)
		r.Post("/signin", authHandler.SignIn)
		r.With(authMiddleware).Post("/signout", authHandler.SignOut)
	})

	r.Route("/events", func(r chi.Router) {
		r.Use(authMiddleware)

		r.With(currentUserCtxMiddleware).Post("/", eventHandler.Create)
		r.Get("/", eventHandler.List)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(eventCtxMiddleware)
			r.Get("/", eventHandler.Get)
			r.With(currentUserCtxMiddleware).Put("/", eventHandler.Update)
			r.With(currentUserCtxMiddleware).Delete("/", eventHandler.Delete)
		})
	})

	r.Route("/bookings", func(r chi.Router) {
		r.Use(authMiddleware)

		r.Route("/user", func(r chi.Router) {
			r.Use(currentUserCtxMiddleware)
			r.Post("/", bookingHandler.Create)
			r.Get("/", bookingHandler.ListByCurrentUser)
			r.Delete("/", bookingHandler.Delete)
		})

		r.With(eventCtxMiddleware).Get("/event/{id}", bookingHandler.ListByEvent)
	})

	return r
}
