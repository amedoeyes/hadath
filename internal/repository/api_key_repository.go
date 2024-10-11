package repository

import (
	"context"

	"github.com/amedoeyes/hadath/internal/database"
	"github.com/amedoeyes/hadath/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIKeyRepository struct {
	db *pgxpool.Pool
}

func NewAPIKeyRepository() *APIKeyRepository {
	return &APIKeyRepository{database.Get()}
}

func (r *APIKeyRepository) Create(ctx context.Context, userID uuid.UUID, key string) error {
	query := `
	INSERT
	INTO api_keys (user_id, key)
	VALUES ($1, $2)
	`

	_, err := r.db.Exec(ctx, query, userID, key)
	if err != nil {
		return err
	}

	return nil
}

func (r *APIKeyRepository) GetByKey(ctx context.Context, key string) (*model.APIKey, error) {
	query := `
	SELECT id, user_id, key
	FROM api_keys
	WHERE key = $1
	`

	apiKey := &model.APIKey{}
	err := r.db.QueryRow(ctx, query, key).Scan(&apiKey.ID, &apiKey.UserID, &apiKey.Key)
	if err != nil {
		return nil, err
	}

	return apiKey, nil
}

func (r *APIKeyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE
	FROM api_keys
	WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
