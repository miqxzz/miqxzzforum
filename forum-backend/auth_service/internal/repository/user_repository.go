package repository

import (
	"context"
	"database/sql"
	"github.com/Engls/forum-project2/auth_service/internal/entity"
	"go.uber.org/zap"
)

type DB interface {
	Exec(query string, args ...any) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type AuthRepository interface {
	Register(user entity.User) error
	GetUserByUsername(username string) (entity.User, error)
	SaveToken(userID int, token string) error
	GetUsernameByID(ctx context.Context, userID int) (string, error)
}

type authRepository struct {
	db     DB
	logger *zap.Logger
}

func NewAuthRepository(db DB, logger *zap.Logger) AuthRepository {
	return &authRepository{db: db, logger: logger}
}

func (r *authRepository) Register(user entity.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", user.Username, user.Password, user.Role)
	if err != nil {
		r.logger.Error("Failed to register user", zap.Error(err), zap.String("username", user.Username))
		return err
	}
	r.logger.Info("User registered successfully", zap.String("username", user.Username))
	return nil
}

func (r *authRepository) GetUserByUsername(username string) (entity.User, error) {
	var user entity.User
	err := r.db.Get(&user, "SELECT id, username, password, role FROM users WHERE username=?", username)
	if err != nil {
		r.logger.Error("Failed to get user by username", zap.Error(err), zap.String("username", username))
		return user, err
	}
	r.logger.Info("User retrieved successfully", zap.String("username", username))
	return user, nil
}

func (r *authRepository) SaveToken(userID int, token string) error {
	_, err := r.db.Exec("INSERT INTO tokens (user_id, token) VALUES (?, ?)", userID, token)
	if err != nil {
		r.logger.Error("Failed to save token", zap.Error(err), zap.Int("userID", userID))
		return err
	}
	r.logger.Info("Token saved successfully", zap.Int("userID", userID))
	return nil
}

func (r *authRepository) GetUsernameByID(ctx context.Context, userID int) (string, error) {
	var username string
	err := r.db.QueryRowContext(ctx, "SELECT username FROM users WHERE id = $1", userID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return username, nil
}
