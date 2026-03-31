package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type Tasks []Task

type User struct {
	ChatID       int64        `json:"chat_id" db:"chat_id"`
	Username     string       `json:"username" db:"username"`
	TGUsername   string       `json:"tg_username" db:"tg_username"`
	Timezone     string       `json:"timezone" db:"timezone"`
	PersonalGoal string       `json:"personal_goal" db:"personal_goal"`
	LastSentAt   sql.NullTime `json:"last_sent_at" db:"last_sent_at"`
	Tasks        Tasks        `json:"tasks" db:"tasks"`
}

func (t *Tasks) Scan(value any) error {
	if value == nil {
		*t = Tasks{}
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, t)
	case string:
		return json.Unmarshal([]byte(v), t)
	default:
		return fmt.Errorf("Unsupported type: %T", v)
	}
}

func (t Tasks) Value() (driver.Value, error) {
	if t == nil {
		return "[]", nil
	}
	return json.Marshal(t)
}

func NewUser() *User {
	return &User{
		Tasks: Tasks{},
	}
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
