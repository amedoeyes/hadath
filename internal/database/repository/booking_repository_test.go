package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBookingRepository(t *testing.T) {
	bookingRepo := NewBookingRepository()
	userRepo := NewUserRepository()
	eventRepo := NewEventRepository()
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		setupTest(t)

		user := setupUser(t, userRepo, ctx)
		event := setupEvent(t, eventRepo, ctx, user.ID)

		err := bookingRepo.Create(ctx, user.ID, event.ID)
		assert.NoError(t, err)

		err = bookingRepo.Create(ctx, user.ID, event.ID)
		assert.Error(t, err)
	})

	t.Run("ListByUser", func(t *testing.T) {
		setupTest(t)

		user := setupUser(t, userRepo, ctx)
		event := setupEvent(t, eventRepo, ctx, user.ID)

		err := bookingRepo.Create(ctx, user.ID, event.ID)
		assert.NoError(t, err)

		events, err := bookingRepo.ListByUser(ctx, user.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, events)
		assert.Equal(t, *event, events[0])
	})

	t.Run("ListByEvent", func(t *testing.T) {
		setupTest(t)

		user := setupUser(t, userRepo, ctx)
		event := setupEvent(t, eventRepo, ctx, user.ID)

		err := bookingRepo.Create(ctx, user.ID, event.ID)
		assert.NoError(t, err)

		users, err := bookingRepo.ListByEvent(ctx, event.ID)
		assert.NoError(t, err)
		assert.NotEmpty(t, users)
		assert.Equal(t, *user, users[0])
	})

	t.Run("Delete", func(t *testing.T) {
		setupTest(t)

		user := setupUser(t, userRepo, ctx)
		event := setupEvent(t, eventRepo, ctx, user.ID)

		err := bookingRepo.Create(ctx, user.ID, event.ID)
		assert.NoError(t, err)

		err = bookingRepo.Delete(ctx, user.ID, event.ID)
		assert.NoError(t, err)

		events, err := bookingRepo.ListByUser(ctx, user.ID)
		assert.NoError(t, err)
		assert.Empty(t, events)
	})
}
