-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
	id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(id) ON DELETE SET NULL,
	name VARCHAR(255) NOT NULL,
	description TEXT,
	address VARCHAR(255) NOT NULL,
	start_time TIMESTAMP NOT NULL,
	end_time TIMESTAMP NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
