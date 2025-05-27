package http

import (
	"bytes"
	"encoding/json"
	"errors"
	utils "github.com/Engls/EnglsJwt"
	"github.com/Engls/forum-project2/auth_service/internal/entity"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/Engls/forum-project2/auth_service/internal/usecase"
	"github.com/Engls/forum-project2/auth_service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestAuthHandler_Register_Success(t *testing.T) {
	logger, _ := zap.NewProduction()

	mockAuthUsecase := new(mocks.AuthUsecase)
	mockAuthUsecase.On("Register", "testuser", "password", "user").Return(nil)

	authHandler := NewAuthHandler(mockAuthUsecase, nil, logger)

	router := gin.Default()
	router.POST("/register", authHandler.Register)

	reqBody := map[string]string{
		"username": "testuser",
		"password": "password",
		"role":     "user",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "User registered successfully")

	mockAuthUsecase.AssertExpectations(t)
}

func TestAuthHandler_Register_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockAuthUsecase := new(mocks.AuthUsecase)
	jwtUtil := utils.NewJWTUtil("secret")

	req := entity.RegisterRequest{
		Username: "testuser",
		Password: "password",
		Role:     "user",
	}

	mockAuthUsecase.On("Register", req.Username, req.Password, req.Role).Return(errors.New("failed to register user"))

	authHandler := NewAuthHandler(mockAuthUsecase, jwtUtil, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{"username":"testuser","password":"password","role":"user"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	authHandler.Register(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to register user")

	mockAuthUsecase.AssertExpectations(t)
}

func TestAuthHandler_Register_BadRequest(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockAuthUsecase := new(mocks.AuthUsecase)
	jwtUtil := utils.NewJWTUtil("secret")

	authHandler := NewAuthHandler(mockAuthUsecase, jwtUtil, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{"{"username":"testuser","password":"","role":""}`))
	c.Request.Header.Set("Content-Type", "application/json")

	authHandler.Register(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "error")

	mockAuthUsecase.AssertExpectations(t)
}

func TestAuthHandler_Login(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockAuthUsecase := new(mocks.AuthUsecase)
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo3LCJyb2xlIjoidXNlciIsImV4cCI6MTc0NTUxNzExOCwiaWF0IjoxNzQ1MjU3OTE4fQ.QmyaHsq-ruAyciGKkgCEgj0xsQZD1J5ER6CLjXhfgQc"
	mockAuthUsecase.On("Login", "testuser", "password").Return(token, nil)
	mockAuthUsecase.On("GetUserRole", "testuser").Return("user", nil)

	jwtUtil := utils.NewJWTUtil("your-secret-key")
	authHandler := NewAuthHandler(mockAuthUsecase, jwtUtil, logger)

	router := gin.Default()
	router.POST("/login", authHandler.Login)

	reqBody := map[string]string{
		"username": "testuser",
		"password": "password",
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")

	mockAuthUsecase.AssertExpectations(t)
}

func TestAuthHandler_Login_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockAuthUsecase := new(mocks.AuthUsecase)
	jwtUtil := utils.NewJWTUtil("secret")

	req := entity.LoginRequest{
		Username: "testuser",
		Password: "password",
	}

	mockAuthUsecase.On("Login", req.Username, req.Password).Return("", errors.New("invalid credentials"))

	authHandler := NewAuthHandler(mockAuthUsecase, jwtUtil, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{"username":"testuser","password":"password"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	authHandler.Login(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid credentials")

	mockAuthUsecase.AssertExpectations(t)
}

func TestAuthHandler_Login_GetUserIDFromTokenFailure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockAuthUsecase := new(mocks.AuthUsecase)
	jwtUtil := utils.NewJWTUtil("secret")

	req := entity.LoginRequest{
		Username: "testuser",
		Password: "password",
	}
	token := "invalid.jwt.token"
	role := "user"

	mockAuthUsecase.On("Login", req.Username, req.Password).Return(token, nil)
	mockAuthUsecase.On("GetUserRole", req.Username).Return(role, nil)

	authHandler := NewAuthHandler(mockAuthUsecase, jwtUtil, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{"username":"testuser","password":"password"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	authHandler.Login(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "invalid character")

	mockAuthUsecase.AssertExpectations(t)
}
