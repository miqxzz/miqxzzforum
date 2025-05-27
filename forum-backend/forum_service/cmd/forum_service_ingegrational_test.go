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
	"github.com/miqxzz/miqxzzforum/forum_service/internal/entity"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/repository"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/usecase"
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

func setupRouter() *gin.Engine {
	db := setupTestDB(nil)
	logger, _ := zap.NewProduction()
	router := gin.New()

	// Инициализация репозиториев и usecases
	postRepo := repository.NewPostRepository(db, logger)
	commentRepo := repository.NewCommentsRepository(db, logger)
	postUsecase := usecase.NewPostUsecase(postRepo, logger)

	// Регистрация маршрутов
	router.POST("/posts", func(c *gin.Context) {
		var post entity.Post
		if err := c.ShouldBindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := postUsecase.CreatePost(c.Request.Context(), post)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, result)
	})

	router.GET("/posts", func(c *gin.Context) {
		limit := 10
		offset := 0
		posts, err := postUsecase.GetPosts(c.Request.Context(), limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, posts)
	})

	router.POST("/posts/:id/comments", func(c *gin.Context) {
		var comment entity.Comment
		if err := c.ShouldBindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		comment.PostId = 1 // Для тестов
		result, err := commentRepo.CreateComment(c.Request.Context(), comment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, result)
	})

	router.GET("/posts/:id/comments", func(c *gin.Context) {
		comments, err := commentRepo.GetComments(c.Request.Context(), 1, 10, 0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, comments)
	})

	return router
}

func TestForumService_Integration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	t.Run("CreatePost", func(t *testing.T) {
		post := entity.Post{
			AuthorId: 1,
			Title:    "Test Post",
			Content:  "This is a test post",
		}

		body, _ := json.Marshal(post)
		req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer test-token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "Test Post")
	})

	t.Run("GetPosts", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/posts", nil)
		req.Header.Set("Authorization", "Bearer test-token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Test Post")
	})

	t.Run("CreateComment", func(t *testing.T) {
		comment := entity.Comment{
			AuthorId: 1,
			PostId:   1,
			Content:  "This is a test comment",
		}

		body, _ := json.Marshal(comment)
		req := httptest.NewRequest(http.MethodPost, "/posts/1/comments", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer test-token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "This is a test comment")
	})

	t.Run("GetCommentsByPostID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/posts/1/comments", nil)
		req.Header.Set("Authorization", "Bearer test-token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "This is a test comment")
	})
}
