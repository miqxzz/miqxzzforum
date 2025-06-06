package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	commonmiqx "github.com/miqxzz/commonmiqx"
	http2 "github.com/miqxzz/miqxzzforum/auth_service/internal/delivery/http"
	"github.com/miqxzz/miqxzzforum/auth_service/internal/entity"
	"github.com/miqxzz/miqxzzforum/auth_service/internal/repository"
	"github.com/miqxzz/miqxzzforum/auth_service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	tmpfile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatalf("Failed to create temp file: %s", err)
	}
	defer tmpfile.Close()

	db, err := sqlx.Open("sqlite3", tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to open database: %s", err)
	}

	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			role TEXT NOT NULL
		);
		CREATE TABLE tokens (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			token TEXT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create tables: %s", err)
	}

	return db
}

func TestAuthService_Integration(t *testing.T) {

	logger, _ := zap.NewProduction()

	db := setupTestDB(t)
	defer db.Close()

	authRepo := repository.NewAuthRepository(db, logger)
	jwtUtil := commonmiqx.NewJWTUtil("secret")
	authUsecase := usecase.NewAuthUsecase(authRepo, jwtUtil, logger)
	authHandler := http2.NewAuthHandler(authUsecase, jwtUtil, logger)

	r := gin.Default()
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	t.Run("RegisterUser", func(t *testing.T) {
		reqBody := entity.RegisterRequest{
			Username: "testuser",
			Password: "password",
			Role:     "user",
		}
		reqBodyBytes, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "User registered successfully")
	})

	t.Run("LoginUser", func(t *testing.T) {
		reqBody := entity.LoginRequest{
			Username: "testuser",
			Password: "password",
		}
		reqBodyBytes, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "token")
		assert.Contains(t, w.Body.String(), "role")
		assert.Contains(t, w.Body.String(), "username")
		assert.Contains(t, w.Body.String(), "userID")
	})

	t.Run("RegisterUser_InvalidPassword", func(t *testing.T) {
		reqBody := entity.RegisterRequest{
			Username: "shortpass",
			Password: "123",
			Role:     "user",
		}
		reqBodyBytes, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "пароль должен содержать минимум 5 символов")
	})

	t.Run("RegisterUser_AlreadyExists", func(t *testing.T) {
		reqBody := entity.RegisterRequest{
			Username: "testuser",
			Password: "password",
			Role:     "user",
		}
		reqBodyBytes, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "UNIQUE constraint failed")
	})

	t.Run("LoginUser_WrongPassword", func(t *testing.T) {
		reqBody := entity.LoginRequest{
			Username: "testuser",
			Password: "wrongpassword",
		}
		reqBodyBytes, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid credentials")
	})

	t.Run("LoginUser_NotExists", func(t *testing.T) {
		reqBody := entity.LoginRequest{
			Username: "nouser",
			Password: "password",
		}
		reqBodyBytes, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "invalid credentials")
	})
}
