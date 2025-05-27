package usecase

import (
	"context"
	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"github.com/Engls/forum-project2/forum_service/internal/repository"
	"go.uber.org/zap"
)

type PostUsecase interface {
	CreatePost(ctx context.Context, post entity.Post) (*entity.Post, error)
	GetPosts(ctx context.Context, limit, offset int) ([]entity.Post, error)
	GetPostByID(ctx context.Context, id int) (*entity.Post, error)
	UpdatePost(ctx context.Context, post entity.Post) (*entity.Post, error)
	DeletePost(ctx context.Context, id int) error
	GetTotalPostsCount(ctx context.Context) (int, error)
}

type postUsecase struct {
	postRepo repository.PostRepository
	logger   *zap.Logger
}

func NewPostUsecase(postRepo repository.PostRepository, logger *zap.Logger) PostUsecase {
	return &postUsecase{postRepo: postRepo, logger: logger}
}

func (u *postUsecase) CreatePost(ctx context.Context, post entity.Post) (*entity.Post, error) {
	u.logger.Info("Creating post",
		zap.Int("authorID", post.AuthorId),
		zap.String("title", post.Title),
		zap.String("content", post.Content),
	)

	createdPost, err := u.postRepo.CreatePost(ctx, post)
	if err != nil {
		u.logger.Error("Failed to create post", zap.Error(err))
		return nil, err
	}

	u.logger.Info("Post created successfully", zap.Int("postID", createdPost.ID))
	return createdPost, nil
}

func (u *postUsecase) GetPosts(ctx context.Context, limit, offset int) ([]entity.Post, error) {
	return u.postRepo.GetPosts(ctx, limit, offset)
}

func (u *postUsecase) GetTotalPostsCount(ctx context.Context) (int, error) {
	return u.postRepo.GetTotalPostsCount(ctx)
}

func (u *postUsecase) GetPostByID(ctx context.Context, id int) (*entity.Post, error) {
	u.logger.Info("Fetching post by ID", zap.Int("postID", id))

	post, err := u.postRepo.GetPostByID(ctx, id)
	if err != nil {
		u.logger.Error("Failed to get post by ID", zap.Error(err), zap.Int("postID", id))
		return nil, err
	}

	u.logger.Info("Post fetched successfully", zap.Int("postID", id))
	return post, nil
}

func (u *postUsecase) UpdatePost(ctx context.Context, post entity.Post) (*entity.Post, error) {
	u.logger.Info("Updating post",
		zap.Int("postID", post.ID),
		zap.String("title", post.Title),
		zap.String("content", post.Content),
	)

	updatedPost, err := u.postRepo.UpdatePost(ctx, post)
	if err != nil {
		u.logger.Error("Failed to update post", zap.Error(err), zap.Int("postID", post.ID))
		return nil, err
	}

	u.logger.Info("Post updated successfully", zap.Int("postID", post.ID))
	return updatedPost, nil
}

func (u *postUsecase) DeletePost(ctx context.Context, id int) error {
	u.logger.Info("Deleting post", zap.Int("postID", id))

	err := u.postRepo.DeletePost(ctx, id)
	if err != nil {
		u.logger.Error("Failed to delete post", zap.Error(err), zap.Int("postID", id))
		return err
	}

	u.logger.Info("Post deleted successfully", zap.Int("postID", id))
	return nil
}
