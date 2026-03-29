package storage

import (
	"agent-care-tg/models"
	"github.com/jmoiron/sqlx"
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
