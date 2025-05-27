package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"github.com/Engls/forum-project2/forum_service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestCommentsUsecases_CreateComment_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockCommentRepo := new(mocks.CommentsRepository)

	commentsUsecases := NewCommentsUsecases(mockCommentRepo, logger)

	comment := entity.Comment{
		PostId:   1,
		AuthorId: 1,
		Content:  "This is a test comment",
	}
	createdComment := comment
	createdComment.ID = 1

	mockCommentRepo.On("CreateComment", mock.Anything, comment).Return(createdComment, nil)

	result, err := commentsUsecases.CreateComment(context.Background(), comment)

	assert.NoError(t, err)
	assert.Equal(t, createdComment, result)

	mockCommentRepo.AssertExpectations(t)
}

func TestCommentsUsecases_CreateComment_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockCommentRepo := new(mocks.CommentsRepository)

	commentsUsecases := NewCommentsUsecases(mockCommentRepo, logger)

	comment := entity.Comment{
		PostId:   1,
		AuthorId: 1,
		Content:  "This is a test comment",
	}

	mockCommentRepo.On("CreateComment", mock.Anything, comment).Return(entity.Comment{}, errors.New("failed to create comment"))

	result, err := commentsUsecases.CreateComment(context.Background(), comment)

	assert.Error(t, err)
	assert.Equal(t, entity.Comment{}, result)

	mockCommentRepo.AssertExpectations(t)
}

func TestCommentsUsecases_GetCommentByPostID_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockCommentRepo := new(mocks.CommentsRepository)

	commentsUsecases := NewCommentsUsecases(mockCommentRepo, logger)

	comments := []entity.Comment{
		{ID: 1, PostId: 1, AuthorId: 1, Content: "Comment 1"},
		{ID: 2, PostId: 1, AuthorId: 2, Content: "Comment 2"},
	}

	mockCommentRepo.On("GetCommentsByPostID", mock.Anything, 1).Return(comments, nil)

	result, err := commentsUsecases.GetCommentByPostID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, comments, result)

	mockCommentRepo.AssertExpectations(t)
}

func TestCommentsUsecases_GetCommentByPostID_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockCommentRepo := new(mocks.CommentsRepository)

	commentsUsecases := NewCommentsUsecases(mockCommentRepo, logger)

	mockCommentRepo.On("GetCommentsByPostID", mock.Anything, 1).Return(nil, errors.New("failed to get comments"))

	result, err := commentsUsecases.GetCommentByPostID(context.Background(), 1)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockCommentRepo.AssertExpectations(t)
}
