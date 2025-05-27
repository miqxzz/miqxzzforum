package repository

import (
	"context"
	"errors"
	"github.com/Engls/forum-project2/forum_service/internal/repository/adapters"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Engls/forum-project2/forum_service/internal/entity"
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
