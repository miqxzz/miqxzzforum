package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"github.com/Engls/forum-project2/forum_service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestChatUsecase_HandleMessage_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockChatRepo := new(mocks.ChatRepository)

	chatUsecase := NewChatUsecase(mockChatRepo, logger)

	userID := 1
	username := "testuser"
	content := "This is a test message"
	message := entity.ChatMessage{
		UserID:    userID,
		Username:  username,
		Content:   content,
		Timestamp: time.Now(),
	}

	mockChatRepo.On("StoreMessage", mock.Anything, message).Return(nil)

	err := chatUsecase.HandleMessage(context.Background(), userID, username, content)

	assert.NoError(t, err)

	mockChatRepo.AssertExpectations(t)
}

func TestChatUsecase_HandleMessage_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockChatRepo := new(mocks.ChatRepository)

	chatUsecase := NewChatUsecase(mockChatRepo, logger)

	userID := 1
	username := "testuser"
	content := "This is a test message"
	message := entity.ChatMessage{
		UserID:    userID,
		Username:  username,
		Content:   content,
		Timestamp: time.Now(),
	}

	mockChatRepo.On("StoreMessage", mock.Anything, message).Return(errors.New("failed to store message"))

	err := chatUsecase.HandleMessage(context.Background(), userID, username, content)

	assert.Error(t, err)

	mockChatRepo.AssertExpectations(t)
}

func TestChatUsecase_GetRecentMessages_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockChatRepo := new(mocks.ChatRepository)

	chatUsecase := NewChatUsecase(mockChatRepo, logger)

	limit := 10
	messages := []entity.ChatMessage{
		{UserID: 1, Username: "user1", Content: "Message 1", Timestamp: time.Now()},
		{UserID: 2, Username: "user2", Content: "Message 2", Timestamp: time.Now()},
	}

	mockChatRepo.On("GetRecentMessages", mock.Anything, limit).Return(messages, nil)

	result, err := chatUsecase.GetRecentMessages(context.Background(), limit)

	assert.NoError(t, err)
	assert.Equal(t, messages, result)

	mockChatRepo.AssertExpectations(t)
}

func TestChatUsecase_GetRecentMessages_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockChatRepo := new(mocks.ChatRepository)

	chatUsecase := NewChatUsecase(mockChatRepo, logger)

	limit := 10

	mockChatRepo.On("GetRecentMessages", mock.Anything, limit).Return(nil, errors.New("failed to get recent messages"))

	result, err := chatUsecase.GetRecentMessages(context.Background(), limit)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockChatRepo.AssertExpectations(t)
}
