package models

import (
	"database/sql"
)

type User struct {
	ChatID       int64        `json:"chat_id" db:"chat_id"`
	Username     string       `json:"username" db:"username"`
	TGUsername   string       `json:"tg_username" db:"tg_username"`
	Timezone     string       `json:"timezone" db:"timezone"`
	PersonalGoal string       `json:"personal_goal" db:"personal_goal"`
	LastSentAt   sql.NullTime `json:"last_sent_at" db:"last_sent_at"`
}
type Task struct {
	ID int64 `json:"id" db:"id"`
	// this is the chat_id of the user, to maintain one to many relationship
	ChatID        int64  `json:"chat_id" db:"chat_id"`
	Name          string `json:"name" db:"name"`
	Description   string `json:"description" db:"description"`
	Tag           string `json:"tag" db:"tag"`
	IsActive      bool   `json:"is_active" db:"is_active"`
	IsDefault     bool   `json:"is_default" db:"is_default"`
	CurrentStreak uint64 `json:"current_streak" db:"current_streak"`
	MaxStreak     uint64 `json:"max_streak" db:"max_streak"`
}

func NewUser() *User {
	return &User{}
}
