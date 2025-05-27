package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	utils "github.com/miqxzz/commonmiqx"

	"github.com/gin-gonic/gin"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/entity"
	"github.com/miqxzz/miqxzzforum/forum_service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestPostHandler_CreatePost_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger, nil)

	post := &entity.Post{
		Title:   "Test Post",
		Content: "This is a test post",
	}
	postJSON, _ := json.Marshal(post)

	mockPostRepo.On("GetUserIDByToken", mock.Anything, "valid.jwt.token").Return(1, nil)
	mockPostUsecase.On("CreatePost", mock.Anything, mock.Anything).Return(post, nil)

	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(postJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	postHandler.CreatePost(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Test Post")

	mockPostRepo.AssertExpectations(t)
	mockPostUsecase.AssertExpectations(t)
}

func TestPostHandler_CreatePost_MissingAuthorizationHeader(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger, nil)

	post := entity.Post{
		Title:   "Test Post",
		Content: "This is a test post",
	}
	postJSON, _ := json.Marshal(post)

	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(postJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	postHandler.CreatePost(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header required")
}

func TestPostHandler_CreatePost_InvalidAuthorizationHeaderFormat(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger, nil)

	post := entity.Post{
		Title:   "Test Post",
		Content: "This is a test post",
	}
	postJSON, _ := json.Marshal(post)

	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(postJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "InvalidFormat")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	postHandler.CreatePost(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid Authorization header format")
}

func TestPostHandler_CreatePost_InvalidToken(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger, nil)

	post := entity.Post{
		Title:   "Test Post",
		Content: "This is a test post",
	}
	postJSON, _ := json.Marshal(post)

	mockPostRepo.On("GetUserIDByToken", mock.Anything, "invalid.jwt.token\t").Return(0, errors.New("invalid token"))

	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(postJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid.jwt.token	")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	postHandler.CreatePost(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")

	mockPostRepo.AssertExpectations(t)
}

func TestPostHandler_CreatePost_FailedToCreatePost(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger, nil)

	post := &entity.Post{
		Title:   "Test Post",
		Content: "This is a test post",
	}
	postJSON, _ := json.Marshal(post)

	mockPostRepo.On("GetUserIDByToken", mock.Anything, "valid.jwt.token").Return(1, nil)
	mockPostUsecase.On("CreatePost", mock.Anything, mock.Anything).Return(post, errors.New("failed to create post"))

	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(postJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid.jwt.token")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	postHandler.CreatePost(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to create post")

	mockPostRepo.AssertExpectations(t)
	mockPostUsecase.AssertExpectations(t)
}

func TestPostHandler_GetPosts_Success(t *testing.T) {
	// Arrange
	logger, _ := zap.NewProduction()
	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	handler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger, nil)

	expectedPosts := []*entity.Post{
		{
			ID:       1,
			AuthorId: 1,
			Title:    "Test Post 1",
			Content:  "Content 1",
		},
		{
			ID:       2,
			AuthorId: 1,
			Title:    "Test Post 2",
			Content:  "Content 2",
		},
	}

	mockPostUsecase.On("GetPosts", mock.Anything, 10, 0).Return(expectedPosts, nil)

	// Act
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/posts", nil)

	handler.GetPosts(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Posts []*entity.Post `json:"posts"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedPosts, response.Posts)
}

func TestPostHandler_GetPosts_FailedToGetPosts(t *testing.T) {
	logger, _ := zap.NewProduction()

	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger, nil)

	mockPostUsecase.On("GetPosts", mock.Anything).Return(nil, errors.New("failed to get posts"))

	req, _ := http.NewRequest("GET", "/posts", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	postHandler.GetPosts(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to get posts")

	mockPostUsecase.AssertExpectations(t)
}

func TestPostHandler_DeletePost_MissingAuthorizationHeader(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger, nil)

	req, _ := http.NewRequest("DELETE", "/posts/1", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	postHandler.DeletePost(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header required")
}

func TestPostHandler_DeletePost_InvalidAuthorizationHeaderFormat(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger, nil)

	req, _ := http.NewRequest("DELETE", "/posts/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "InvalidFormat")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	postHandler.DeletePost(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid Authorization header format")
}

func TestPostHandler_DeletePost_Success_Owner(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger, nil)

	token, err := jwtUtil.GenerateToken(1, "user")
	assert.NoError(t, err)

	mockPostRepo.On("GetPostByID", mock.Anything, 1).Return(&entity.Post{ID: 1, AuthorId: 1}, nil)
	mockPostUsecase.On("DeletePost", mock.Anything, 1).Return(nil)

	req, _ := http.NewRequest("DELETE", "/posts/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	postHandler.DeletePost(c)

	assert.Equal(t, http.StatusOK, w.Code)

	mockPostRepo.AssertExpectations(t)
	mockPostUsecase.AssertExpectations(t)
}

func TestPostHandler_DeletePost_Success_Admin(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger, nil)

	token, err := jwtUtil.GenerateToken(1, "admin")
	assert.NoError(t, err)

	mockPostRepo.On("GetPostByID", mock.Anything, 1).Return(entity.Post{ID: 1, AuthorId: 2}, nil)
	mockPostUsecase.On("DeletePost", mock.Anything, 1).Return(nil)

	req, _ := http.NewRequest("DELETE", "/posts/1", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	postHandler.DeletePost(c)

	assert.Equal(t, http.StatusOK, w.Code)

	mockPostUsecase.AssertExpectations(t)
}

func TestPostHandler_GetTotalPostsCount_Success(t *testing.T) {
	logger, _ := zap.NewProduction()

	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger, nil)

	mockPostUsecase.On("GetTotalPostsCount", mock.Anything).Return(10, nil)

	req, _ := http.NewRequest("GET", "/posts/count", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	postHandler.GetTotalPostsCount(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]int
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 10, response["count"])

	mockPostUsecase.AssertExpectations(t)
}

func TestPostHandler_GetTotalPostsCount_Error(t *testing.T) {
	logger, _ := zap.NewProduction()

	mockPostUsecase := new(mocks.PostUsecase)
	mockPostRepo := new(mocks.PostRepository)
	jwtUtil := utils.NewJWTUtil("secret")

	postHandler := NewPostHandler(mockPostUsecase, mockPostRepo, jwtUtil, logger, nil)

	mockPostUsecase.On("GetTotalPostsCount", mock.Anything).Return(0, errors.New("failed to get count"))

	req, _ := http.NewRequest("GET", "/posts/count", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	postHandler.GetTotalPostsCount(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to get count")

	mockPostUsecase.AssertExpectations(t)
}
