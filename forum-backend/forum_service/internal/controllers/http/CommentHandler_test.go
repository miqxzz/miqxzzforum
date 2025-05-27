package http

import (
	"bytes"
	"encoding/json"
	"errors"
	utils "github.com/Engls/EnglsJwt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"github.com/Engls/forum-project2/forum_service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestCommentHandler_CreateComment_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger)

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

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger)

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

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger)

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

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger)

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

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger)

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

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger)

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

func TestCommentHandler_GetCommentsByPostID_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger)

	comments := []entity.Comment{
		{ID: 1, PostId: 1, Content: "Comment 1"},
		{ID: 2, PostId: 1, Content: "Comment 2"},
	}

	mockCommentUsecase.On("GetCommentByPostID", mock.Anything, 1).Return(comments, nil)

	req, _ := http.NewRequest("GET", "/posts/1/comments", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	commentHandler.GetCommentsByPostID(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseComments []entity.Comment
	err := json.Unmarshal(w.Body.Bytes(), &responseComments)
	assert.NoError(t, err)
	assert.Equal(t, comments, responseComments)

	mockCommentUsecase.AssertExpectations(t)
}

func TestCommentHandler_GetCommentsByPostID_InvalidPostID(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger)

	req, _ := http.NewRequest("GET", "/posts/invalid/comments", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "invalid"}}

	commentHandler.GetCommentsByPostID(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid post ID")
}

func TestCommentHandler_GetCommentsByPostID_FailedToGetComments(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockCommentUsecase := new(mocks.CommentsUsecases)
	jwtUtil := utils.NewJWTUtil("secret")

	commentHandler := NewCommentHandler(mockCommentUsecase, jwtUtil, logger)

	mockCommentUsecase.On("GetCommentByPostID", mock.Anything, 1).Return(nil, errors.New("failed to get comments"))

	req, _ := http.NewRequest("GET", "/posts/1/comments", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	commentHandler.GetCommentsByPostID(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to get comments")

	mockCommentUsecase.AssertExpectations(t)
}
