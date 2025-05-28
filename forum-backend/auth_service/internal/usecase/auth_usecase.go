package usecase

import (
	"errors"

	utils "github.com/miqxzz/commonmiqx"
	entity "github.com/miqxzz/miqxzzforum/auth_service/internal/entity"
	repository "github.com/miqxzz/miqxzzforum/auth_service/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(username, password, role string) error
	Login(username, password string) (string, error)
	GetUserRole(username string) (string, error)
	UpdateUserRole(userID int, newRole string) error
}

type authUsecase struct {
	authRepo repository.AuthRepository
	jwtUtil  *utils.JWTUtil
	logger   *zap.Logger
}

func NewAuthUsecase(authRepo repository.AuthRepository, jwtUtil *utils.JWTUtil, logger *zap.Logger) AuthUsecase {
	return &authUsecase{authRepo: authRepo, jwtUtil: jwtUtil, logger: logger}
}

func validatePassword(password string) error {
	if len(password) < 5 {
		return errors.New("пароль должен содержать минимум 5 символов")
	}
	return nil
}

func (u *authUsecase) Register(username, password, role string) error {
	if err := validatePassword(password); err != nil {
		u.logger.Error("Invalid password format", zap.Error(err), zap.String("username", username))
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error("Failed to hash password", zap.Error(err), zap.String("username", username))
		return err
	}
	user := entity.User{Username: username, Password: string(hashedPassword), Role: role}
	if err := u.authRepo.Register(user); err != nil {
		u.logger.Error("Failed to register user", zap.Error(err), zap.String("username", username))
		return err
	}
	u.logger.Info("User registered successfully", zap.String("username", username))
	return nil
}

func (u *authUsecase) Login(username, password string) (string, error) {
	user, err := u.authRepo.GetUserByUsername(username)
	if err != nil {
		u.logger.Error("Failed to get user by username", zap.Error(err), zap.String("username", username))
		return "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		u.logger.Error("Invalid password", zap.String("username", username))
		return "", errors.New("invalid credentials")
	}
	token, err := u.jwtUtil.GenerateToken(user.ID, user.Role)
	if err != nil {
		u.logger.Error("Failed to generate token", zap.Error(err), zap.String("username", username))
		return "", err
	}
	if err := u.authRepo.SaveToken(user.ID, token); err != nil {
		u.logger.Error("Failed to save token", zap.Error(err), zap.String("username", username))
		return "", err
	}
	u.logger.Info("User logged in successfully", zap.String("username", username))
	return token, nil
}

func (u *authUsecase) GetUserRole(username string) (string, error) {
	user, err := u.authRepo.GetUserByUsername(username)
	if err != nil {
		u.logger.Error("Failed to get user role", zap.Error(err), zap.String("username", username))
		return "", errors.New("invalid credentials")
	}
	u.logger.Info("User role retrieved successfully", zap.String("username", username), zap.String("role", user.Role))
	return user.Role, nil
}

func (u *authUsecase) UpdateUserRole(userID int, newRole string) error {
	if newRole != "user" && newRole != "admin" && newRole != "moderator" {
		u.logger.Error("Invalid role", zap.String("role", newRole))
		return errors.New("недопустимая роль пользователя")
	}

	if err := u.authRepo.UpdateUserRole(userID, newRole); err != nil {
		u.logger.Error("Failed to update user role", zap.Error(err), zap.Int("userID", userID), zap.String("newRole", newRole))
		return err
	}
	u.logger.Info("User role updated successfully", zap.Int("userID", userID), zap.String("newRole", newRole))
	return nil
}
