package models

import "database/sql"

type User struct {
	ChatID       int64        `json:"chat_id" db:"chat_id"`
	Username     string       `json:"username" db:"username"`
	TGUsername   string       `json:"tg_username" db:"tg_username"`
	Timezone     string       `json:"timezone" db:"timezone"`
	PersonalGoal string       `json:"personal_goal" db:"personal_goal"`
	LastSentAt   sql.NullTime `json:"last_sent_at" db:"last_sent_at"`
}

type Task struct {
	Name      string `json:"name"`
	Performed bool   `json:"performed"`
	Streak    int64  `json:"streak"`
	Time      string `json:"time"`
	Expiry    string `json:"expiry,omitempty"`
}
