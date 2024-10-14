package service

import (
	"context"
	"testing"
	"time"

	"github.com/amedoeyes/hadath/internal/database/model"
	"github.com/amedoeyes/hadath/internal/database/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBookingService(t *testing.T) {
	userRepo := repository.NewUserRepository()
	eventRepo := repository.NewEventRepository()
	bookingRepo := repository.NewBookingRepository()
	service := NewBookingService(bookingRepo)

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
		event := createEvent(t, ctx, user.ID)

		ctx = context.WithValue(ctx, "user", user)
		ctx = context.WithValue(ctx, "event", event)

		err := service.Create(ctx)
		assert.NoError(t, err)

		err = service.Create(ctx)
		assert.Error(t, err)
	})

	t.Run("ListByUser", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()

		user := createUser(t, ctx)
		event := createEvent(t, ctx, user.ID)

		ctx = context.WithValue(ctx, "user", user)
		ctx = context.WithValue(ctx, "event", event)

		err := service.Create(ctx)
		assert.NoError(t, err)

		bookings, err := service.ListByUser(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, bookings)
	})

	t.Run("ListByEvent", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()

		user := createUser(t, ctx)
		event := createEvent(t, ctx, user.ID)

		ctx = context.WithValue(ctx, "user", user)
		ctx = context.WithValue(ctx, "event", event)

		err := service.Create(ctx)
		assert.NoError(t, err)

		bookings, err := service.ListByEvent(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, bookings)
	})

	t.Run("Delete", func(t *testing.T) {
		setupTest(t)
		ctx := context.Background()

		user := createUser(t, ctx)
		event := createEvent(t, ctx, user.ID)

		ctx = context.WithValue(ctx, "user", user)
		ctx = context.WithValue(ctx, "event", event)

		err := service.Delete(ctx)
		assert.NoError(t, err)

		userBookings, err := service.ListByUser(ctx)
		assert.NoError(t, err)
		assert.Empty(t, userBookings)

		eventBookings, err := service.ListByEvent(ctx)
		assert.NoError(t, err)
		assert.Empty(t, eventBookings)
	})
}
