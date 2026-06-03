-- +goose Up
CREATE TABLE tasks (
  id SERIAL PRIMARY KEY,
  chat_id BIGINT NOT NULL REFERENCES users(chat_id),
  name VARCHAR NOT NULL,
  description VARCHAR NOT NULL,
  max_streak INT NOT NULL DEFAULT 0,
  current_streak INT NOT NULL DEFAULT 0,
  is_default BOOLEAN NOT NULL DEFAULT TRUE,
  is_active BOOLEAN NOT NULL DEFAULT TRUE
);
-- +goose Down
DROP TABLE tasks;

