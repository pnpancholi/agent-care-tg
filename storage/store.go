package storage

import (
	"agent-care-tg/models"
	"database/sql"
	"fmt"
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
func (s *Store) IncrementStreak(chatID int64, taskTag string) error {
	query := `UPDATE tasks SET current_streak = current_streak + 1 WHERE chat_id = $1 AND tag = $2`
	_, err := s.db.Exec(query, chatID, taskTag)

	if err != nil {
		slog.Error("Failed to update user's streak", "err", err)
		return fmt.Errorf("Failed to update user's streak: %w", err)
	}
	return nil
}

func (s *Store) GenerateDefaultTasks(userID int64) error {
	tx, err := s.db.Beginx()
	if err != nil {
		slog.Error("[GenerateDefaultTasks]: Failed to begin transaction", "error", err)
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	query := `INSERT INTO tasks (chat_id, name, description, tag, is_active, is_default, current_streak, max_streak) VALUES (:chat_id, :name, :description, :tag, :is_active, :is_default, :current_streak, :max_streak)`

	defaultTasks := []models.Task{
		{ChatID: userID, Name: "Morning Routine", Description: "Morning Routine", Tag: "daily_morning", IsActive: true, IsDefault: true, CurrentStreak: 0, MaxStreak: 0},
		{ChatID: userID, Name: "Sunlight", Description: "Sunlight Routine", Tag: "daily_sunlight", IsActive: true, IsDefault: true, CurrentStreak: 0, MaxStreak: 0},
		{ChatID: userID, Name: "Workout", Description: "Workout", Tag: "daily_excercise", IsActive: true, IsDefault: true, CurrentStreak: 0, MaxStreak: 0},
		{ChatID: userID, Name: "Healthy Meal", Description: "Healthy Meal", Tag: "daily_meal", IsActive: true, IsDefault: true, CurrentStreak: 0, MaxStreak: 0},
		{ChatID: userID, Name: "Personal Goal", Description: "Personal Goal", Tag: "daily_personal", IsActive: true, IsDefault: true, CurrentStreak: 0, MaxStreak: 0},
	}

	for _, task := range defaultTasks {
		_, err := tx.NamedExec(query, task)
		if err != nil {
			slog.Error("[GenerateDefaultTasks]: Failed to execute query", "error", err)
			return fmt.Errorf("Failed to execute query: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		slog.Error("[GenerateDefaultTasks]: Failed to commit transaction", "error", err)
		return fmt.Errorf("Failed to commit transaction: %w", err)
	}
	return nil
}

func (s *Store) SaveUser(user *models.User) error {
	query := `
		INSERT INTO users (chat_id, tg_username, username, personal_goal, timezone)
	VALUES (:chat_id, :tg_username, :username, :personal_goal, :timezone )
	`

	_, err := s.db.NamedExec(query, user)
	if err != nil {
		slog.Error("[SaveUser]: Failed to save user", "error", err)
		return fmt.Errorf("Failed to save user: %w", err)
	}

	err = s.GenerateDefaultTasks(user.ChatID)
	if err != nil {
		slog.Error("[SaveUser]: Failed to generate default tasks", "error", err)
		return fmt.Errorf("Failed to generate default tasks: %w", err)
	}

	return nil
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
