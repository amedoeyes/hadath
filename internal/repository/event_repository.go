package repository

import (
	"context"
	"time"

	"github.com/amedoeyes/hadath/internal/database"
	"github.com/amedoeyes/hadath/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository() *EventRepository {
	return &EventRepository{database.Get()}
}

func (r *EventRepository) Create(ctx context.Context, user_id uuid.UUID, name, description, address string, startTime, endTime time.Time) error {
	query := "INSERT INTO events (user_id, name, description, address, start_time, end_time) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := r.db.Exec(ctx, query, user_id, name, description, address, startTime, endTime)
	if err != nil {
		return err
	}

	return nil
}

func (r *EventRepository) GetAll(ctx context.Context) ([]model.Event, error) {
	query := "SELECT * FROM events"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	events, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Event])
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *EventRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	query := "SELECT * FROM events WHERE id = $1"

	rows, err := r.db.Query(ctx, query, id)
	event, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[model.Event])
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (r *EventRepository) UpdateByID(ctx context.Context, id uuid.UUID, name, description, address string, startTime, endTime time.Time) error {
	query := "UPDATE events SET name = $1, description = $2, address = $3, start_time = $4, end_time = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $6"

	_, err := r.db.Exec(ctx, query, name, description, address, startTime, endTime, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *EventRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM events WHERE id = $1"

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
