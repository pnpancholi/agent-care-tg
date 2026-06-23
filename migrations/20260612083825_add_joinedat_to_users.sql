-- +goose Up
ALTER TABLE users ADD COLUMN joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

-- +goose Down
AlTER TABLE users DROP COLUMN joined_at;
