package repository

import (
	"context"
	"testing"

	"github.com/amedoeyes/hadath/internal/database/model"
	"github.com/stretchr/testify/assert"
)

var apiKeyData = model.APIKey{
	Key: "secret_key",
}

func TestAPIKeyRepository(t *testing.T) {
	apiKeyRepo := NewAPIKeyRepository()
	userRepo := NewUserRepository()
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		setupTest(t)

		user := setupUser(t, userRepo, ctx)

		err := apiKeyRepo.Create(ctx, user.ID, apiKeyData.Key)
		assert.NoError(t, err)
	})

	t.Run("GetByKey", func(t *testing.T) {
		setupTest(t)

		user := setupUser(t, userRepo, ctx)

		err := apiKeyRepo.Create(ctx, user.ID, apiKeyData.Key)
		assert.NoError(t, err)

		apiKey, err := apiKeyRepo.GetByKey(ctx, apiKeyData.Key)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, apiKey.UserID)
	})

	t.Run("Delete", func(t *testing.T) {
		setupTest(t)

		user := setupUser(t, userRepo, ctx)

		err := apiKeyRepo.Create(ctx, user.ID, apiKeyData.Key)
		assert.NoError(t, err)

		apiKey, err := apiKeyRepo.GetByKey(ctx, apiKeyData.Key)
		assert.NoError(t, err)

		err = apiKeyRepo.Delete(ctx, apiKey.ID)
		assert.NoError(t, err)

		_, err = apiKeyRepo.GetByKey(ctx, apiKeyData.Key)
		assert.Error(t, err)
	})
}
