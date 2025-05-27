package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/miqxzz/miqxzzforum/forum_service/internal/entity"
	"github.com/miqxzz/miqxzzforum/forum_service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestCommentsUsecase_CreateComment_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(mocks.CommentsRepository)
	usecase := NewCommentsUsecases(mockRepo, logger)

	comment := entity.Comment{
		AuthorId: 1,
		PostId:   1,
		Content:  "Test comment",
	}

	expectedComment := entity.Comment{
		ID:        1,
		AuthorId:  1,
		PostId:    1,
		Content:   "Test comment",
		CreatedAt: time.Now(),
	}

	mockRepo.On("CreateComment", mock.Anything, comment).Return(expectedComment, nil)

	createdComment, err := usecase.CreateComment(context.Background(), comment)
	assert.NoError(t, err)
	assert.Equal(t, expectedComment, createdComment)
	mockRepo.AssertExpectations(t)
}

func TestCommentsUsecase_CreateComment_Failure(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(mocks.CommentsRepository)
	usecase := NewCommentsUsecases(mockRepo, logger)

	comment := entity.Comment{
		AuthorId: 1,
		PostId:   1,
		Content:  "Test comment",
	}

	mockRepo.On("CreateComment", mock.Anything, comment).Return(entity.Comment{}, assert.AnError)

	createdComment, err := usecase.CreateComment(context.Background(), comment)
	assert.Error(t, err)
	assert.Equal(t, entity.Comment{}, createdComment)
	mockRepo.AssertExpectations(t)
}

func TestCommentsUsecase_GetComments_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(mocks.CommentsRepository)
	usecase := NewCommentsUsecases(mockRepo, logger)

	expectedComments := []entity.Comment{
		{
			ID:        1,
			AuthorId:  1,
			PostId:    1,
			Content:   "Comment 1",
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			AuthorId:  2,
			PostId:    1,
			Content:   "Comment 2",
			CreatedAt: time.Now(),
		},
	}

	mockRepo.On("GetComments", mock.Anything, 1, 10, 0).Return(expectedComments, nil)

	comments, err := usecase.GetComments(context.Background(), 1, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, expectedComments, comments)
	mockRepo.AssertExpectations(t)
}

func TestCommentsUsecase_GetComments_Failure(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(mocks.CommentsRepository)
	usecase := NewCommentsUsecases(mockRepo, logger)

	mockRepo.On("GetComments", mock.Anything, 1, 10, 0).Return(nil, assert.AnError)

	comments, err := usecase.GetComments(context.Background(), 1, 10, 0)
	assert.Error(t, err)
	assert.Nil(t, comments)
	mockRepo.AssertExpectations(t)
}

func TestCommentsUsecase_GetTotalCommentsCount_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(mocks.CommentsRepository)
	usecase := NewCommentsUsecases(mockRepo, logger)

	expectedCount := 5

	mockRepo.On("GetTotalCommentsCount", mock.Anything, 1).Return(expectedCount, nil)

	count, err := usecase.GetTotalCommentsCount(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	mockRepo.AssertExpectations(t)
}

func TestCommentsUsecase_GetTotalCommentsCount_Failure(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(mocks.CommentsRepository)
	usecase := NewCommentsUsecases(mockRepo, logger)

	mockRepo.On("GetTotalCommentsCount", mock.Anything, 1).Return(0, assert.AnError)

	count, err := usecase.GetTotalCommentsCount(context.Background(), 1)
	assert.Error(t, err)
	assert.Equal(t, 0, count)
	mockRepo.AssertExpectations(t)
}
