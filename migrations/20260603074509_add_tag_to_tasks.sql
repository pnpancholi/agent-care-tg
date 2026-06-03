-- +goose Up
ALTER TABLE tasks ADD COLUMN tag text;

-- +goose Down
AlTER TABLE tasks DROP COLUMN tag;
