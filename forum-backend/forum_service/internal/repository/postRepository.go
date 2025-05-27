package repository

import (
	"context"
	"github.com/Engls/forum-project2/forum_service/internal/entity"
	"go.uber.org/zap"
)

type DBposts interface {
	Get(dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type PostRepository interface {
	CreatePost(ctx context.Context, post entity.Post) (*entity.Post, error)
	GetPosts(ctx context.Context, limit, offset int) ([]entity.Post, error)
	GetPostByID(ctx context.Context, id int) (*entity.Post, error)
	UpdatePost(ctx context.Context, post entity.Post) (*entity.Post, error)
	DeletePost(ctx context.Context, id int) error
	GetUserIDByToken(ctx context.Context, token string) (int, error)
	GetTotalPostsCount(ctx context.Context) (int, error)
}

type postRepository struct {
	db     DB
	logger *zap.Logger
}

func NewPostRepository(db DB, logger *zap.Logger) PostRepository {
	return &postRepository{db: db, logger: logger}
}

func (r *postRepository) CreatePost(ctx context.Context, post entity.Post) (*entity.Post, error) {
	query := `INSERT INTO posts (author_id, title, content) VALUES (?, ?, ?)`
	result, err := r.db.ExecContext(ctx, query, post.AuthorId, post.Title, post.Content)
	if err != nil {
		r.logger.Error("Failed to create post", zap.Error(err), zap.Int("authorID", post.AuthorId))
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		r.logger.Error("Failed to get last insert ID", zap.Error(err))
		return nil, err
	}

	post.ID = int(id)
	r.logger.Info("Post created successfully", zap.Int("postID", post.ID), zap.Int("authorID", post.AuthorId))
	return &post, nil
}

func (r *postRepository) GetPosts(ctx context.Context, limit, offset int) ([]entity.Post, error) {
	query := `SELECT id, title, content, author_id FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var post entity.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorId); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *postRepository) GetTotalPostsCount(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM posts`).Scan(&count)
	return count, err
}

func (r *postRepository) GetPostByID(ctx context.Context, id int) (*entity.Post, error) {
	query := `SELECT id, author_id, title, content FROM posts WHERE id = ?`
	var post entity.Post
	err := r.db.GetContext(ctx, &post, query, id)
	if err != nil {
		r.logger.Error("Failed to get post by ID", zap.Error(err), zap.Int("postID", id))
		return nil, err
	}
	r.logger.Info("Post retrieved successfully", zap.Int("postID", id))
	return &post, nil
}

func (r *postRepository) UpdatePost(ctx context.Context, post entity.Post) (*entity.Post, error) {
	query := `UPDATE posts SET title = ?, content = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, post.Title, post.Content, post.ID) //TODO посмотреть что происходит при обновлении не сущ. записи
	if err != nil {
		r.logger.Error("Failed to update post", zap.Error(err), zap.Int("postID", post.ID))
		return nil, err
	}
	r.logger.Info("Post updated successfully", zap.Int("postID", post.ID))
	return &post, nil
}

func (r *postRepository) DeletePost(ctx context.Context, id int) error {
	query := `DELETE FROM posts WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete post", zap.Error(err), zap.Int("postID", id))
		return err
	}
	r.logger.Info("Post deleted successfully", zap.Int("postID", id))
	return nil
}

func (r *postRepository) GetUserIDByToken(ctx context.Context, token string) (int, error) {
	query := `SELECT user_id FROM tokens WHERE token = ?`
	var userID int
	err := r.db.GetContext(ctx, &userID, query, token)
	if err != nil {
		r.logger.Error("Failed to get user ID by token", zap.Error(err), zap.String("token", token))
		return 0, err
	}
	r.logger.Info("User ID retrieved successfully", zap.String("token", token), zap.Int("userID", userID))
	return userID, nil
}
