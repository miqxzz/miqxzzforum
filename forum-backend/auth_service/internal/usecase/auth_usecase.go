package usecase

import (
	"errors"
	utils "github.com/Engls/EnglsJwt"
	"github.com/Engls/forum-project2/auth_service/internal/entity"
	"github.com/Engls/forum-project2/auth_service/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(username, password, role string) error
	Login(username, password string) (string, error)
	GetUserRole(username string) (string, error)
}

type authUsecase struct {
	authRepo repository.AuthRepository
	jwtUtil  *utils.JWTUtil
	logger   *zap.Logger
}

func NewAuthUsecase(authRepo repository.AuthRepository, jwtUtil *utils.JWTUtil, logger *zap.Logger) AuthUsecase {
	return &authUsecase{authRepo: authRepo, jwtUtil: jwtUtil, logger: logger}
}

func (u *authUsecase) Register(username, password, role string) error {
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
