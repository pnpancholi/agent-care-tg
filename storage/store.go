package storage

import (
	"agent-care-tg/models"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"time"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) SaveUser(user *models.User) error {
	query := `
		INSERT INTO users (chat_id, tg_username, username, personal_goal, timezone, tasks)
	VALUES (:chat_id, :tg_username, :username, :personal_goal, :timezone, :tasks)
	`

	_, err := s.db.NamedExec(query, user)
	return err
}

func (s *Store) GetAllUsers() ([]models.User, error) {
	var users []models.User
	query := `
		SELECT * FROM users
	`
	err := s.db.Select(&users, query)

	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Store) UpdateLastSentAt(user *models.User) error {
	lastSeenAt := sql.NullTime{Time: time.Now().UTC(), Valid: true}
	chatId := user.ChatID
	query := `UPDATE users SET last_sent_at = $1 WHERE chat_id = $2`

	_, err := s.db.Exec(query, lastSeenAt, chatId)
	return err

}
