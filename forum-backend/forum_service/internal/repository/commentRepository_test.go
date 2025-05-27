package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/entity"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/repository/adapters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCommentsRepository_CreateComment_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dbAdapter := &adapters.DbAdapter{DB: db}
	repo := NewCommentsRepository(dbAdapter, logger)

	comment := entity.Comment{
		PostId:   1,
		AuthorId: 1,
		Content:  "Test comment",
	}

	rows := sqlmock.NewRows([]string{"id", "created_at"}).
		AddRow(1, time.Now())

	mock.ExpectQuery("INSERT INTO comments").
		WithArgs(1, 1, "Test comment").
		WillReturnRows(rows)

	createdComment, err := repo.CreateComment(context.Background(), comment)
	assert.NoError(t, err)
	assert.Equal(t, 1, createdComment.ID)
	assert.Equal(t, 1, createdComment.AuthorId)
	assert.Equal(t, 1, createdComment.PostId)
	assert.Equal(t, "Test comment", createdComment.Content)
	assert.NotZero(t, createdComment.CreatedAt)
}

func TestCommentsRepository_CreateComment_Failure(t *testing.T) {
	logger, _ := zap.NewProduction()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dbAdapter := &adapters.DbAdapter{DB: db}
	repo := NewCommentsRepository(dbAdapter, logger)

	comment := entity.Comment{
		PostId:   1,
		AuthorId: 1,
		Content:  "Test comment",
	}

	mock.ExpectQuery("INSERT INTO comments").
		WithArgs(1, 1, "Test comment").
		WillReturnError(assert.AnError)

	createdComment, err := repo.CreateComment(context.Background(), comment)
	assert.Error(t, err)
	assert.Equal(t, entity.Comment{}, createdComment)
}

func TestCommentsRepository_GetComments_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dbAdapter := &adapters.DbAdapter{DB: db}
	repo := NewCommentsRepository(dbAdapter, logger)

	expectedComments := []entity.Comment{
		{
			ID:        1,
			AuthorId:  1,
			PostId:    1,
			Content:   "Comment 1",
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			AuthorId:  2,
			PostId:    1,
			Content:   "Comment 2",
			CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "author_id", "post_id", "content", "created_at"}).
		AddRow(expectedComments[0].ID, expectedComments[0].AuthorId, expectedComments[0].PostId, expectedComments[0].Content, expectedComments[0].CreatedAt).
		AddRow(expectedComments[1].ID, expectedComments[1].AuthorId, expectedComments[1].PostId, expectedComments[1].Content, expectedComments[1].CreatedAt)

	mock.ExpectQuery("SELECT id, author_id, post_id, content, created_at FROM comments").
		WithArgs(1, 10, 0).
		WillReturnRows(rows)

	comments, err := repo.GetComments(context.Background(), 1, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, expectedComments, comments)
}

func TestCommentsRepository_GetComments_Failure(t *testing.T) {
	logger, _ := zap.NewProduction()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dbAdapter := &adapters.DbAdapter{DB: db}
	repo := NewCommentsRepository(dbAdapter, logger)

	mock.ExpectQuery("SELECT id, author_id, post_id, content, created_at FROM comments").
		WithArgs(1, 10, 0).
		WillReturnError(assert.AnError)

	comments, err := repo.GetComments(context.Background(), 1, 10, 0)
	assert.Error(t, err)
	assert.Nil(t, comments)
}

func TestCommentsRepository_GetTotalCommentsCount_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dbAdapter := &adapters.DbAdapter{DB: db}
	repo := NewCommentsRepository(dbAdapter, logger)

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(5)

	mock.ExpectQuery("SELECT COUNT").
		WithArgs(1).
		WillReturnRows(rows)

	count, err := repo.GetTotalCommentsCount(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, 5, count)
}

func TestCommentsRepository_GetTotalCommentsCount_Failure(t *testing.T) {
	logger, _ := zap.NewProduction()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	dbAdapter := &adapters.DbAdapter{DB: db}
	repo := NewCommentsRepository(dbAdapter, logger)

	mock.ExpectQuery("SELECT COUNT").
		WithArgs(1).
		WillReturnError(assert.AnError)

	count, err := repo.GetTotalCommentsCount(context.Background(), 1)
	assert.Error(t, err)
	assert.Equal(t, 0, count)
}
