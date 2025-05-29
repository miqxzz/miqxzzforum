package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/entity"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/repository/adapters"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCommentsRepository_CreateComment_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()
	dbAdapter := adapters.DbAdapter{db}

	commentsRepo := NewCommentsRepository(&dbAdapter, logger)

	comment := entity.Comment{
		PostId:   1,
		AuthorId: 1,
		Content:  "This is a test comment",
	}
	createdComment := comment
	createdComment.ID = 1
	createdComment.CreatedAt = time.Now()

	mock.ExpectQuery(`INSERT INTO comments \(post_id, author_id, content\) VALUES \(\$1, \$2, \$3\) RETURNING id, created_at`).
		WithArgs(comment.PostId, comment.AuthorId, comment.Content).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(createdComment.ID, createdComment.CreatedAt))

	result, err := commentsRepo.CreateComment(context.Background(), comment)

	assert.NoError(t, err)
	assert.Equal(t, createdComment, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentsRepository_CreateComment_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := adapters.DbAdapter{db}

	commentsRepo := NewCommentsRepository(&dbAdapter, logger)

	comment := entity.Comment{
		PostId:   1,
		AuthorId: 1,
		Content:  "This is a test comment",
	}

	mock.ExpectQuery(`INSERT INTO comments \(post_id, author_id, content\) VALUES \(\$1, \$2, \$3\) RETURNING id, created_at`).
		WithArgs(comment.PostId, comment.AuthorId, comment.Content).
		WillReturnError(errors.New("failed to create comment"))

	result, err := commentsRepo.CreateComment(context.Background(), comment)

	assert.Error(t, err)
	assert.Equal(t, entity.Comment{}, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentsRepository_GetCommentsByPostID_Success(t *testing.T) {

	logger, _ := zap.NewProduction()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := adapters.DbAdapter{db}

	commentsRepo := NewCommentsRepository(&dbAdapter, logger)

	postID := 1
	comments := []entity.Comment{
		{ID: 1, PostId: postID, AuthorId: 1, Content: "Comment 1", CreatedAt: time.Now()},
		{ID: 2, PostId: postID, AuthorId: 2, Content: "Comment 2", CreatedAt: time.Now()},
	}

	rows := sqlmock.NewRows([]string{"id", "post_id", "author_id", "content", "created_at"})
	for _, comment := range comments {
		rows.AddRow(comment.ID, comment.PostId, comment.AuthorId, comment.Content, comment.CreatedAt)
	}
	mock.ExpectQuery(`SELECT id, post_id, author_id, content, created_at FROM comments WHERE post_id = \$1 ORDER BY created_at ASC`).
		WithArgs(postID).
		WillReturnRows(rows)

	result, err := commentsRepo.GetCommentsByPostID(context.Background(), postID)

	assert.NoError(t, err)
	assert.Equal(t, comments, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentsRepository_GetCommentsByPostID_Failure(t *testing.T) {

	logger, _ := zap.NewProduction()

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := adapters.DbAdapter{db}

	commentsRepo := NewCommentsRepository(&dbAdapter, logger)

	postID := 1

	mock.ExpectQuery(`SELECT id, post_id, author_id, content, created_at FROM comments WHERE post_id = \$1 ORDER BY created_at ASC`).
		WithArgs(postID).
		WillReturnError(errors.New("failed to get comments"))

	result, err := commentsRepo.GetCommentsByPostID(context.Background(), postID)

	assert.Error(t, err)
	assert.Nil(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentsRepository_GetCommentByID_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := adapters.DbAdapter{db}
	commentsRepo := NewCommentsRepository(&dbAdapter, logger)

	comment := entity.Comment{ID: 1, PostId: 1, AuthorId: 1, Content: "Test", CreatedAt: time.Now()}
	rows := sqlmock.NewRows([]string{"id", "post_id", "author_id", "content", "created_at"}).
		AddRow(comment.ID, comment.PostId, comment.AuthorId, comment.Content, comment.CreatedAt)
	mock.ExpectQuery(`SELECT id, post_id, author_id, content, created_at FROM comments WHERE id = \$1`).
		WithArgs(comment.ID).
		WillReturnRows(rows)

	result, err := commentsRepo.GetCommentByID(context.Background(), comment.ID)
	assert.NoError(t, err)
	assert.Equal(t, comment, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentsRepository_GetCommentByID_NotFound(t *testing.T) {
	logger, _ := zap.NewProduction()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := adapters.DbAdapter{db}
	commentsRepo := NewCommentsRepository(&dbAdapter, logger)

	mock.ExpectQuery(`SELECT id, post_id, author_id, content, created_at FROM comments WHERE id = \$1`).
		WithArgs(42).
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "author_id", "content", "created_at"}))

	result, err := commentsRepo.GetCommentByID(context.Background(), 42)
	assert.Error(t, err)
	assert.Equal(t, entity.Comment{}, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentsRepository_DeleteComment_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := adapters.DbAdapter{db}
	commentsRepo := NewCommentsRepository(&dbAdapter, logger)

	mock.ExpectExec(`DELETE FROM comments WHERE id = \$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	err = commentsRepo.DeleteComment(context.Background(), 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentsRepository_DeleteComment_Failure(t *testing.T) {
	logger, _ := zap.NewProduction()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := adapters.DbAdapter{db}
	commentsRepo := NewCommentsRepository(&dbAdapter, logger)

	mock.ExpectExec(`DELETE FROM comments WHERE id = \$1`).WithArgs(2).WillReturnError(errors.New("delete error"))

	err = commentsRepo.DeleteComment(context.Background(), 2)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentsRepository_GetTotalCommentsCount_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := adapters.DbAdapter{db}
	commentsRepo := NewCommentsRepository(&dbAdapter, logger)

	rows := sqlmock.NewRows([]string{"count"}).AddRow(5)
	mock.ExpectQuery(`SELECT COUNT\(\*\) as count FROM comments WHERE post_id = \$1`).WithArgs(1).WillReturnRows(rows)

	count, err := commentsRepo.GetTotalCommentsCount(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, 5, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCommentsRepository_GetTotalCommentsCount_Failure(t *testing.T) {
	logger, _ := zap.NewProduction()
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbAdapter := adapters.DbAdapter{db}
	commentsRepo := NewCommentsRepository(&dbAdapter, logger)

	mock.ExpectQuery(`SELECT COUNT\(\*\) as count FROM comments WHERE post_id = \$1`).WithArgs(1).WillReturnError(errors.New("count error"))

	count, err := commentsRepo.GetTotalCommentsCount(context.Background(), 1)
	assert.Error(t, err)
	assert.Equal(t, 0, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}
