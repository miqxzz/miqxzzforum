package usecase

import (
	"context"
	"time"

	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"github.com/Engls/forum-project2/forum_service/internal/repository"
	"go.uber.org/zap"
)

type ChatUsecase interface {
	HandleMessage(ctx context.Context, userID int, username, content string) error
	GetRecentMessages(ctx context.Context, limit int) ([]entity.ChatMessage, error)
}

type chatUsecase struct {
	repo   repository.ChatRepository
	logger *zap.Logger
}

func NewChatUsecase(repo repository.ChatRepository, logger *zap.Logger) ChatUsecase {
	return &chatUsecase{repo: repo, logger: logger}
}

func (uc *chatUsecase) HandleMessage(ctx context.Context, userID int, username, content string) error {
	message := entity.ChatMessage{
		UserID:    userID,
		Username:  username,
		Content:   content,
		Timestamp: time.Now(),
	}

	uc.logger.Info("Handling message",
		zap.Int("userID", userID),
		zap.String("username", username),
		zap.String("content", content),
	)

	err := uc.repo.StoreMessage(ctx, message)
	if err != nil {
		uc.logger.Error("Failed to store message", zap.Error(err))
		return err
	}

	uc.logger.Info("Message stored successfully", zap.Int("userID", userID), zap.String("username", username))
	return nil
}

func (uc *chatUsecase) GetRecentMessages(ctx context.Context, limit int) ([]entity.ChatMessage, error) {
	uc.logger.Info("Fetching recent messages", zap.Int("limit", limit))

	messages, err := uc.repo.GetRecentMessages(ctx, limit)
	if err != nil {
		uc.logger.Error("Failed to get recent messages", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("Recent messages fetched successfully", zap.Int("count", len(messages)))
	return messages, nil
}
