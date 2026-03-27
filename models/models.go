package models

import "database/sql"

type User struct {
	ChatID       int64        `json:"chat_id" db:"chat_id"`
	Username     string       `json:"username" db:"username"`
	TGUsername   string       `json:"tg_username" db:"tg_username"`
	Timezone     string       `json:"timezone" db:"timezone"`
	PersonalGoal string       `json:"personal_goal" db:"personal_goal"`
	LastSentAt   sql.NullTime `json:"last_sent_at" db:"last_sent_at"`
	Tasks        []Task       `json:"tasks" db:"tasks"`
}

type Task struct {
	Name            string `json:"name"`
	ScheduledHour   uint8  `json:"scheduled_hour"`
	IsActive        bool   `json:"is_active"`
	isDefault       bool   `json:"is_default"`
	CurrentStreak   uint64 `json:"current_streak"`
	MaxStreak       uint64 `json:"max_streak"`
	LastCompletedAt string `json:"last_completed_at"`
}
