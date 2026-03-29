-- +goose Up
ALTER TABLE users ADD COLUMN tasks JSONB DEFAULT '[]'

-- +goose Down
ALTER TABLE users DROP COLUMN tasks;
