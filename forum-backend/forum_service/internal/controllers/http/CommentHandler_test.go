package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	utils "github.com/miqxzz/commonmiqx"

	"github.com/gin-gonic/gin"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/entity"
	"github.com/miqxzz/miqxzzforum/forum_service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestCommentHandler_CreateComment_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger, nil)

	token, err := jwtUtil.GenerateToken(1, "user")
	assert.NoError(t, err)

	comment := entity.Comment{
		Content: "This is a test comment",
	}
	commentJSON, _ := json.Marshal(comment)

	mockCommentUsecase.On("CreateComment", mock.Anything, mock.Anything).Return(comment, nil)

	req, _ := http.NewRequest("POST", "/posts/1/comments", bytes.NewBuffer(commentJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	commentHandler.CreateComment(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	var responseComment entity.Comment
	err = json.Unmarshal(w.Body.Bytes(), &responseComment)
	assert.NoError(t, err)
	assert.Equal(t, comment, responseComment)

	mockCommentUsecase.AssertExpectations(t)
}

func TestCommentHandler_CreateComment_InvalidPostID(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger, nil)

	token, err := jwtUtil.GenerateToken(1, "user")
	assert.NoError(t, err)

	comment := entity.Comment{
		Content: "This is a test comment",
	}
	commentJSON, _ := json.Marshal(comment)

	req, _ := http.NewRequest("POST", "/posts/invalid/comments", bytes.NewBuffer(commentJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid"}}

	commentHandler.CreateComment(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid post ID")
}

func TestCommentHandler_CreateComment_MissingAuthorizationHeader(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger, nil)

	comment := entity.Comment{
		Content: "This is a test comment",
	}
	commentJSON, _ := json.Marshal(comment)

	req, _ := http.NewRequest("POST", "/posts/1/comments", bytes.NewBuffer(commentJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	commentHandler.CreateComment(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header required")
}

func TestCommentHandler_CreateComment_InvalidAuthorizationHeaderFormat(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger, nil)

	comment := entity.Comment{
		Content: "This is a test comment",
	}
	commentJSON, _ := json.Marshal(comment)

	req, _ := http.NewRequest("POST", "/posts/1/comments", bytes.NewBuffer(commentJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "InvalidFormat")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	commentHandler.CreateComment(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid Authorization header format")
}

func TestCommentHandler_CreateComment_InvalidToken(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger, nil)

	comment := entity.Comment{
		Content: "This is a test comment",
	}
	commentJSON, _ := json.Marshal(comment)

	req, _ := http.NewRequest("POST", "/posts/1/comments", bytes.NewBuffer(commentJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer invalid.jwt.token")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	commentHandler.CreateComment(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token or user ID")
}

func TestCommentHandler_CreateComment_FailedToCreateComment(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger, nil)

	token, err := jwtUtil.GenerateToken(1, "user")
	assert.NoError(t, err)

	comment := entity.Comment{
		Content: "This is a test comment",
	}
	commentJSON, _ := json.Marshal(comment)

	mockCommentUsecase.On("CreateComment", mock.Anything, mock.Anything).Return(entity.Comment{}, errors.New("failed to create comment"))

	req, _ := http.NewRequest("POST", "/posts/1/comments", bytes.NewBuffer(commentJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	commentHandler.CreateComment(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to create comment")

	mockCommentUsecase.AssertExpectations(t)
}

func TestCommentHandler_GetComments_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	expectedComments := []*entity.Comment{
		{
			ID:        1,
			PostId:    1,
			Content:   "Comment 1",
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			PostId:    1,
			Content:   "Comment 2",
			CreatedAt: time.Now(),
		},
	}

	mockCommentUsecase.On("GetComments", mock.Anything, "1").Return(expectedComments, nil)

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("GET", "/posts/1/comments", nil)

	commentHandler.GetComments(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Comments []*entity.Comment `json:"comments"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedComments, response.Comments)

	mockCommentUsecase.AssertExpectations(t)
}

func TestCommentHandler_GetComments_InvalidPostID(t *testing.T) {
	logger, _ := zap.NewProduction()

	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger, nil)

	req, _ := http.NewRequest("GET", "/posts/invalid/comments", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid"}}

	commentHandler.GetComments(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid post ID")
}

func TestCommentHandler_GetComments_FailedToGetComments(t *testing.T) {
	logger, _ := zap.NewProduction()

	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger, nil)

	mockCommentUsecase.On("GetComments", mock.Anything, 1).Return(nil, errors.New("failed to get comments"))

	req, _ := http.NewRequest("GET", "/posts/1/comments", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	commentHandler.GetComments(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to get comments")

	mockCommentUsecase.AssertExpectations(t)
}

func TestCommentHandler_GetTotalCommentsCount_Success(t *testing.T) {
	logger, _ := zap.NewProduction()

	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger, nil)

	mockCommentUsecase.On("GetTotalCommentsCount", mock.Anything, 1).Return(10, nil)

	req, _ := http.NewRequest("GET", "/posts/1/comments/count", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	commentHandler.GetTotalCommentsCount(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]int
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 10, response["count"])

	mockCommentUsecase.AssertExpectations(t)
}

func TestCommentHandler_GetTotalCommentsCount_InvalidPostID(t *testing.T) {
	logger, _ := zap.NewProduction()

	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger, nil)

	req, _ := http.NewRequest("GET", "/posts/invalid/comments/count", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid"}}

	commentHandler.GetTotalCommentsCount(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid post ID")
}

func TestCommentHandler_GetTotalCommentsCount_Error(t *testing.T) {
	logger, _ := zap.NewProduction()

	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger, nil)

	mockCommentUsecase.On("GetTotalCommentsCount", mock.Anything, 1).Return(0, errors.New("failed to get count"))

	req, _ := http.NewRequest("GET", "/posts/1/comments/count", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	commentHandler.GetTotalCommentsCount(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to get count")

	mockCommentUsecase.AssertExpectations(t)
}
