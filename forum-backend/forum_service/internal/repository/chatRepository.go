package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"go.uber.org/zap"
)

type DB interface {
	Exec(query string, args ...any) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type ChatRepository interface {
	StoreMessage(ctx context.Context, msg entity.ChatMessage) error
	GetRecentMessages(ctx context.Context, limit int) ([]entity.ChatMessage, error)
}

type chatRepo struct {
	db     DB
	logger *zap.Logger
}

func NewChatRepository(db DB, logger *zap.Logger) ChatRepository {
	return &chatRepo{db: db, logger: logger}
}

func (r *chatRepo) StoreMessage(ctx context.Context, msg entity.ChatMessage) error {
	r.logger.Info("Saving message",
		zap.Int("userID", msg.UserID),
		zap.String("username", msg.Username),
		zap.String("content", msg.Content),
		zap.Time("timestamp", msg.Timestamp),
	)

	query := `INSERT INTO chat_messages (user_id, username, content, timestamp) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, msg.UserID, msg.Username, msg.Content, msg.Timestamp.Format(time.RFC3339))
	if err != nil {
		r.logger.Error("Failed to store message", zap.Error(err))
		return err
	}
	return nil
}

func (r *chatRepo) GetRecentMessages(ctx context.Context, limit int) ([]entity.ChatMessage, error) {
	query := `
        SELECT id, user_id, username, content,
               timestamp
        FROM chat_messages
        ORDER BY timestamp DESC
        LIMIT ?`

	var messages []entity.ChatMessage
	err := r.db.SelectContext(ctx, &messages, query, limit)
	if err != nil {
		r.logger.Error("Failed to get recent messages", zap.Error(err))
		return nil, err
	}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	r.logger.Info("Recent messages retrieved successfully", zap.Int("count", len(messages)))
	return messages, nil
}
