package storage

import (
	"agent-care-tg/models"
	"database/sql"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GenerateDefaultTasks(userID int64) error {
	tx, err := s.db.Beginx()
	if err != nil {
		slog.Error("[GenerateDefaultTasks]: Failed to begin transaction", "error", err)
	}

	defer tx.Rollback()

	query := `INSERT INTO tasks (chat_id, name, description, isActive, isDefault, currentStreak, maxStreak) VALUES (:chat_id, :name, :description, :isActive, :isDefault, :currentStreak, :maxStreak)`

	defaultTaks := []models.Task{
		{ChatID: userID, Name: "Morning Routine", Description: "Morning Routine", IsActive: true, IsDefault: true, CurrentStreak: 0, MaxStreak: 0},
		{ChatID: userID, Name: "Sunlight", Description: "Sunlight Routine", IsActive: true, IsDefault: true, CurrentStreak: 0, MaxStreak: 0},
		{ChatID: userID, Name: "Workout", Description: "Workout", IsActive: true, IsDefault: true, CurrentStreak: 0, MaxStreak: 0},
		{ChatID: userID, Name: "Healthy Meal", Description: "Healthy Meal", IsActive: true, IsDefault: true, CurrentStreak: 0, MaxStreak: 0},
		{ChatID: userID, Name: "Personal Goal", Description: "Personal Goal", IsActive: true, IsDefault: true, CurrentStreak: 0, MaxStreak: 0},
	}

	for _, task := range defaultTaks {
		_, err := tx.NamedExec(query, task)
		if err != nil {
			slog.Error("[GenerateDefaultTasks]: Failed to execute query", "error", err)
		}
	}
	return tx.Commit()
}

func (s *Store) SaveUser(user *models.User) error {
	query := `
		INSERT INTO users (chat_id, tg_username, username, personal_goal, timezone, tasks)
	VALUES (:chat_id, :tg_username, :username, :personal_goal, :timezone, :tasks)
	`

	_, err := s.db.NamedExec(query, user)
	if err != nil {
		slog.Error("[SaveUser]: Failed to save user", "error", err)
	}

	error := s.GenerateDefaultTasks(user.ChatID)
	if error != nil {
		slog.Error("[SaveUser]: Failed to generate default tasks", "error", err)
	}

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
