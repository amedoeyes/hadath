package repository

import (
	"context"
	"testing"

	"github.com/amedoeyes/hadath/internal/database/model"
	"github.com/stretchr/testify/assert"
)

var userData = model.User{
	Name:     "Cat",
	Email:    "cat@meow.purr",
	Password: "meowmeowmeow",
}

func setupUser(t *testing.T, repo *UserRepository, ctx context.Context) *model.User {
	t.Helper()
	err := repo.Create(ctx, userData.Name, userData.Email, userData.Password)
	assert.NoError(t, err)
	user, err := repo.GetByEmail(ctx, userData.Email)
	assert.NoError(t, err)
	return user
}

func TestUserRepository(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		setupTest(t)

		err := repo.Create(ctx, userData.Name, userData.Email, userData.Password)
		assert.NoError(t, err)
	})

	t.Run("GetByEmail", func(t *testing.T) {
		setupTest(t)

		err := repo.Create(ctx, userData.Name, userData.Email, userData.Password)
		assert.NoError(t, err)

		user, err := repo.GetByEmail(ctx, userData.Email)
		assert.NoError(t, err)
		assert.Equal(t, userData.Name, user.Name)
		assert.Equal(t, userData.Email, user.Email)
		assert.Equal(t, userData.Password, user.Password)
	})

	t.Run("Get", func(t *testing.T) {
		setupTest(t)

		u := setupUser(t, repo, ctx)

		user, err := repo.Get(ctx, u.ID)
		assert.NoError(t, err)
		assert.Equal(t, u, user)
	})

	t.Run("Update", func(t *testing.T) {
		setupTest(t)

		u := setupUser(t, repo, ctx)

		newData := model.User{
			Name:     "Cool Cat",
			Email:    "coolcat@meow.purr",
			Password: "catnip4lyfe",
		}

		err := repo.Update(ctx, u.ID, newData.Name, newData.Email, newData.Password)
		assert.NoError(t, err)

		updatedUser, err := repo.Get(ctx, u.ID)
		assert.NoError(t, err)
		assert.Equal(t, newData.Name, updatedUser.Name)
		assert.Equal(t, newData.Email, updatedUser.Email)
		assert.Equal(t, newData.Password, updatedUser.Password)
	})

	t.Run("Delete", func(t *testing.T) {
		setupTest(t)

		u := setupUser(t, repo, ctx)

		err := repo.Delete(ctx, u.ID)
		assert.NoError(t, err)

		_, err = repo.Get(ctx, u.ID)
		assert.Error(t, err)
	})
}
