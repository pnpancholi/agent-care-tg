-- +goose Up
CREATE TABLE IF NOT EXISTS users (
  chat_id BIGINT PRIMARY KEY,
  tg_username VARCHAR(100),
  username VARCHAR(100),
  personal_goal TEXT,
  timezone VARCHAR(50)
);
-- +goose Down
DROP TABLE IF EXISTS users;
