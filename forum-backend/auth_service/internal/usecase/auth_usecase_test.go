package usecase

import (
	"errors"
	"testing"

	"github.com/Engls/EnglsJwt"
	"github.com/Engls/forum-project2/auth_service/internal/entity"
	"github.com/Engls/forum-project2/auth_service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthUsecase_Register_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := EnglsJwt.NewJWTUtil("secret")

	username := "testuser"
	password := "password"
	role := "user"

	mockAuthRepo.On("Register", mock.AnythingOfType("entity.User")).Return(nil)

	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	err := authUsecase.Register(username, password, role)

	assert.NoError(t, err)

	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUsecase_Register_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := EnglsJwt.NewJWTUtil("secret")

	username := "testuser"
	password := "password"
	role := "user"

	mockAuthRepo.On("Register", mock.AnythingOfType("entity.User")).Return(errors.New("failed to register user"))

	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	err := authUsecase.Register(username, password, role)

	assert.Error(t, err)

	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := EnglsJwt.NewJWTUtil("secret")

	username := "testuser"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := entity.User{ID: 1, Username: username, Password: string(hashedPassword), Role: "user"}

	mockAuthRepo.On("GetUserByUsername", username).Return(user, nil)
	mockAuthRepo.On("SaveToken", user.ID, mock.Anything).Return(nil)

	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	resultToken, err := authUsecase.Login(username, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, resultToken)

	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login_Failure_InvalidCredentials(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := EnglsJwt.NewJWTUtil("secret")

	username := "testuser"
	password := "password"

	mockAuthRepo.On("GetUserByUsername", username).Return(entity.User{}, errors.New("user not found"))

	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	resultToken, err := authUsecase.Login(username, password)

	assert.Error(t, err)
	assert.Equal(t, "", resultToken)

	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login_Failure_InvalidPassword(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := EnglsJwt.NewJWTUtil("secret")

	username := "testuser"
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("wrongpassword"), bcrypt.DefaultCost)
	user := entity.User{ID: 1, Username: username, Password: string(hashedPassword), Role: "user"}

	mockAuthRepo.On("GetUserByUsername", username).Return(user, nil)

	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	resultToken, err := authUsecase.Login(username, password)

	assert.Error(t, err)
	assert.Equal(t, "", resultToken)

	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUsecase_GetUserRole_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := EnglsJwt.NewJWTUtil("secret")

	username := "testuser"
	user := entity.User{ID: 1, Username: username, Password: "hashedpassword", Role: "user"}

	mockAuthRepo.On("GetUserByUsername", username).Return(user, nil)

	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	role, err := authUsecase.GetUserRole(username)

	assert.NoError(t, err)
	assert.Equal(t, user.Role, role)

	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUsecase_GetUserRole_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := EnglsJwt.NewJWTUtil("secret")

	username := "testuser"

	mockAuthRepo.On("GetUserByUsername", username).Return(entity.User{}, errors.New("user not found"))

	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	role, err := authUsecase.GetUserRole(username)

	assert.Error(t, err)
	assert.Equal(t, "", role)

	mockAuthRepo.AssertExpectations(t)
}
