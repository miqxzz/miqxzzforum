package usecase

import (
	"context"
	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"github.com/Engls/forum-project2/forum_service/internal/repository"
	"go.uber.org/zap"
)

type CommentsUsecases interface {
	CreateComment(ctx context.Context, comment entity.Comment) (entity.Comment, error)
	GetComments(ctx context.Context, postID, limit, offset int) ([]entity.Comment, error)
	GetTotalCommentsCount(ctx context.Context, postID int) (int, error)
}

type commentsUsecases struct {
	commentRepo repository.CommentsRepository
	logger      *zap.Logger
}

func NewCommentsUsecases(commentRepo repository.CommentsRepository, logger *zap.Logger) CommentsUsecases {
	return &commentsUsecases{commentRepo: commentRepo, logger: logger}
}

func (u *commentsUsecases) CreateComment(ctx context.Context, comment entity.Comment) (entity.Comment, error) {
	u.logger.Info("Creating comment",
		zap.Int("postID", comment.PostId),
		zap.Int("authorID", comment.AuthorId),
		zap.String("content", comment.Content),
	)

	createdComment, err := u.commentRepo.CreateComment(ctx, comment)
	if err != nil {
		u.logger.Error("Failed to create comment", zap.Error(err))
		return entity.Comment{}, err
	}

	u.logger.Info("Comment created successfully", zap.Int("commentID", createdComment.ID), zap.Int("postID", createdComment.PostId))
	return createdComment, nil
}

func (u *commentsUsecases) GetComments(ctx context.Context, postID, limit, offset int) ([]entity.Comment, error) {
	return u.commentRepo.GetComments(ctx, postID, limit, offset)
}

func (u *commentsUsecases) GetTotalCommentsCount(ctx context.Context, postID int) (int, error) {
	return u.commentRepo.GetTotalCommentsCount(ctx, postID)
}
