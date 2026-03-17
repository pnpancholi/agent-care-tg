-- +goose Up
ALTER TABLE users ADD COLUMN last_sent_at TIMESTAMP;

-- +goose Down
ALTER TABLE users DROP COLUMN last_sent_at;
