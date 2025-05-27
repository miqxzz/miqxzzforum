package repository

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"github.com/Engls/forum-project2/forum_service/internal/repository/adapters"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestPostRepository_CreatePost_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := adapters.DbAdapter{db}

	postRepo := NewPostRepository(&dbAdapter, logger)

	post := entity.Post{
		AuthorId: 1,
		Title:    "Test Post",
		Content:  "This is a test post",
	}
	createdPost := post
	createdPost.ID = 1

	mock.ExpectExec(`INSERT INTO posts \(author_id, title, content\) VALUES \(\?, \?, \?\)`).
		WithArgs(post.AuthorId, post.Title, post.Content).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := postRepo.CreatePost(context.Background(), post)

	assert.NoError(t, err)
	assert.Equal(t, &createdPost, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_CreatePost_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := &adapters.DbAdapter{DB: db}

	postRepo := NewPostRepository(dbAdapter, logger)

	post := entity.Post{
		AuthorId: 1,
		Title:    "Test Post",
		Content:  "This is a test post",
	}

	mock.ExpectExec(`INSERT INTO posts \(author_id, title, content\) VALUES \(\?, \?, \?\)`).
		WithArgs(post.AuthorId, post.Title, post.Content).
		WillReturnError(errors.New("failed to create post"))

	result, err := postRepo.CreatePost(context.Background(), post)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetPosts_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := &adapters.DbAdapter{DB: db}

	postRepo := NewPostRepository(dbAdapter, logger)

	posts := []entity.Post{
		{ID: 1, AuthorId: 1, Title: "Post 1", Content: "Content 1"},
		{ID: 2, AuthorId: 2, Title: "Post 2", Content: "Content 2"},
	}

	rows := sqlmock.NewRows([]string{"id", "author_id", "title", "content"})
	for _, post := range posts {
		rows.AddRow(post.ID, post.AuthorId, post.Title, post.Content)
	}
	mock.ExpectQuery(`SELECT id, author_id, title, content FROM posts`).
		WillReturnRows(rows)

	result, err := postRepo.GetPosts(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, posts, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetPosts_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := &adapters.DbAdapter{DB: db}

	postRepo := NewPostRepository(dbAdapter, logger)

	mock.ExpectQuery(`SELECT id, author_id, title, content FROM posts`).
		WillReturnError(errors.New("failed to get posts"))

	result, err := postRepo.GetPosts(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetPostByID_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := &adapters.DbAdapter{DB: db}

	postRepo := NewPostRepository(dbAdapter, logger)

	postID := 1

	mock.ExpectQuery(`SELECT id, author_id, title, content FROM posts WHERE id = \?`).
		WithArgs(postID).
		WillReturnError(errors.New("failed to get post"))

	result, err := postRepo.GetPostByID(context.Background(), postID)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_UpdatePost_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := &adapters.DbAdapter{DB: db}

	postRepo := NewPostRepository(dbAdapter, logger)

	post := entity.Post{
		ID:       1,
		AuthorId: 1,
		Title:    "Updated Post",
		Content:  "This is an updated post",
	}

	mock.ExpectExec(`UPDATE posts SET title = \?, content = \? WHERE id = \?`).
		WithArgs(post.Title, post.Content, post.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := postRepo.UpdatePost(context.Background(), post)

	assert.NoError(t, err)
	assert.Equal(t, &post, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_UpdatePost_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := &adapters.DbAdapter{DB: db}

	postRepo := NewPostRepository(dbAdapter, logger)

	post := entity.Post{
		ID:       1,
		AuthorId: 1,
		Title:    "Updated Post",
		Content:  "This is an updated post",
	}

	mock.ExpectExec(`UPDATE posts SET title = \?, content = \? WHERE id = \?`).
		WithArgs(post.Title, post.Content, post.ID).
		WillReturnError(errors.New("failed to update post"))

	result, err := postRepo.UpdatePost(context.Background(), post)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_DeletePost_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := &adapters.DbAdapter{DB: db}

	postRepo := NewPostRepository(dbAdapter, logger)

	postID := 1

	mock.ExpectExec(`DELETE FROM posts WHERE id = \?`).
		WithArgs(postID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = postRepo.DeletePost(context.Background(), postID)

	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_DeletePost_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := &adapters.DbAdapter{DB: db}

	postRepo := NewPostRepository(dbAdapter, logger)

	postID := 1

	mock.ExpectExec(`DELETE FROM posts WHERE id = \?`).
		WithArgs(postID).
		WillReturnError(errors.New("failed to delete post"))

	err = postRepo.DeletePost(context.Background(), postID)

	assert.Error(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetUserIDByToken_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := &adapters.DbAdapter{DB: db}

	postRepo := NewPostRepository(dbAdapter, logger)

	token := "test-token"
	userID := 1

	rows := sqlmock.NewRows([]string{"user_id"}).AddRow(userID)
	mock.ExpectQuery(`SELECT user_id FROM tokens WHERE token = \?`).
		WithArgs(token).
		WillReturnRows(rows)

	result, err := postRepo.GetUserIDByToken(context.Background(), token)

	assert.NoError(t, err)
	assert.Equal(t, userID, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostRepository_GetUserIDByToken_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := &adapters.DbAdapter{DB: db}

	postRepo := NewPostRepository(dbAdapter, logger)

	token := "test-token"

	mock.ExpectQuery(`SELECT user_id FROM tokens WHERE token = \?`).
		WithArgs(token).
		WillReturnError(errors.New("failed to get user ID by token"))

	result, err := postRepo.GetUserIDByToken(context.Background(), token)

	assert.Error(t, err)
	assert.Equal(t, 0, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}
