package repository

import (
	"context"
	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"go.uber.org/zap"
)

type CommentsRepository interface {
	CreateComment(ctx context.Context, comment entity.Comment) (entity.Comment, error)
	GetComments(ctx context.Context, postID, limit, offset int) ([]entity.Comment, error)
	GetTotalCommentsCount(ctx context.Context, postID int) (int, error)
}

type commentsRepository struct {
	db     DB
	logger *zap.Logger
}

func NewCommentsRepository(db DB, logger *zap.Logger) CommentsRepository {
	return &commentsRepository{db: db, logger: logger}
}

func (r *commentsRepository) CreateComment(ctx context.Context, comment entity.Comment) (entity.Comment, error) {
	query := `
		INSERT INTO comments (post_id, author_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(ctx, query, comment.PostId, comment.AuthorId, comment.Content).Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		r.logger.Error("Failed to create comment", zap.Error(err), zap.Int("postID", comment.PostId), zap.Int("authorID", comment.AuthorId))
		return entity.Comment{}, err
	}
	r.logger.Info("Comment created successfully", zap.Int("commentID", comment.ID), zap.Int("postID", comment.PostId), zap.Int("authorID", comment.AuthorId))
	return comment, nil
}

func (r *commentsRepository) GetComments(ctx context.Context, postID, limit, offset int) ([]entity.Comment, error) {
	query := `
        SELECT id, content, author_id, post_id, created_at 
        FROM comments 
        WHERE post_id = $1 
        ORDER BY created_at DESC 
        LIMIT $2 OFFSET $3
    `
	rows, err := r.db.QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []entity.Comment
	for rows.Next() {
		var comment entity.Comment
		if err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.AuthorId,
			&comment.PostId,
			&comment.CreatedAt,
		); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (r *commentsRepository) GetTotalCommentsCount(ctx context.Context, postID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM comments WHERE post_id = $1`
	err := r.db.QueryRowContext(ctx, query, postID).Scan(&count)
	return count, err
}
