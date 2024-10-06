package repository

import (
	"context"

	"github.com/amedoeyes/hadath/internal/database"
	"github.com/amedoeyes/hadath/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository() *UserRepository {
	return &UserRepository{database.Get()}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, created_at, updated_at"
	err := r.db.QueryRow(ctx, query, user.Email, user.Password).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint32) (*model.User, error) {
	query := "SELECT id, email, password, created_at, updated_at FROM users WHERE id = $1"
	user := &model.User{}
	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) List(ctx context.Context) ([]model.User, error) {
	query := "SELECT id, email, password, created_at, updated_at FROM users"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := "UPDATE users SET email = $1, password = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING updated_at"
	err := r.db.QueryRow(
		ctx,
		query,
		user.Email,
		user.Password,
		user.ID,
	).Scan(&user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uint32) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
