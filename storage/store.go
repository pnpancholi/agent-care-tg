package storage

import (
	"agent-care-tg/models"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"time"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) ResetStreak(chatID int64, taskTag string) error {
	query := `UPDATE tasks SET current_streak = 0 WHERE chat_id = $1 AND tag = $2`
	_, err := s.db.Exec(query, chatID, taskTag)

	if err != nil {
		slog.Error("Failed to reset streak", "err", err)
		return fmt.Errorf("Failed to reset streak : %w", err)
	}
	return nil
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

func (s *Store) GetUserByChatID(chatID int64) (models.User, error) {
	var user models.User
	query := `SELECT * from users WHERE chat_id = $1`

	err := s.db.Get(&user, query, chatID)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *Store) GetAllTasksForUserByChatID(chatID int64) ([]models.Task, error) {
	var tasks []models.Task

	query := `SELECT * FROM tasks WHERE chat_id = $1`
	err := s.db.Select(&tasks, query, chatID)

	if err != nil {
		return tasks, err
	}
	slog.Info("tasks", "tasks", tasks)
	return tasks, nil
}

func (s *Store) GetTask(chatID int64, taskTag string) (models.Task, error) {
	var task models.Task
	query := `SELECT * FROM tasks WHERE chat_id = $1 AND tag = $2`
	err := s.db.Get(&task, query, chatID, taskTag)

	if err != nil {
		return task, err
	}
	return task, nil
}

func (s *Store) UpdateMaxStreak(taskID int64, newStreak int64) error {
	query := `UPDATE tasks SET max_streak = $1 WHERE id = $2`

	res, err := s.db.Exec(query, newStreak, taskID)

	if err != nil {
		slog.Error("Failed to update max streak", "error", err)
		return fmt.Errorf("Failed to update max streak: %w", err)
	}

	rows, _ := res.RowsAffected()

	if rows != 1 {
		slog.Error("Failed to update max streak", "error", err)
		return fmt.Errorf("Failed to update max streak: %w", err)
	}

	return nil
}
