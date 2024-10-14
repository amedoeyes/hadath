package service

import (
	"context"
	"testing"
	"time"

	"github.com/amedoeyes/hadath/internal/api/request"
	"github.com/amedoeyes/hadath/internal/database/model"
	"github.com/amedoeyes/hadath/internal/database/repository"
	"github.com/stretchr/testify/assert"
)

func TestEventService(t *testing.T) {
	userRepo := repository.NewUserRepository()
	eventRepo := repository.NewEventRepository()
	service := NewEventService(eventRepo)

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
			name    string
			request request.EventRequest
			wantErr bool
		}{
			{"ValidRequest", validEventRequest, false},
			{"EmptyRequest", request.EventRequest{}, true},
			{"EmptyName", func() request.EventRequest { r := validEventRequest; r.Name = ""; return r }(), true},
			{"EmptyDescription", func() request.EventRequest { r := validEventRequest; r.Description = ""; return r }(), false},
			{"EmptyAddress", func() request.EventRequest { r := validEventRequest; r.Address = ""; return r }(), true},
			{"EmptyStartTime", func() request.EventRequest { r := validEventRequest; r.StartTime = time.Time{}; return r }(), true},
			{"EmptyEndTime", func() request.EventRequest { r := validEventRequest; r.EndTime = time.Time{}; return r }(), true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := service.Create(ctx, &tt.request)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("List", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()
		user := createUser(t, ctx)
		ctx = context.WithValue(ctx, "user", user)

		for i := 0; i < 4; i++ {
			err := service.Create(ctx, &validEventRequest)
			assert.NoError(t, err)
		}

		events, err := service.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, events, 4)
	})

	t.Run("Get", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()
		user := createUser(t, ctx)
		ctx = context.WithValue(ctx, "user", user)

		err := service.Create(ctx, &validEventRequest)
		assert.NoError(t, err)

		events, err := service.List(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, events)

		event, err := eventRepo.Get(ctx, events[0].ID)
		assert.NoError(t, err)
		ctx = context.WithValue(ctx, "event", event)

		res := service.Get(ctx)
		assert.Equal(t, validEventRequest.Name, res.Name)
		assert.Equal(t, validEventRequest.Description, res.Description)
		assert.Equal(t, validEventRequest.Address, res.Address)
		assert.Equal(t, validEventRequest.StartTime.Unix(), res.StartTime.Unix())
		assert.Equal(t, validEventRequest.EndTime.Unix(), res.EndTime.Unix())
	})

	t.Run("Update", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()
		user := createUser(t, ctx)
		ctx = context.WithValue(ctx, "user", user)

		err := service.Create(ctx, &validEventRequest)
		assert.NoError(t, err)

		events, err := service.List(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, events)

		event, err := eventRepo.Get(ctx, events[0].ID)
		assert.NoError(t, err)
		ctx = context.WithValue(ctx, "event", event)

		updatedReq := request.EventRequest{
			Name:        "Catnip Convention 2.0",
			Description: "All cats zero dogs and rats",
			Address:     "321 cat street",
			StartTime:   time.Now().UTC(),
			EndTime:     time.Now().UTC().Add(time.Hour * 2),
		}
		err = service.Update(ctx, &updatedReq)
		assert.NoError(t, err)

		event, err = eventRepo.Get(ctx, events[0].ID)
		assert.NoError(t, err)
		ctx = context.WithValue(ctx, "event", event)

		updatedEvent := service.Get(ctx)
		assert.Equal(t, updatedReq.Name, updatedEvent.Name)
		assert.Equal(t, updatedReq.Description, updatedEvent.Description)
		assert.Equal(t, updatedReq.Address, updatedEvent.Address)
		assert.Equal(t, updatedReq.StartTime.Unix(), updatedEvent.StartTime.Unix())
		assert.Equal(t, updatedReq.EndTime.Unix(), updatedEvent.EndTime.Unix())
	})

	t.Run("Delete", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()
		user := createUser(t, ctx)
		ctx = context.WithValue(ctx, "user", user)

		err := service.Create(ctx, &validEventRequest)
		assert.NoError(t, err)

		events, err := service.List(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, events)

		event, err := eventRepo.Get(ctx, events[0].ID)
		assert.NoError(t, err)
		ctx = context.WithValue(ctx, "event", event)

		err = service.Delete(ctx)
		assert.NoError(t, err)

		events, err = service.List(ctx)
		assert.NoError(t, err)
		assert.Empty(t, events)
	})
}
