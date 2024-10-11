-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE bookings (
	user_id UUID REFERENCES users(id) ON DELETE SET NULL,
	event_id UUID REFERENCES events(id) ON DELETE SET NULL,
	UNIQUE (user_id, event_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bookings;
-- +goose StatementEnd
