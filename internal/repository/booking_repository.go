package repository

import (
	"context"

	"github.com/amedoeyes/hadath/internal/database"
	"github.com/amedoeyes/hadath/internal/model"
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

func (r *BookingRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]model.Booking, error) {
	query := `
	SELECT id, user_id, event_id
	FROM bookings
	WHERE user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	var bookings []model.Booking
	for rows.Next() {
		var booking model.Booking
		err := rows.Scan(&booking.ID, &booking.UserID, &booking.EventID)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *BookingRepository) GetAllByEventID(ctx context.Context, eventID uuid.UUID) ([]model.Booking, error) {
	query := `
	SELECT id, user_id, event_id
	FROM bookings
	WHERE event_id = $1
	`

	rows, err := r.db.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}

	var bookings []model.Booking
	for rows.Next() {
		var booking model.Booking
		err := rows.Scan(&booking.ID, &booking.UserID, &booking.EventID)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *BookingRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE 
	FROM bookings
	WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
