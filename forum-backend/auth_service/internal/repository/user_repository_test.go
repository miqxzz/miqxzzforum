package repository

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/Engls/forum-project2/auth_service/internal/entity"
	"github.com/Engls/forum-project2/auth_service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestAuthRepository_Register_Success(t *testing.T) {
	logger, _ := zap.NewProduction()

	mockDB := new(mocks.DB)

	user := entity.User{Username: "testuser", Password: "hashedpassword", Role: "user"}

	mockDB.On("Exec", "INSERT INTO users (username, password, role) VALUES (?, ?, ?)", user.Username, user.Password, user.Role).Return(sql.Result(nil), nil)

	authRepo := NewAuthRepository(mockDB, logger)

	err := authRepo.Register(user)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

func TestAuthRepository_Register_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockDB := new(mocks.DB)

	user := entity.User{Username: "testuser", Password: "hashedpassword", Role: "user"}

	mockDB.On("Exec", "INSERT INTO users (username, password, role) VALUES (?, ?, ?)", user.Username, user.Password, user.Role).Return(nil, errors.New("failed to register user"))

	authRepo := NewAuthRepository(mockDB, logger)

	err := authRepo.Register(user)

	assert.Error(t, err)

	mockDB.AssertExpectations(t)
}

func TestAuthRepository_GetUserByUsername_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockDB := new(mocks.DB)

	user := entity.User{ID: 1, Username: "testuser", Password: "hashedpassword", Role: "user"}

	mockDB.On("Get", mock.Anything, "SELECT id, username, password, role FROM users WHERE username=?", user.Username).Run(func(args mock.Arguments) {
		dest := args.Get(0).(*entity.User)
		*dest = user
	}).Return(nil)

	authRepo := NewAuthRepository(mockDB, logger)

	result, err := authRepo.GetUserByUsername(user.Username)

	assert.NoError(t, err)
	assert.Equal(t, user, result)

	mockDB.AssertExpectations(t)
}

func TestAuthRepository_GetUserByUsername_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockDB := new(mocks.DB)

	username := "testuser"

	mockDB.On("Get", mock.Anything, "SELECT id, username, password, role FROM users WHERE username=?", username).Return(errors.New("failed to get user by username"))

	authRepo := NewAuthRepository(mockDB, logger)

	result, err := authRepo.GetUserByUsername(username)

	assert.Error(t, err)
	assert.Equal(t, entity.User{}, result)

	mockDB.AssertExpectations(t)
}

func TestAuthRepository_SaveToken_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockDB := new(mocks.DB)

	userID := 1
	token := "valid.jwt.token"

	mockDB.On("Exec", "INSERT INTO tokens (user_id, token) VALUES (?, ?)", userID, token).Return(sql.Result(nil), nil)

	authRepo := NewAuthRepository(mockDB, logger)

	err := authRepo.SaveToken(userID, token)

	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

func TestAuthRepository_SaveToken_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockDB := new(mocks.DB)

	userID := 1
	token := "valid.jwt.token"

	mockDB.On("Exec", "INSERT INTO tokens (user_id, token) VALUES (?, ?)", userID, token).Return(nil, errors.New("failed to save token"))

	authRepo := NewAuthRepository(mockDB, logger)

	err := authRepo.SaveToken(userID, token)

	assert.Error(t, err)

	mockDB.AssertExpectations(t)
}
