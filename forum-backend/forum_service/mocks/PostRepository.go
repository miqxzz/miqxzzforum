package mocks

import (
	"context"

	"github.com/miqxzz/miqxzzforum/forum_service/internal/entity"
	"github.com/stretchr/testify/mock"
)

type PostRepository struct {
	mock.Mock
}

func (m *PostRepository) CreatePost(ctx context.Context, post entity.Post) (*entity.Post, error) {
	args := m.Called(ctx, post)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *PostRepository) GetPosts(ctx context.Context, limit, offset int) ([]entity.Post, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]entity.Post), args.Error(1)
}

func (m *PostRepository) GetPostByID(ctx context.Context, id int) (*entity.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *PostRepository) UpdatePost(ctx context.Context, post entity.Post) (*entity.Post, error) {
	args := m.Called(ctx, post)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *PostRepository) DeletePost(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *PostRepository) GetUserIDByToken(ctx context.Context, token string) (int, error) {
	args := m.Called(ctx, token)
	return args.Int(0), args.Error(1)
}

func (m *PostRepository) GetTotalPostsCount(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}
