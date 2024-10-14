package repository

import (
	"context"
	"testing"
	"time"

	"github.com/amedoeyes/hadath/internal/database/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var eventData = model.Event{
	Name:        "Test Event",
	Description: "Test Description",
	Address:     "Test Address",
	StartTime:   time.Now().UTC(),
	EndTime:     time.Now().UTC().Add(time.Hour),
}

func setupEvent(t *testing.T, repo *EventRepository, ctx context.Context, userID uuid.UUID) *model.Event {
	t.Helper()
	err := repo.Create(ctx, userID, eventData.Name, eventData.Description, eventData.Address, eventData.StartTime, eventData.EndTime)
	assert.NoError(t, err)
	events, err := repo.List(ctx)
	assert.NoError(t, err)
	event := events[0]
	return &event
}

func TestEventRepository(t *testing.T) {
	eventRepo := NewEventRepository()
	userRepo := NewUserRepository()
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		setupTest(t)

		user := setupUser(t, userRepo, ctx)
		err := eventRepo.Create(ctx, user.ID, eventData.Name, eventData.Description, eventData.Address, eventData.StartTime, eventData.EndTime)
		assert.NoError(t, err)
	})

	t.Run("List", func(t *testing.T) {
		setupTest(t)

		user := setupUser(t, userRepo, ctx)
		err := eventRepo.Create(ctx, user.ID, eventData.Name, eventData.Description, eventData.Address, eventData.StartTime, eventData.EndTime)
		assert.NoError(t, err)

		events, err := eventRepo.List(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, events)
	})

	t.Run("Get", func(t *testing.T) {
		setupTest(t)

		user := setupUser(t, userRepo, ctx)
		err := eventRepo.Create(ctx, user.ID, eventData.Name, eventData.Description, eventData.Address, eventData.StartTime, eventData.EndTime)
		assert.NoError(t, err)

		events, err := eventRepo.List(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, events)

		e := events[0]
		event, err := eventRepo.Get(ctx, e.ID)
		assert.NoError(t, err)
		assert.Equal(t, e.ID, event.ID)
		assert.Equal(t, eventData.Name, event.Name)
		assert.Equal(t, eventData.Description, event.Description)
		assert.Equal(t, eventData.Address, event.Address)
		assert.Equal(t, eventData.StartTime.Unix(), event.StartTime.Unix())
		assert.Equal(t, eventData.EndTime.Unix(), event.EndTime.Unix())
	})

	t.Run("Update", func(t *testing.T) {
		setupTest(t)

		user := setupUser(t, userRepo, ctx)
		err := eventRepo.Create(ctx, user.ID, eventData.Name, eventData.Description, eventData.Address, eventData.StartTime, eventData.EndTime)
		assert.NoError(t, err)

		events, err := eventRepo.List(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, events)

		e := events[0]
		newData := struct {
			name        string
			description string
			address     string
			startTime   time.Time
			endTime     time.Time
		}{
			name:        "Updated Event",
			description: "Updated Description",
			address:     "Updated Address",
			startTime:   time.Now().UTC().Add(time.Hour),
			endTime:     time.Now().UTC().Add(2 * time.Hour),
		}

		err = eventRepo.Update(ctx, e.ID, newData.name, newData.description, newData.address, newData.startTime, newData.endTime)
		assert.NoError(t, err)

		updatedEvent, err := eventRepo.Get(ctx, e.ID)
		assert.NoError(t, err)
		assert.Equal(t, newData.name, updatedEvent.Name)
		assert.Equal(t, newData.description, updatedEvent.Description)
		assert.Equal(t, newData.address, updatedEvent.Address)
		assert.Equal(t, newData.startTime.Unix(), updatedEvent.StartTime.Unix())
		assert.Equal(t, newData.endTime.Unix(), updatedEvent.EndTime.Unix())
	})

	t.Run("Delete", func(t *testing.T) {
		setupTest(t)

		user := setupUser(t, userRepo, ctx)
		err := eventRepo.Create(ctx, user.ID, eventData.Name, eventData.Description, eventData.Address, eventData.StartTime, eventData.EndTime)
		assert.NoError(t, err)

		events, err := eventRepo.List(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, events)

		e := events[0]
		err = eventRepo.Delete(ctx, e.ID)
		assert.NoError(t, err)

		_, err = eventRepo.Get(ctx, e.ID)
		assert.Error(t, err)
	})
}
