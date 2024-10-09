package repository

import (
	"context"

	"github.com/amedoeyes/hadath/internal/database"
	"github.com/amedoeyes/hadath/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository() *UserRepository {
	return &UserRepository{database.Get()}
}

func (r *UserRepository) Create(ctx context.Context, name, email, password string) error {
	query := `
	INSERT
	INTO users (name, email, password)
	VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(ctx, query, name, email, password)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
	SELECT id, name, email, created_at, updated_at
	FROM users
	WHERE email = $1
	`

	rows, err := r.db.Query(ctx, query, email)
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	rows.Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `
	SELECT id, name, email, created_at, updated_at
	FROM users
	WHERE id = $1
	`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	rows.Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, id uuid.UUID, name, email, password string) error {
	query := `
	UPDATE users
	SET name = $1, email = $2, password = $3, updated_at = CURRENT_TIMESTAMP
	WHERE id = $4 RETURNING updated_at
	`

	_, err := r.db.Exec(ctx, query, name, email, password, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE
	FROM users
	WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
