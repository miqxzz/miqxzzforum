package mocks

import (
	"context"

	"github.com/miqxzz/miqxzzforum/forum_service/internal/controllers/grpc"
	"github.com/stretchr/testify/mock"
)

type UserClient struct {
	mock.Mock
}

func NewUserClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *grpc.UserClient {
	mock := &UserClient{}
	mock.Mock.Test(t)
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return &grpc.UserClient{}
}

// GetUsername provides a mock function with given fields: ctx, userID
func (_m *UserClient) GetUsername(ctx context.Context, userID int) (string, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetUsername")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (string, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) string); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Close provides a mock function with given fields:
func (_m *UserClient) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
