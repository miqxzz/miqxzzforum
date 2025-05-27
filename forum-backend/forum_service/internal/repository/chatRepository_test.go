package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"github.com/Engls/forum-project2/forum_service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestChatRepo_StoreMessage_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockDB := new(mocks.DB)

	chatRepo := NewChatRepository(mockDB, logger)

	msg := entity.ChatMessage{
		UserID:    1,
		Username:  "testuser",
		Content:   "This is a test message",
		Timestamp: time.Date(2025, time.April, 22, 23, 51, 38, 843016900, time.Local),
	}

	mockDB.On("ExecContext", mock.Anything, mock.Anything, msg.UserID, msg.Username, msg.Content, msg.Timestamp.Format(time.RFC3339)).Return(sql.Result(nil), nil)

	err := chatRepo.StoreMessage(context.Background(), msg)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

func TestChatRepo_StoreMessage_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockDB := new(mocks.DB)

	chatRepo := NewChatRepository(mockDB, logger)

	msg := entity.ChatMessage{
		UserID:    1,
		Username:  "testuser",
		Content:   "This is a test message",
		Timestamp: time.Now(),
	}

	mockDB.On("ExecContext", mock.Anything, mock.Anything, msg.UserID, msg.Username, msg.Content, msg.Timestamp.Format(time.RFC3339)).Return(nil, errors.New("failed to store message"))

	err := chatRepo.StoreMessage(context.Background(), msg)

	assert.Error(t, err)

	mockDB.AssertExpectations(t)
}

func TestChatRepo_GetRecentMessages_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockDB := new(mocks.DB)

	chatRepo := NewChatRepository(mockDB, logger)

	limit := 10
	messages := []entity.ChatMessage{
		{ID: 1, UserID: 1, Username: "user1", Content: "Message 1", Timestamp: time.Now()},
		{ID: 2, UserID: 2, Username: "user2", Content: "Message 2", Timestamp: time.Now()},
	}

	mockDB.On("SelectContext", mock.Anything, mock.Anything, mock.Anything, limit).Return(nil).Run(func(args mock.Arguments) {
		dest := args.Get(1).(*[]entity.ChatMessage)
		*dest = messages
	})

	result, err := chatRepo.GetRecentMessages(context.Background(), limit)

	assert.NoError(t, err)
	assert.Equal(t, messages, result)

	mockDB.AssertExpectations(t)
}

func TestChatRepo_GetRecentMessages_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockDB := new(mocks.DB)

	chatRepo := NewChatRepository(mockDB, logger)

	limit := 10

	mockDB.On("SelectContext", mock.Anything, mock.Anything, mock.Anything, limit).Return(errors.New("failed to get recent messages"))

	result, err := chatRepo.GetRecentMessages(context.Background(), limit)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockDB.AssertExpectations(t)
}
