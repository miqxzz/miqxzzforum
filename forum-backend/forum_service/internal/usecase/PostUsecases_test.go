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

func TestPostUsecase_CreatePost_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockPostRepo := new(mocks.PostRepository)

	postUsecase := NewPostUsecase(mockPostRepo, logger)

	post := entity.Post{
		AuthorId: 1,
		Title:    "Test Post",
		Content:  "This is a test post",
	}
	createdPost := post
	createdPost.ID = 1

	mockPostRepo.On("CreatePost", mock.Anything, post).Return(&createdPost, nil)

	result, err := postUsecase.CreatePost(context.Background(), post)

	assert.NoError(t, err)
	assert.Equal(t, &createdPost, result)

	mockPostRepo.AssertExpectations(t)
}

func TestPostUsecase_CreatePost_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockPostRepo := new(mocks.PostRepository)

	postUsecase := NewPostUsecase(mockPostRepo, logger)

	post := entity.Post{
		AuthorId: 1,
		Title:    "Test Post",
		Content:  "This is a test post",
	}

	mockPostRepo.On("CreatePost", mock.Anything, post).Return(nil, errors.New("failed to create post"))

	result, err := postUsecase.CreatePost(context.Background(), post)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockPostRepo.AssertExpectations(t)
}

func TestPostUsecase_GetPosts_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockPostRepo := new(mocks.PostRepository)

	postUsecase := NewPostUsecase(mockPostRepo, logger)

	posts := []entity.Post{
		{ID: 1, AuthorId: 1, Title: "Post 1", Content: "Content 1"},
		{ID: 2, AuthorId: 2, Title: "Post 2", Content: "Content 2"},
	}

	mockPostRepo.On("GetPosts", mock.Anything).Return(posts, nil)

	result, err := postUsecase.GetPosts(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, posts, result)

	mockPostRepo.AssertExpectations(t)
}

func TestPostUsecase_GetPosts_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockPostRepo := new(mocks.PostRepository)

	postUsecase := NewPostUsecase(mockPostRepo, logger)

	mockPostRepo.On("GetPosts", mock.Anything).Return(nil, errors.New("failed to get posts"))

	result, err := postUsecase.GetPosts(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)

	mockPostRepo.AssertExpectations(t)
}

func TestPostUsecase_GetPostByID_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockPostRepo := new(mocks.PostRepository)

	postUsecase := NewPostUsecase(mockPostRepo, logger)

	post := entity.Post{
		ID:       1,
		AuthorId: 1,
		Title:    "Test Post",
		Content:  "This is a test post",
	}

	mockPostRepo.On("GetPostByID", mock.Anything, 1).Return(&post, nil)

	result, err := postUsecase.GetPostByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, &post, result)

	mockPostRepo.AssertExpectations(t)
}

func TestPostUsecase_GetPostByID_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockPostRepo := new(mocks.PostRepository)

	postUsecase := NewPostUsecase(mockPostRepo, logger)

	mockPostRepo.On("GetPostByID", mock.Anything, 1).Return(nil, errors.New("failed to get post"))

	result, err := postUsecase.GetPostByID(context.Background(), 1)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockPostRepo.AssertExpectations(t)
}

func TestPostUsecase_UpdatePost_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockPostRepo := new(mocks.PostRepository)

	postUsecase := NewPostUsecase(mockPostRepo, logger)

	post := entity.Post{
		ID:       1,
		AuthorId: 1,
		Title:    "Updated Post",
		Content:  "This is an updated post",
	}

	mockPostRepo.On("UpdatePost", mock.Anything, post).Return(&post, nil)

	result, err := postUsecase.UpdatePost(context.Background(), post)

	assert.NoError(t, err)
	assert.Equal(t, &post, result)

	mockPostRepo.AssertExpectations(t)
}

func TestPostUsecase_UpdatePost_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockPostRepo := new(mocks.PostRepository)

	postUsecase := NewPostUsecase(mockPostRepo, logger)

	post := entity.Post{
		ID:       1,
		AuthorId: 1,
		Title:    "Updated Post",
		Content:  "This is an updated post",
	}

	mockPostRepo.On("UpdatePost", mock.Anything, post).Return(nil, errors.New("failed to update post"))

	result, err := postUsecase.UpdatePost(context.Background(), post)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockPostRepo.AssertExpectations(t)
}

func TestPostUsecase_DeletePost_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockPostRepo := new(mocks.PostRepository)

	postUsecase := NewPostUsecase(mockPostRepo, logger)

	mockPostRepo.On("DeletePost", mock.Anything, 1).Return(nil)

	err := postUsecase.DeletePost(context.Background(), 1)

	assert.NoError(t, err)

	mockPostRepo.AssertExpectations(t)
}

func TestPostUsecase_DeletePost_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockPostRepo := new(mocks.PostRepository)

	postUsecase := NewPostUsecase(mockPostRepo, logger)

	mockPostRepo.On("DeletePost", mock.Anything, 1).Return(errors.New("failed to delete post"))

	err := postUsecase.DeletePost(context.Background(), 1)

	assert.Error(t, err)

	mockPostRepo.AssertExpectations(t)
}
