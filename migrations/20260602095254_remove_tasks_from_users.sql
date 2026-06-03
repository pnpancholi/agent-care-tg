-- +goose Up
ALTER TABLE users DROP COLUMN tasks;

-- +goose Down
ALTER TABLE users ADD COLUMN tasks JSONB;
