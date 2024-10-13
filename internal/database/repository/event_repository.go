package repository

import (
	"context"
	"time"

	"github.com/amedoeyes/hadath/internal/database"
	"github.com/amedoeyes/hadath/internal/database/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository() *EventRepository {
	return &EventRepository{database.Get()}
}

func (r *EventRepository) Create(ctx context.Context, user_id uuid.UUID, name, description, address string, startTime, endTime time.Time) error {
	query := `
	INSERT
	INTO events (user_id, name, description, address, start_time, end_time) 
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query, user_id, name, description, address, startTime, endTime)
	if err != nil {
		return err
	}

	return nil
}

func (r *EventRepository) List(ctx context.Context) ([]model.Event, error) {
	query := `
	SELECT 
		events.id AS event_id,
		events.name AS event_name,
		events.description AS event_description,
		events.address AS event_address,
		events.start_time AS event_start_time,
		events.end_time AS event_end_time,
		user_id,
		users.name AS user_name
	FROM events
	JOIN users ON events.user_id = users.id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.Event
	for rows.Next() {
		var event model.Event
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Description,
			&event.Address,
			&event.StartTime,
			&event.EndTime,
			&event.User.ID,
			&event.User.Name,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) Get(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	query := `
	SELECT 
		events.id AS event_id,
		events.name AS event_name,
		events.description AS event_description,
		events.address AS event_address,
		events.start_time AS event_start_time,
		events.end_time AS event_end_time,
		user_id,
		users.name AS user_name
	FROM events
	JOIN users ON events.user_id = users.id
	WHERE events.id = $1
	`

	event := &model.Event{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&event.ID,
		&event.Name,
		&event.Description,
		&event.Address,
		&event.StartTime,
		&event.EndTime,
		&event.User.ID,
		&event.User.Name,
	)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (r *EventRepository) Update(ctx context.Context, id uuid.UUID, name, description, address string, startTime, endTime time.Time) error {
	query := `
	UPDATE events
	SET name = $1, description = $2, address = $3, start_time = $4, end_time = $5, updated_at = CURRENT_TIMESTAMP
	WHERE id = $6
	`

	_, err := r.db.Exec(ctx, query, name, description, address, startTime, endTime, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *EventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE
	FROM events
	WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
