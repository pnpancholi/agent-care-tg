-- +goose Up
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_sent_at TIMESTAMP;

-- +goose Down
ALTER TABLE users DROP COLUMN IF NOT EXISTS last_sent_at;
