package usecase

import (
	"context"
	"errors"
	"testing"

	commonmiqx "github.com/miqxzz/commonmiqx"
	entity "github.com/miqxzz/miqxzzforum/auth_service/internal/entity"
	mocks "github.com/miqxzz/miqxzzforum/auth_service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Мок репозитория
type mockAuthRepo struct{ mock.Mock }

func (m *mockAuthRepo) Register(user entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}
func (m *mockAuthRepo) Login(username, password string) (string, error) { return "token", nil }
func (m *mockAuthRepo) GetUserRole(username string) (string, error)     { return "user", nil }
func (m *mockAuthRepo) UpdateUserRole(userID int, newRole string) error {
	args := m.Called(userID, newRole)
	return args.Error(0)
}
func (m *mockAuthRepo) GetUserByUsername(username string) (entity.User, error) {
	args := m.Called(username)
	return args.Get(0).(entity.User), args.Error(1)
}
func (m *mockAuthRepo) SaveToken(userID int, token string) error { return nil }
func (m *mockAuthRepo) GetUsernameByID(ctx context.Context, userID int) (string, error) {
	return "", nil
}

func TestAuthUsecase_Register_Success(t *testing.T) {
	logger, _ := zap.NewProduction()

	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := commonmiqx.NewJWTUtil("secret")

	username := "testuser"
	password := "12345"
	role := "user"

	mockAuthRepo.On("Register", mock.AnythingOfType("entity.User")).Return(nil)

	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	err := authUsecase.Register(username, password, role)

	assert.NoError(t, err)

	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUsecase_Register_InvalidPassword(t *testing.T) {
	logger, _ := zap.NewProduction()

	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := commonmiqx.NewJWTUtil("secret")

	testCases := []struct {
		name     string
		password string
		errorMsg string
	}{
		{
			name:     "Too short",
			password: "1234",
			errorMsg: "пароль должен содержать минимум 5 символов",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)
			err := authUsecase.Register("testuser", tc.password, "user")
			assert.Error(t, err)
			assert.Equal(t, tc.errorMsg, err.Error())
		})
	}
}

func TestAuthUsecase_Register_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := commonmiqx.NewJWTUtil("secret")

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
	jwtUtil := commonmiqx.NewJWTUtil("secret")

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
	jwtUtil := commonmiqx.NewJWTUtil("secret")

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
	jwtUtil := commonmiqx.NewJWTUtil("secret")

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
	jwtUtil := commonmiqx.NewJWTUtil("secret")

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
	jwtUtil := commonmiqx.NewJWTUtil("secret")

	username := "testuser"

	mockAuthRepo.On("GetUserByUsername", username).Return(entity.User{}, errors.New("user not found"))

	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	role, err := authUsecase.GetUserRole(username)

	assert.Error(t, err)
	assert.Equal(t, "", role)

	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUsecase_UpdateUserRole_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := commonmiqx.NewJWTUtil("secret")

	userID := 1
	newRole := "admin"

	mockAuthRepo.On("UpdateUserRole", userID, newRole).Return(nil)

	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	err := authUsecase.UpdateUserRole(userID, newRole)

	assert.NoError(t, err)
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUsecase_UpdateUserRole_InvalidRole(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := commonmiqx.NewJWTUtil("secret")

	userID := 1
	invalidRole := "invalid_role"

	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	err := authUsecase.UpdateUserRole(userID, invalidRole)

	assert.Error(t, err)
	assert.Equal(t, "недопустимая роль пользователя", err.Error())
	mockAuthRepo.AssertExpectations(t)
}

func TestAuthUsecase_UpdateUserRole_RepositoryError(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockAuthRepo := new(mocks.AuthRepository)
	jwtUtil := commonmiqx.NewJWTUtil("secret")

	userID := 1
	newRole := "admin"

	mockAuthRepo.On("UpdateUserRole", userID, newRole).Return(errors.New("database error"))

	authUsecase := NewAuthUsecase(mockAuthRepo, jwtUtil, logger)

	err := authUsecase.UpdateUserRole(userID, newRole)

	assert.Error(t, err)
	mockAuthRepo.AssertExpectations(t)
}

func TestRegister_Success(t *testing.T) {
	repo := new(mockAuthRepo)
	jwtUtil := commonmiqx.NewJWTUtil("secret")
	logger, _ := zap.NewProduction()
	uc := NewAuthUsecase(repo, jwtUtil, logger)

	repo.On("Register", mock.AnythingOfType("entity.User")).Return(nil)

	err := uc.Register("test", "12345", "user")
	assert.NoError(t, err)
}

func TestRegister_InvalidPassword(t *testing.T) {
	repo := new(mockAuthRepo)
	jwtUtil := commonmiqx.NewJWTUtil("secret")
	logger, _ := zap.NewProduction()
	uc := NewAuthUsecase(repo, jwtUtil, logger)

	err := uc.Register("test", "123", "user")
	assert.Error(t, err)
}

func TestRegister_RepoError(t *testing.T) {
	repo := new(mockAuthRepo)
	jwtUtil := commonmiqx.NewJWTUtil("secret")
	logger, _ := zap.NewProduction()
	uc := NewAuthUsecase(repo, jwtUtil, logger)

	repo.On("Register", mock.AnythingOfType("entity.User")).Return(errors.New("db error"))
	err := uc.Register("test", "12345", "user")
	assert.Error(t, err)
}

func TestUpdateUserRole_Success(t *testing.T) {
	repo := new(mockAuthRepo)
	jwtUtil := commonmiqx.NewJWTUtil("secret")
	logger, _ := zap.NewProduction()
	uc := NewAuthUsecase(repo, jwtUtil, logger)

	repo.On("UpdateUserRole", 1, "admin").Return(nil)
	err := uc.UpdateUserRole(1, "admin")
	assert.NoError(t, err)
}

func TestUpdateUserRole_InvalidRole(t *testing.T) {
	repo := new(mockAuthRepo)
	jwtUtil := commonmiqx.NewJWTUtil("secret")
	logger, _ := zap.NewProduction()
	uc := NewAuthUsecase(repo, jwtUtil, logger)

	err := uc.UpdateUserRole(1, "superuser")
	assert.Error(t, err)
}

func TestUpdateUserRole_RepoError(t *testing.T) {
	repo := new(mockAuthRepo)
	jwtUtil := commonmiqx.NewJWTUtil("secret")
	logger, _ := zap.NewProduction()
	uc := NewAuthUsecase(repo, jwtUtil, logger)

	repo.On("UpdateUserRole", 1, "admin").Return(errors.New("db error"))
	err := uc.UpdateUserRole(1, "admin")
	assert.Error(t, err)
}
