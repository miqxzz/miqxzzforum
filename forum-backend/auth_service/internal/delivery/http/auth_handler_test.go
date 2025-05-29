package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	utils "github.com/miqxzz/commonmiqx"
	entity "github.com/miqxzz/miqxzzforum/auth_service/internal/entity"

	"github.com/gin-gonic/gin"
	_ "github.com/miqxzz/miqxzzforum/auth_service/internal/usecase"
	"github.com/miqxzz/miqxzzforum/auth_service/mocks"
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

func TestAuthHandler_UpdateUserRole_Success(t *testing.T) {
	// Создаем мок для AuthUsecase
	mockAuthUsecase := new(mocks.AuthUsecase)
	logger, _ := zap.NewProduction()
	jwtUtil := utils.NewJWTUtil("secret")

	// Настраиваем ожидаемое поведение
	mockAuthUsecase.On("GetUserIDFromToken", "valid-token").Return(1, nil)
	mockAuthUsecase.On("UpdateUserRole", 1, "moderator").Return(nil)

	// Создаем тестируемый обработчик
	handler := NewAuthHandler(mockAuthUsecase, jwtUtil, logger)

	// Создаем тестовый HTTP-запрос
	req := httptest.NewRequest("PUT", "/update-role", strings.NewReader(`{"role": "moderator"}`))
	req.Header.Set("Authorization", "Bearer valid-token")

	// Создаем тестовый HTTP-ответ
	w := httptest.NewRecorder()

	// Создаем тестовый маршрутизатор
	router := gin.New()
	router.PUT("/update-role", handler.UpdateUserRole)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)

	// Проверяем, что все ожидаемые вызовы были выполнены
	mockAuthUsecase.AssertExpectations(t)
}

func TestAuthHandler_UpdateUserRole_NoToken(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockAuthUsecase := new(mocks.AuthUsecase)
	jwtUtil := utils.NewJWTUtil("secret")

	authHandler := NewAuthHandler(mockAuthUsecase, jwtUtil, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/update-role", bytes.NewBufferString(`{"user_id":2,"new_role":"moderator"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	authHandler.UpdateUserRole(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "требуется авторизация")

	mockAuthUsecase.AssertExpectations(t)
}

func TestAuthHandler_UpdateUserRole_NotAdmin(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockAuthUsecase := new(mocks.AuthUsecase)
	jwtUtil := utils.NewJWTUtil("secret")

	// Генерируем токен для обычного пользователя
	token, _ := jwtUtil.GenerateToken(1, "user")

	authHandler := NewAuthHandler(mockAuthUsecase, jwtUtil, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/update-role", bytes.NewBufferString(`{"user_id":2,"new_role":"moderator"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("Authorization", token)

	authHandler.UpdateUserRole(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "недостаточно прав")

	mockAuthUsecase.AssertExpectations(t)
}

func TestAuthHandler_UpdateUserRole_InvalidRequest(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockAuthUsecase := new(mocks.AuthUsecase)
	jwtUtil := utils.NewJWTUtil("secret")

	// Генерируем токен для админа
	token, _ := jwtUtil.GenerateToken(1, "admin")

	authHandler := NewAuthHandler(mockAuthUsecase, jwtUtil, logger)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/update-role", bytes.NewBufferString(`invalid json`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("Authorization", token)

	authHandler.UpdateUserRole(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	mockAuthUsecase.AssertExpectations(t)
}

func TestAuthHandler_UpdateUserRole_Error(t *testing.T) {
	// Создаем мок для AuthUsecase
	mockAuthUsecase := new(mocks.AuthUsecase)
	logger, _ := zap.NewProduction()
	jwtUtil := utils.NewJWTUtil("secret")

	// Настраиваем ожидаемое поведение
	mockAuthUsecase.On("GetUserIDFromToken", "valid-token").Return(1, nil)
	mockAuthUsecase.On("UpdateUserRole", 1, "invalid-role").Return(errors.New("invalid role"))

	// Создаем тестируемый обработчик
	handler := NewAuthHandler(mockAuthUsecase, jwtUtil, logger)

	// Создаем тестовый HTTP-запрос
	req := httptest.NewRequest("PUT", "/update-role", strings.NewReader(`{"role": "invalid-role"}`))
	req.Header.Set("Authorization", "Bearer valid-token")

	// Создаем тестовый HTTP-ответ
	w := httptest.NewRecorder()

	// Создаем тестовый маршрутизатор
	router := gin.New()
	router.PUT("/update-role", handler.UpdateUserRole)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Проверяем, что все ожидаемые вызовы были выполнены
	mockAuthUsecase.AssertExpectations(t)
}

func TestAuthHandler_UpdateUserRole_InvalidToken(t *testing.T) {
	// Создаем мок для AuthUsecase
	mockAuthUsecase := new(mocks.AuthUsecase)
	logger, _ := zap.NewProduction()
	jwtUtil := utils.NewJWTUtil("secret")

	// Настраиваем ожидаемое поведение
	mockAuthUsecase.On("GetUserIDFromToken", "invalid-token").Return(0, errors.New("invalid token"))

	// Создаем тестируемый обработчик
	handler := NewAuthHandler(mockAuthUsecase, jwtUtil, logger)

	// Создаем тестовый HTTP-запрос
	req := httptest.NewRequest("PUT", "/update-role", strings.NewReader(`{"role": "moderator"}`))
	req.Header.Set("Authorization", "Bearer invalid-token")

	// Создаем тестовый HTTP-ответ
	w := httptest.NewRecorder()

	// Создаем тестовый маршрутизатор
	router := gin.New()
	router.PUT("/update-role", handler.UpdateUserRole)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Проверяем, что все ожидаемые вызовы были выполнены
	mockAuthUsecase.AssertExpectations(t)
}

func TestAuthHandler_UpdateUserRole_InvalidJSON(t *testing.T) {
	// Создаем мок для AuthUsecase
	mockAuthUsecase := new(mocks.AuthUsecase)
	logger, _ := zap.NewProduction()
	jwtUtil := utils.NewJWTUtil("secret")

	// Создаем тестируемый обработчик
	handler := NewAuthHandler(mockAuthUsecase, jwtUtil, logger)

	// Создаем тестовый HTTP-запрос с неверным JSON
	req := httptest.NewRequest("PUT", "/update-role", strings.NewReader(`{"role": "moderator"`))
	req.Header.Set("Authorization", "Bearer valid-token")

	// Создаем тестовый HTTP-ответ
	w := httptest.NewRecorder()

	// Создаем тестовый маршрутизатор
	router := gin.New()
	router.PUT("/update-role", handler.UpdateUserRole)

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Проверяем, что все ожидаемые вызовы были выполнены
	mockAuthUsecase.AssertExpectations(t)
}
