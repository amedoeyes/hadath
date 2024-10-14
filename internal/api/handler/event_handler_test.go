package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/amedoeyes/hadath/internal/api/middleware"
	"github.com/amedoeyes/hadath/internal/api/request"
	"github.com/amedoeyes/hadath/internal/database/model"
	"github.com/amedoeyes/hadath/internal/database/repository"
	"github.com/amedoeyes/hadath/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestEventHandler(t *testing.T) {
	userRepo := repository.NewUserRepository()
	eventRepo := repository.NewEventRepository()
	eventService := service.NewEventService(eventRepo)
	eventHandler := NewEventHandler(eventService)

	r := chi.NewRouter()
	r.Route("/events", func(r chi.Router) {
		r.Post("/", eventHandler.Create)
		r.Get("/", eventHandler.List)
		r.Route("/{id}", func(r chi.Router) {
			r.Use(middleware.EventCtx(eventRepo))
			r.Get("/", eventHandler.Get)
			r.Put("/", eventHandler.Update)
			r.Delete("/", eventHandler.Delete)
		})
	})

	createUser := func(t *testing.T, ctx context.Context) *model.User {
		err := userRepo.Create(ctx, "Cat", "cat@meow.purr", "meowmeowmeow")
		assert.NoError(t, err)
		user, err := userRepo.GetByEmail(ctx, "cat@meow.purr")
		assert.NoError(t, err)
		return user
	}

	validEventRequest := request.EventRequest{
		Name:        "Catnip Convention",
		Description: "All cats zero dogs",
		Address:     "123 cat street",
		StartTime:   time.Now().UTC(),
		EndTime:     time.Now().UTC().Add(time.Hour),
	}

	t.Run("Create", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()
		user := createUser(t, ctx)
		ctx = context.WithValue(ctx, "user", user)

		tests := []struct {
			name       string
			request    request.EventRequest
			wantStatus int
		}{
			{"ValidRequest", validEventRequest, http.StatusCreated},
			{"EmptyRequest", request.EventRequest{}, http.StatusBadRequest},
			{"EmptyName", func() request.EventRequest { r := validEventRequest; r.Name = ""; return r }(), http.StatusBadRequest},
			{"EmptyDescription", func() request.EventRequest { r := validEventRequest; r.Description = ""; return r }(), http.StatusCreated},
			{"EmptyAddress", func() request.EventRequest { r := validEventRequest; r.Address = ""; return r }(), http.StatusBadRequest},
			{"EmptyStartTime", func() request.EventRequest { r := validEventRequest; r.StartTime = time.Time{}; return r }(), http.StatusBadRequest},
			{"EmptyEndTime", func() request.EventRequest { r := validEventRequest; r.EndTime = time.Time{}; return r }(), http.StatusBadRequest},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				reqBody, _ := json.Marshal(tt.request)
				req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(reqBody))
				req = req.WithContext(ctx)
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				r.ServeHTTP(rr, req)

				assert.Equal(t, tt.wantStatus, rr.Code)
			})
		}
	})

	t.Run("List", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()
		user := createUser(t, ctx)
		ctx = context.WithValue(ctx, "user", user)

		for i := 0; i < 4; i++ {
			err := eventService.Create(ctx, &validEventRequest)
			assert.NoError(t, err)
		}

		req := httptest.NewRequest(http.MethodGet, "/events", nil)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var events []request.EventRequest
		err := json.NewDecoder(rr.Body).Decode(&events)
		assert.NoError(t, err)
		assert.Len(t, events, 4)
	})

	t.Run("Get", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()
		user := createUser(t, ctx)
		ctx = context.WithValue(ctx, "user", user)

		err := eventService.Create(ctx, &validEventRequest)
		assert.NoError(t, err)

		events, err := eventRepo.List(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, events)

		eventID := events[0].ID
		req := httptest.NewRequest(http.MethodGet, "/events/"+eventID.String(), nil)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var event request.EventRequest
		err = json.NewDecoder(rr.Body).Decode(&event)
		assert.NoError(t, err)
		assert.Equal(t, validEventRequest.Name, event.Name)
	})

	t.Run("Update", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()
		user := createUser(t, ctx)
		ctx = context.WithValue(ctx, "user", user)

		err := eventService.Create(ctx, &validEventRequest)
		assert.NoError(t, err)

		events, err := eventRepo.List(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, events)

		updatedReq := request.EventRequest{
			Name:        "Catnip Convention 2.0",
			Description: "All cats zero dogs and rats",
			Address:     "321 cat street",
			StartTime:   time.Now().UTC(),
			EndTime:     time.Now().UTC().Add(time.Hour * 2),
		}

		eventID := events[0].ID
		reqBody, _ := json.Marshal(updatedReq)
		req := httptest.NewRequest(http.MethodPut, "/events/"+eventID.String(), bytes.NewBuffer(reqBody))
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)

		req = httptest.NewRequest(http.MethodGet, "/events/"+eventID.String(), nil)
		req = req.WithContext(ctx)
		rr = httptest.NewRecorder()

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		var event request.EventRequest
		err = json.NewDecoder(rr.Body).Decode(&event)
		assert.NoError(t, err)
		assert.Equal(t, updatedReq.Name, event.Name)
	})

	t.Run("Delete", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()
		user := createUser(t, ctx)
		ctx = context.WithValue(ctx, "user", user)

		err := eventService.Create(ctx, &validEventRequest)
		assert.NoError(t, err)

		events, err := eventRepo.List(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, events)

		eventID := events[0].ID
		req := httptest.NewRequest(http.MethodDelete, "/events/"+eventID.String(), nil)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		r.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)

		req = httptest.NewRequest(http.MethodGet, "/events/"+eventID.String(), nil)
		req = req.WithContext(ctx)
		rr = httptest.NewRecorder()

		r.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
}
