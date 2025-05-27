package entity

import "time"

type ChatMessage struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"userID" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	Content   string    `json:"content" db:"content"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

func (m *ChatMessage) Validate() error {
	if m.UserID <= 0 {
		return ErrInvalidUserID
	}
	if m.Username == "" {
		return ErrEmptyUsername
	}
	if m.Content == "" {
		return ErrEmptyContent
	}
	return nil
}
