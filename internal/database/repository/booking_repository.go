package repository

import (
	"context"

	"github.com/amedoeyes/hadath/internal/database"
	"github.com/amedoeyes/hadath/internal/database/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository struct {
	db *pgxpool.Pool
}

func NewBookingRepository() *BookingRepository {
	return &BookingRepository{database.Get()}
}

func (r *BookingRepository) Create(ctx context.Context, userID, eventID uuid.UUID) error {
	query := `
	INSERT 
	INTO bookings (user_id, event_id) 
	VALUES ($1, $2)
	`

	_, err := r.db.Exec(ctx, query, userID, eventID)
	if err != nil {
		return err
	}

	return nil
}

func (r *BookingRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Event, error) {
	query := `
	SELECT
		event_id,
		events.name,
		events.description,
		events.address,
		events.start_time,
		events.end_time,
		users.id AS user_id,
		users.name AS user_name
	FROM bookings
	JOIN events ON events.id = bookings.event_id
	JOIN users ON users.id = events.user_id
	WHERE bookings.user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

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

func (r *BookingRepository) ListByEvent(ctx context.Context, eventID uuid.UUID) ([]model.User, error) {
	query := `
	SELECT
		user_id,
		users.name AS user_name
	FROM bookings
	JOIN users ON users.id = bookings.user_id
	WHERE bookings.event_id = $1
	`

	rows, err := r.db.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *BookingRepository) Delete(ctx context.Context, userID, eventID uuid.UUID) error {
	query := `
	DELETE 
	FROM bookings
	WHERE user_id = $1 AND event_id = $2
	`

	_, err := r.db.Exec(ctx, query, userID, eventID)
	if err != nil {
		return err
	}

	return nil
}
