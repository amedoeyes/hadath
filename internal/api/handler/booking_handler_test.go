package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/amedoeyes/hadath/internal/api/middleware"
	"github.com/amedoeyes/hadath/internal/api/response"
	"github.com/amedoeyes/hadath/internal/database/model"
	"github.com/amedoeyes/hadath/internal/database/repository"
	"github.com/amedoeyes/hadath/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBookingHandler(t *testing.T) {
	userRepo := repository.NewUserRepository()
	eventRepo := repository.NewEventRepository()
	bookingRepo := repository.NewBookingRepository()
	bookingService := service.NewBookingService(bookingRepo)
	bookingHandler := NewBookingHandler(bookingService)

	r := chi.NewRouter()
	r.Route("/bookings", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Get("/", bookingHandler.ListByUser)
		})

		r.Route("/event/{id}", func(r chi.Router) {
			r.Use(middleware.EventCtx(eventRepo))
			r.Post("/", bookingHandler.Create)
			r.Get("/", bookingHandler.ListByEvent)
			r.Delete("/", bookingHandler.Delete)
		})
	})

	createUser := func(t *testing.T, ctx context.Context) *model.User {
		err := userRepo.Create(ctx, "Cat", "cat@meow.purr", "meowmeowmeow")
		assert.NoError(t, err)
		user, err := userRepo.GetByEmail(ctx, "cat@meow.purr")
		assert.NoError(t, err)
		return user
	}

	createEvent := func(t *testing.T, ctx context.Context, userID uuid.UUID) *model.Event {
		err := eventRepo.Create(ctx, userID, "Catnip Convention", "All cats zero dogs", "123 cat street", time.Now().UTC(), time.Now().UTC().Add(time.Hour))
		assert.NoError(t, err)
		events, err := eventRepo.List(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, events)
		return &events[0]
	}

	t.Run("Create", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()
		user := createUser(t, ctx)
		ctx = context.WithValue(ctx, "user", user)
		event := createEvent(t, ctx, user.ID)
		ctx = context.WithValue(ctx, "event", event)

		req := httptest.NewRequest(http.MethodPost, "/bookings/event/"+event.ID.String(), nil)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)

		req = httptest.NewRequest(http.MethodPost, "/bookings/event/"+event.ID.String(), nil)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/json")

		rr = httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
	})

	t.Run("ListByUser", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()
		user := createUser(t, ctx)
		ctx = context.WithValue(ctx, "user", user)
		event := createEvent(t, ctx, user.ID)
		ctx = context.WithValue(ctx, "event", event)

		err := bookingService.Create(ctx)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/bookings/user/", nil)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var bookings []response.EventResponse
		err = json.NewDecoder(rr.Body).Decode(&bookings)
		assert.NoError(t, err)
		assert.NotEmpty(t, bookings)
	})

	t.Run("ListByEvent", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()
		user := createUser(t, ctx)
		ctx = context.WithValue(ctx, "user", user)
		event := createEvent(t, ctx, user.ID)
		ctx = context.WithValue(ctx, "event", event)

		err := bookingService.Create(ctx)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/bookings/event/"+event.ID.String(), nil)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var bookings []response.UserResponse
		err = json.NewDecoder(rr.Body).Decode(&bookings)
		assert.NoError(t, err)
		assert.NotEmpty(t, bookings)
	})

	t.Run("Delete", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()
		user := createUser(t, ctx)
		ctx = context.WithValue(ctx, "user", user)
		event := createEvent(t, ctx, user.ID)
		ctx = context.WithValue(ctx, "event", event)

		err := bookingService.Create(ctx)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodDelete, "/bookings/event/"+event.ID.String(), nil)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})
}
