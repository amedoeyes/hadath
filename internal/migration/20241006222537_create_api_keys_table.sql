-- +goose Up
-- +goose StatementBegin
CREATE TABLE api_keys (
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL REFERENCES users(id),
	key VARCHAR(255) NOT NULL,
	FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS api_keys
-- +goose StatementEnd
