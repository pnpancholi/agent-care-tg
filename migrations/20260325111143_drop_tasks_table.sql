-- +goose Up
DROP TABLE tasks;

-- +goose Down
CREATE TABLE tasks (
  id SERIAL PRIMARY KEY,
  chat_id BIGINT NOT NULL REFERENCES users(chat_id),
  name VARCHAR NOT NULL,
  description VARCHAR NOT NULL,
  scheduled_hour VARCHAR NOT NULL,
  is_default BOOLEAN NOT NULL,
  is_active BOOLEAN NOT NULL
);

