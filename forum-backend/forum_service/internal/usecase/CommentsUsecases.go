package usecase

import (
	"context"

	"github.com/miqxzz/miqxzzforum/forum_service/internal/entity"
	"github.com/miqxzz/miqxzzforum/forum_service/internal/repository"
	"go.uber.org/zap"
)

type CommentsUsecases interface {
	CreateComment(ctx context.Context, comment entity.Comment) (entity.Comment, error)
	GetCommentByPostID(ctx context.Context, postId int) ([]entity.Comment, error)
	GetComments(ctx context.Context, postID, limit, offset int) ([]entity.Comment, error)
	GetTotalCommentsCount(ctx context.Context, postID int) (int, error)
	DeleteComment(ctx context.Context, id int) error
	GetCommentByID(ctx context.Context, id int) (entity.Comment, error)
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

func (u *commentsUsecases) GetCommentByPostID(ctx context.Context, postId int) ([]entity.Comment, error) {
	u.logger.Info("Fetching comments by post ID", zap.Int("postID", postId))

	comments, err := u.commentRepo.GetCommentsByPostID(ctx, postId)
	if err != nil {
		u.logger.Error("Failed to get comments by post ID", zap.Error(err), zap.Int("postID", postId))
		return nil, err
	}

	u.logger.Info("Comments fetched successfully", zap.Int("postID", postId), zap.Int("count", len(comments)))
	return comments, nil
}

func (u *commentsUsecases) DeleteComment(ctx context.Context, id int) error {
	u.logger.Info("Deleting comment", zap.Int("commentID", id))

	err := u.commentRepo.DeleteComment(ctx, id)
	if err != nil {
		u.logger.Error("Failed to delete comment", zap.Error(err), zap.Int("commentID", id))
		return err
	}

	u.logger.Info("Comment deleted successfully", zap.Int("commentID", id))
	return nil
}

func (u *commentsUsecases) GetCommentByID(ctx context.Context, id int) (entity.Comment, error) {
	u.logger.Info("Fetching comment by ID", zap.Int("commentID", id))

	comment, err := u.commentRepo.GetCommentByID(ctx, id)
	if err != nil {
		u.logger.Error("Failed to get comment by ID", zap.Error(err), zap.Int("commentID", id))
		return entity.Comment{}, err
	}

	u.logger.Info("Comment fetched successfully", zap.Int("commentID", id))
	return comment, nil
}
