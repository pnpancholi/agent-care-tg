package models

type User struct {
	ChatID        int64  `json:"chat_id"`
	PrefferedName string `json:"preffered_name"`
	Timezone      string `json:"timezone"`
}

type Task struct {
	Name      string `json:"name"`
	Performed bool   `json:"performed"`
	Streak    int64  `json:"streak"`
	Time      string `json:"time"`
	Expiry    string `json:"expiry,omitempty"`
}
