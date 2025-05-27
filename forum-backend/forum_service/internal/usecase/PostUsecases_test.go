package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/miqxzz/miqxzzforum/forum_service/internal/entity"
	"github.com/miqxzz/miqxzzforum/forum_service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestPostUsecase_CreatePost_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(mocks.PostRepository)
	usecase := NewPostUsecase(mockRepo, logger)

	post := entity.Post{
		AuthorId: 1,
		Title:    "Test Post",
		Content:  "This is a test post",
	}

	expectedPost := entity.Post{
		ID:       1,
		AuthorId: 1,
		Title:    "Test Post",
		Content:  "This is a test post",
	}

	mockRepo.On("CreatePost", mock.Anything, post).Return(expectedPost, nil)

	createdPost, err := usecase.CreatePost(context.Background(), post)
	assert.NoError(t, err)
	assert.Equal(t, expectedPost, createdPost)
	mockRepo.AssertExpectations(t)
}

func TestPostUsecase_CreatePost_Failure(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(mocks.PostRepository)
	usecase := NewPostUsecase(mockRepo, logger)

	post := entity.Post{
		AuthorId: 1,
		Title:    "Test Post",
		Content:  "This is a test post",
	}

	mockRepo.On("CreatePost", mock.Anything, post).Return(entity.Post{}, assert.AnError)

	createdPost, err := usecase.CreatePost(context.Background(), post)
	assert.Error(t, err)
	assert.Equal(t, entity.Post{}, createdPost)
	mockRepo.AssertExpectations(t)
}

func TestPostUsecase_GetPosts_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(mocks.PostRepository)
	usecase := NewPostUsecase(mockRepo, logger)

	expectedPosts := []entity.Post{
		{
			ID:       1,
			AuthorId: 1,
			Title:    "Post 1",
			Content:  "Content 1",
		},
		{
			ID:       2,
			AuthorId: 2,
			Title:    "Post 2",
			Content:  "Content 2",
		},
	}

	mockRepo.On("GetPosts", mock.Anything, 10, 0).Return(expectedPosts, nil)

	posts, err := usecase.GetPosts(context.Background(), 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, expectedPosts, posts)
	mockRepo.AssertExpectations(t)
}

func TestPostUsecase_GetPosts_Failure(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(mocks.PostRepository)
	usecase := NewPostUsecase(mockRepo, logger)

	mockRepo.On("GetPosts", mock.Anything, 10, 0).Return(nil, assert.AnError)

	posts, err := usecase.GetPosts(context.Background(), 10, 0)
	assert.Error(t, err)
	assert.Nil(t, posts)
	mockRepo.AssertExpectations(t)
}

func TestPostUsecase_GetPostByID_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(mocks.PostRepository)
	usecase := NewPostUsecase(mockRepo, logger)

	expectedPost := &entity.Post{
		ID:       1,
		AuthorId: 1,
		Title:    "Test Post",
		Content:  "This is a test post",
	}

	mockRepo.On("GetPostByID", mock.Anything, 1).Return(expectedPost, nil)

	post, err := usecase.GetPostByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedPost, post)
	mockRepo.AssertExpectations(t)
}

func TestPostUsecase_GetPostByID_Failure(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(mocks.PostRepository)
	usecase := NewPostUsecase(mockRepo, logger)

	mockRepo.On("GetPostByID", mock.Anything, 1).Return(nil, assert.AnError)

	post, err := usecase.GetPostByID(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, post)
	mockRepo.AssertExpectations(t)
}

func TestPostUsecase_UpdatePost_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(mocks.PostRepository)
	usecase := NewPostUsecase(mockRepo, logger)

	post := entity.Post{
		ID:       1,
		AuthorId: 1,
		Title:    "Updated Post",
		Content:  "This is an updated post",
	}

	expectedPost := &entity.Post{
		ID:       1,
		AuthorId: 1,
		Title:    "Updated Post",
		Content:  "This is an updated post",
	}

	mockRepo.On("UpdatePost", mock.Anything, post).Return(expectedPost, nil)

	updatedPost, err := usecase.UpdatePost(context.Background(), post)
	assert.NoError(t, err)
	assert.Equal(t, expectedPost, updatedPost)
	mockRepo.AssertExpectations(t)
}

func TestPostUsecase_UpdatePost_Failure(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(mocks.PostRepository)
	usecase := NewPostUsecase(mockRepo, logger)

	post := entity.Post{
		ID:       1,
		AuthorId: 1,
		Title:    "Updated Post",
		Content:  "This is an updated post",
	}

	mockRepo.On("UpdatePost", mock.Anything, post).Return(nil, assert.AnError)

	updatedPost, err := usecase.UpdatePost(context.Background(), post)
	assert.Error(t, err)
	assert.Nil(t, updatedPost)
	mockRepo.AssertExpectations(t)
}

func TestPostUsecase_DeletePost_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(mocks.PostRepository)
	usecase := NewPostUsecase(mockRepo, logger)

	mockRepo.On("DeletePost", mock.Anything, 1).Return(nil)

	err := usecase.DeletePost(context.Background(), 1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPostUsecase_DeletePost_Failure(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(mocks.PostRepository)
	usecase := NewPostUsecase(mockRepo, logger)

	mockRepo.On("DeletePost", mock.Anything, 1).Return(assert.AnError)

	err := usecase.DeletePost(context.Background(), 1)
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPostUsecase_GetTotalPostsCount_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(mocks.PostRepository)
	usecase := NewPostUsecase(mockRepo, logger)

	mockRepo.On("GetTotalPostsCount", mock.Anything).Return(42, nil)

	count, err := usecase.GetTotalPostsCount(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 42, count)
	mockRepo.AssertExpectations(t)
}

func TestPostUsecase_GetTotalPostsCount_Failure(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockRepo := new(mocks.PostRepository)
	usecase := NewPostUsecase(mockRepo, logger)

	mockRepo.On("GetTotalPostsCount", mock.Anything).Return(0, errors.New("failed to count posts"))

	count, err := usecase.GetTotalPostsCount(context.Background())
	assert.Error(t, err)
	assert.Equal(t, 0, count)
	mockRepo.AssertExpectations(t)
}
