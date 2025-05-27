package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Engls/EnglsJwt"
	"github.com/Engls/forum-project2/forum_service/internal/controllers/chat"
	http2 "github.com/Engls/forum-project2/forum_service/internal/controllers/http"
	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"github.com/Engls/forum-project2/forum_service/internal/repository"
	"github.com/Engls/forum-project2/forum_service/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
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
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			role TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			author_id INTEGER,
			title TEXT,
			content TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (author_id) REFERENCES users(id)
		);
		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			author_id INTEGER,
			post_id INTEGER,
			content TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS chat_messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			username TEXT NOT NULL,
			content TEXT NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS tokens (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			token TEXT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create tables: %s", err)
	}

	return db
}

func TestForumService_Integration(t *testing.T) {

	logger, _ := zap.NewProduction()

	db := setupTestDB(t)
	defer db.Close()

	postRepo := repository.NewPostRepository(db, logger)
	commentRepo := repository.NewCommentsRepository(db, logger)
	chatRepo := repository.NewChatRepository(db, logger)
	postUsecase := usecase.NewPostUsecase(postRepo, logger)
	commentUsecase := usecase.NewCommentsUsecases(commentRepo, logger)
	hub := chat.NewHub()
	chatUsecase := usecase.NewChatUsecase(chatRepo, logger)
	jwtUtil := EnglsJwt.NewJWTUtil("secret")

	postHandler := http2.NewPostHandler(postUsecase, postRepo, jwtUtil, logger)
	commentHandler := http2.NewCommentHandler(commentUsecase, jwtUtil, logger)
	chatHandler := http2.NewChatHandler(hub, chatUsecase, jwtUtil, logger)

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Content-type", "Origin", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/ws", chatHandler.ServeWS)
	router.POST("/posts", postHandler.CreatePost)
	router.GET("/posts", postHandler.GetPosts)
	router.DELETE("/posts/:id", postHandler.DeletePost)
	router.POST("/posts/:id/comments", commentHandler.CreateComment)
	router.GET("/posts/:id/comments", commentHandler.GetCommentsByPostID)

	token, err := jwtUtil.GenerateToken(1, "user")
	if err != nil {
		t.Fatalf("Failed to generate token: %s", err)
	}

	_, err = db.Exec("INSERT INTO tokens (user_id, token) VALUES (?, ?)", 1, token)
	if err != nil {
		t.Fatalf("Failed to save token: %s", err)
	}

	t.Run("CreatePost", func(t *testing.T) {
		reqBody := entity.Post{
			AuthorId: 1,
			Title:    "Test Post",
			Content:  "This is a test post",
		}
		reqBodyBytes, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "This is a test post")
	})

	t.Run("GetPosts", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/posts", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Test Post")
	})

	t.Run("CreateComment", func(t *testing.T) {
		reqBody := entity.Comment{
			PostId:   1,
			AuthorId: 1,
			Content:  "This is a test comment",
		}
		reqBodyBytes, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/posts/1/comments", bytes.NewBuffer(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "This is a test comment")
	})

	t.Run("GetCommentsByPostID", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/posts/1/comments", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "This is a test comment")
	})
}
